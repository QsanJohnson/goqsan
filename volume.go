// @2022 QSAN Inc. All rights reserved

package goqsan

import (
	"context"
	"encoding/json"
	"net/http"
)

// VolumeOp handles volume related methods of the QSAN storage.
type VolumeOp struct {
	client *AuthClient
}

type VolumeData struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	PoolID                string `json:"poolId"`
	LunID                 string `json:"lunId"`
	TargetID              string `json:"targetId"`
	Online                bool   `json:"online"`
	Health                string `json:"health"`
	Provision             string `json:"provision"`
	TotalSize             uint64 `json:"totalSize"`
	UsedSize              uint64 `json:"usedSize"`
	BlockSize             uint64 `json:"blockSize"`
	StripeSize            uint64 `json:"stripeSize"`
	CacheMode             string `json:"cacheMode"`
	IoPriority            string `json:"ioPriority"`
	BgIoPriority          string `json:"bgIoPriority"`
	EnableReadAhead       bool   `json:"enableReadAhead"`
	EraseData             string `json:"eraseData"`
	EnableFastRaidRebuild bool   `json:"enableFastRaidRebuild"`
	Tags                  struct {
		Wwn  string `json:"wwn"`
		Type string `json:"type"`
	} `json:"tags"`
}

type VolumeCreateOptions struct {
	Name            string `json:"name"`
	TotalSize       uint64 `json:"totalSize"`
	BlockSize       uint64 `json:"blockSize"`
	PoolID          string `json:"poolId"`
	IoPriority      string `json:"ioPriority,omitempty"`
	BgIoPriority    string `json:"bgIoPriority,omitempty"`
	CacheMode       string `json:"cacheMode,omitempty"`
	EnableReadAhead *bool  `json:"enableReadAhead,omitempty"`
}

type VolumeModifyOptions struct {
	Name            string `json:"name,omitempty"`
	TotalSize       uint64 `json:"totalSize,omitempty"`
	IoPriority      string `json:"ioPriority,omitempty"`
	BgIoPriority    string `json:"bgIoPriority,omitempty"`
	CacheMode       string `json:"cacheMode,omitempty"`
	EnableReadAhead *bool  `json:"enableReadAhead,omitempty"`
}

//return value of GET /rest/v2/storage/qos/volumes
//Patch /rest/v2/storage/qos/volumes
type QoSData struct {
	EnableQos bool   `json:"enableQos"`
	QosRule   string `json:"qosRule"`
}

// NewVolume returns volume operation
func NewVolume(client *AuthClient) *VolumeOp {
	return &VolumeOp{client}
}

// ListVolumes list all volumes
func (v *VolumeOp) ListVolumes(ctx context.Context) (*[]VolumeData, error) {

	req, err := v.client.NewRequest(ctx, http.MethodGet, "/rest/v2/storage/block/volumes", nil)
	if err != nil {
		return nil, err
	}

	res := []VolumeData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// ListVolumeByID list a dedicated volume with volId
func (v *VolumeOp) ListVolumeByID(ctx context.Context, volId string) (*VolumeData, error) {

	req, err := v.client.NewRequest(ctx, http.MethodGet, "/rest/v2/storage/block/volumes/"+volId, nil)
	if err != nil {
		return nil, err
	}

	res := VolumeData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		// resterr, ok := err.(*RestError)
		// if ok {
		// 	fmt.Printf("[ListVolumeByID] StatusCode=%d ErrResp=%+v\n", resterr.StatusCode, resterr.ErrResp)
		// }
		return nil, err
	}
	return &res, nil
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
func (v *VolumeOp) CreateVolume(ctx context.Context, poolID, volname string, volsize uint64, options *VolumeCreateOptions) (*VolumeData, error) {

	options.PoolID = poolID
	options.Name = volname
	options.TotalSize = volsize

	rawdata, _ := json.Marshal(options)
	req, err := v.client.NewRequest(ctx, http.MethodPost, "/rest/v2/storage/block/volumes", string(rawdata))
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

func (v *VolumeOp) GetQoS(ctx context.Context) (*QoSData, error) {

	req, err := v.client.NewRequest(ctx, http.MethodGet, "/rest/v2/storage/qos/volumes", nil)
	if err != nil {
		return nil, err
	}

	res := QoSData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (v *VolumeOp) PatchQoS(ctx context.Context, qosEnable bool, qosRule string) (*QoSData, error) {

	options := QoSData{}
	options.EnableQos = qosEnable
	options.QosRule = qosRule

	rawdata, _ := json.Marshal(options)
	req, err := v.client.NewRequest(ctx, http.MethodPatch, "/rest/v2/storage/qos/volumes", string(rawdata))
	if err != nil {
		return nil, err
	}

	res := QoSData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
