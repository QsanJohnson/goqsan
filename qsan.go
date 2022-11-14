// @2022 QSAN Inc. All rights reserved

package goqsan

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang/glog"
)

const (
	defaultHttpPort  = 80
	defaultHttpsPort = 443
)

// QSAN client without authentication
type Client struct {
	apiKey     string
	baseURL    string
	HTTPClient *http.Client
}

// ClientOptions are options for QSAN http client.
type ClientOptions struct {
	Https      bool
	Port       int
	ReqTimeout time.Duration
}

// QSAN client with authentication
type AuthClient struct {
	Client
	user, passwd string
	accessToken  string
	refreshToken string
}

// For authentication
type AuthRes struct {
	AccessToken  string `json:"accessToken"`
	ExpireTime   int    `json:"expireTime"`
	RefreshToken string `json:"refreshToken"`
}

// Empty response data
type EmptyData []interface{}

type errorResponse struct {
	Error struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error"`
}

type RestError struct {
	ReqUrl     string
	StatusCode int
	ErrResp    errorResponse
	Err        error
}

func (r *RestError) Error() string {
	return fmt.Sprintf("[%s] status %d: %v (%d)", r.ReqUrl, r.StatusCode, r.ErrResp.Error.Message, r.ErrResp.Error.Code)
}

// NewClient returns QSAN client with given URL
func NewClient(ip string, opts ClientOptions) *Client {
	client := &Client{}
	if opts.Https {
		port := defaultHttpsPort
		if opts.Port != 0 {
			port = opts.Port
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &Client{
			HTTPClient: &http.Client{Transport: tr},
			baseURL:    fmt.Sprintf("https://%s:%d", ip, port),
		}
	} else {
		port := defaultHttpPort
		if opts.Port != 0 {
			port = opts.Port
		}

		client = &Client{
			HTTPClient: &http.Client{},
			baseURL:    fmt.Sprintf("http://%s:%d", ip, port),
		}
	}

	if opts.ReqTimeout != 0 {
		client.HTTPClient.Timeout = opts.ReqTimeout
	}

	return client
}

// If body format is url.Values, then body data will be sent using x-www-form-urlencoded format.
// If body format is string, then body data will be sent using raw data with JSON format.
func (c *Client) NewRequest(ctx context.Context, method, urlPath string, body interface{}) (*http.Request, error) {
	var (
		req *http.Request
		err error
	)

	urlStr := c.baseURL + urlPath
	glog.V(2).Infof("[NewRequest] %s url: %s\n", method, urlStr)
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	if body != nil {
		glog.V(3).Infof("[NewRequest] body: %v\n", body)
		switch body := body.(type) {
		case url.Values:
			req, err = http.NewRequest(method, u.String(), strings.NewReader(body.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		case string:
			// raw data
			req, err = http.NewRequest(method, u.String(), strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
		default:
			return nil, fmt.Errorf("Unknow body format! Only url.Values and string formats are supported.\n")
		}
	} else {
		req, err = http.NewRequest(method, u.String(), nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *AuthClient) SendRequest(ctx context.Context, req *http.Request, v interface{}) error {
	resterr := RestError{ReqUrl: req.Host + req.URL.Path}
	res, err := c.doSendRequest(ctx, req, v)
	if err != nil {
		resterr.Err = err
		return &resterr
	}

	resterr.StatusCode = res.StatusCode
	if res.StatusCode == 401 {
		res.Body.Close()

		if req.URL.Path != "/auth/refresh" {
			// When the existing access token expired, generate a new access token.
			glog.V(2).Infof("[AuthSendRequest] generate new access token. (%s%s)\n", req.Host, req.URL.Path)
			authRes, err := c.genAccessToken(ctx, c.refreshToken)
			if err != nil {
				resterr.Err = fmt.Errorf("genAccessToken failed: %v\n", err)
				return &resterr
			}

			// Update new access token then send request again
			c.accessToken = authRes.AccessToken
			c.apiKey = authRes.AccessToken
			glog.V(2).Infof("[AuthSendRequest] SendRequest again (%s%s)\n", req.Host, req.URL.Path)
			res, err = c.doSendRequest(ctx, req, v)
		} else {
			// When refresh token expired, renew a new access token and refresh token.
			glog.V(2).Infof("[AuthSendRequest] renew new access token and refresh token.\n")
			res, err := c.login(ctx, c.user, c.passwd)
			if err != nil {
				resterr.Err = fmt.Errorf("renew access token failed: %v\n", err)
				return &resterr
			}

			// Update new access token and refresh token
			c.accessToken = res.AccessToken
			c.apiKey = res.AccessToken
			c.refreshToken = res.RefreshToken

			authRes, ok := v.(*AuthRes)
			if ok {
				*authRes = *res
			} else {
				glog.Errorf("[AuthSendRequest] Should no be here. (%s%s)\n", req.Host, req.URL.Path)
			}

			return nil
		}

	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		errRes := errorResponse{}
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			resterr.ErrResp = errRes
			return &resterr
		}

		resterr.Err = fmt.Errorf("unknown error, status code: %d", res.StatusCode)
		return &resterr
	}

	if err = json.NewDecoder(res.Body).Decode(v); err != nil {
		resterr.Err = err
		return &resterr
	}

	return nil

}

func (c *Client) SendRequest(ctx context.Context, req *http.Request, v interface{}) error {
	res, err := c.doSendRequest(ctx, req, v)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		errRes := errorResponse{}
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Error.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(v); err != nil {
		return err
	}

	return err

}

func (c *Client) doSendRequest(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	if c.apiKey != "" {
		glog.V(5).Infof("[doSendRequest] apiKey: %s\n", c.apiKey)
		req.Header.Set("Authorization", c.apiKey)
	}

	req = req.WithContext(ctx)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		glog.Errorf("[doSendRequest] err: %v\n", err)
		return nil, err
	}

	glog.V(4).Infof("[doSendRequest] StatusCode: %d (%s%s)\n", res.StatusCode, req.Host, req.URL.Path)
	return res, nil
}

func (c *Client) login(ctx context.Context, user string, passwd string) (*AuthRes, error) {
	params := url.Values{}
	params.Add("user", user)
	params.Add("password", passwd)
	params.Add("offlineAccess", "true")

	req, err := c.NewRequest(ctx, http.MethodPost, "/auth/get", params)
	if err != nil {
		return nil, err
	}

	res := AuthRes{}
	if err := c.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Generate a new access token from refresh token
func (c *AuthClient) genAccessToken(ctx context.Context, t string) (*AuthRes, error) {
	params := url.Values{}
	params.Add("refreshToken", t)

	req, err := c.NewRequest(ctx, http.MethodPost, "/auth/refresh", params)
	if err != nil {
		return nil, err
	}

	res := AuthRes{}
	if err := c.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetAuthClient(ctx context.Context, user string, passwd string) (*AuthClient, error) {
	res, err := c.login(ctx, user, passwd)
	if err != nil {
		return nil, fmt.Errorf("login failed: %v\n", err)
	}

	glog.V(3).Infof("AccessToken: %s\n", res.AccessToken)

	return &AuthClient{
		Client: Client{
			apiKey:     res.AccessToken,
			baseURL:    c.baseURL,
			HTTPClient: c.HTTPClient,
		},
		user:         user,
		passwd:       passwd,
		accessToken:  res.AccessToken,
		refreshToken: res.RefreshToken,
	}, nil
}
