// @2022 QSAN Inc. All rights reserved

package goqsan

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	//"reflect"
)

// VolumeOp handles volume related methods of the QSAN storage.
type VolumeOp struct {
	client *AuthClient
}

type VolumeData struct {
	ID                    string      `json:"id"`
	Name                  string      `json:"name"`
	PoolID                string      `json:"poolId"`
	LunID                 interface{} `json:"lunId"`
	Online                bool        `json:"online"`
	Health                string      `json:"health"`
	Provision             string      `json:"provision"`
	TotalSize             uint64      `json:"totalSize"`
	UsedSize              uint64      `json:"usedSize"`
	BlockSize             uint64      `json:"blockSize"`
	StripeSize            uint64      `json:"stripeSize"`
	CacheMode             string      `json:"cacheMode"`
	IoPriority            string      `json:"ioPriority"`
	BgIoPriority          string      `json:"bgIoPriority"`
	EnableReadAhead       bool        `json:"enableReadAhead"`
	EraseData             string      `json:"eraseData"`
	EnableFastRaidRebuild bool        `json:"enableFastRaidRebuild"`
	Tags                  struct {
		Wwn  string `json:"wwn"`
		Type string `json:"type"`
	} `json:"tags"`
}

type VolumeCreateOptions struct {
	BlockSize    uint64 // recordsize: 1024, 2048 ..., 65536
	PoolId       uint64
	IoPriority   string
	BgIoPriority string
	CacheMode    string
}

type VolumeModifyOptions struct {
	Name            string `json:"name,omitempty"`
	UsedSize        uint64 `json:"usedSize,omitempty"`
	IoPriority      string `json:"ioPriority,omitempty"`
	BgIoPriority    string `json:"bgIoPriority,omitempty"`
	CacheMode       string `json:"cacheMode,omitempty"`
	EnableReadAhead *bool  `json:"enableReadAhead,omitempty"`
}

// NewVolume returns volume operation
func NewVolume(client *AuthClient) *VolumeOp {
	return &VolumeOp{client}
}

// ListVolumes list all volumes or a dedicated volume with volId
func (v *VolumeOp) ListVolumes(ctx context.Context, volId string) (*[]VolumeData, error) {

	req, err := v.client.NewRequest(ctx, http.MethodGet, "/rest/v2/storage/block/volumes/"+volId, nil)
	if err != nil {
		return nil, err
	}

	if volId == "" {
		//list All volumes
		res := []VolumeData{}
		if err := v.client.SendRequest(ctx, req, &res); err != nil {
			return nil, err
		}
		return &res, nil
	} else {
		//list certain volume
		singleres := VolumeData{}
		if err := v.client.SendRequest(ctx, req, &singleres); err != nil {
			return nil, err
		}
		res := []VolumeData{singleres}
		return &res, nil
	}
}

// list volumes under given PoolID
func (v *VolumeOp) ListVolumesByPoolID(ctx context.Context, poolID string) (*[]VolumeData, error) {

	req, err := v.client.NewRequest(ctx, http.MethodGet, "/rest/v2/storage/block/volumes?q=poolId='"+poolID+"'", nil)
	if err != nil {
		return nil, err
	}

	res := []VolumeData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// CreateVolume create a volume on a storage container
func (v *VolumeOp) CreateVolume(ctx context.Context, name string, size uint64, options *VolumeCreateOptions) (*VolumeData, error) {
	params := url.Values{}
	params.Add("name", name)
	params.Add("usedSize", strconv.FormatUint(size, 10))

	var optionMap map[string]interface{}
	data, _ := json.Marshal(options)
	json.Unmarshal(data, &optionMap)

	// 'int' type will become 'float64' type after struct to map[string]interface{} conversion
	if optionMap["BlockSize"] != float64(0) {
		v := uint64(optionMap["BlockSize"].(float64))
		params.Add("blockSize", strconv.FormatUint(v, 10))
	}

	if optionMap["PoolId"] != "" {
		v := uint64(optionMap["PoolId"].(float64))
		params.Add("poolId", strconv.FormatUint(v, 10))
	}

	req, err := v.client.NewRequest(ctx, http.MethodPost, "/rest/v2/storage/block/volumes", params)
	if err != nil {
		return nil, err
	}
	res := VolumeData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// DeleteVolume delete a volume from a storage container
func (v *VolumeOp) DeleteVolume(ctx context.Context, volId string) error {
	req, err := v.client.NewRequest(ctx, http.MethodDelete, "/rest/v2/storage/block/volumes/"+volId, nil)
	if err != nil {
		return err
	}

	res := EmptyData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return err
	}

	return nil
}

func (v *VolumeOp) ModifyVolume(ctx context.Context, volId string, options *VolumeModifyOptions) (*VolumeData, error) {
	rawdata, _ := json.Marshal(options)
	req, err := v.client.NewRequest(ctx, http.MethodPatch, "/rest/v2/storage/block/volumes/"+volId, string(rawdata))
	if err != nil {
		return nil, err
	}

	res := VolumeData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
