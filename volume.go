// @2022 QSAN Inc. All rights reserved

package goqsan

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"net/http"
)

// VolumeOp handles volume related methods of the QSAN storage.
type VolumeOp struct {
	client *AuthClient
}

type VolumeMetadata struct {
	Status    string `json:"status,omitempty"`
	Type      string `json:"type,omitempty"`
	Content   string `json:"content,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

type VolumeData struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	PoolID                string `json:"poolId"`
	LunID                 string `json:"lunId"`
	TargetID              string `json:"targetId"`
	Online                bool   `json:"online"`
	State                 string `json:"state"`
	Progress              int    `json:"progress"`
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
	TargetResponseTime    uint64 `json:"targetResponseTime"`
	MaxIops               uint64 `json:"maxIops"`
	MaxThroughtput        uint64 `json:"maxThroughtput"`
	Tags                  struct {
		Wwn  string `json:"wwn"`
		Type string `json:"type"`
	} `json:"tags"`
	Metadata VolumeMetadata `json:"metadata"`
}

type VolumeCreateOptions struct {
	Name            string         `json:"name"`
	TotalSize       uint64         `json:"totalSize"`
	BlockSize       uint64         `json:"blockSize"`
	PoolID          string         `json:"poolId"`
	IoPriority      string         `json:"ioPriority,omitempty"`
	BgIoPriority    string         `json:"bgIoPriority,omitempty"`
	CacheMode       string         `json:"cacheMode,omitempty"`
	EnableReadAhead *bool          `json:"enableReadAhead,omitempty"`
	Metadata        VolumeMetadata `json:"metadata,omitempty"`
}

//Patch /rest/v2/storage/block/volumes/_volumes
type Tag struct {
	Type string `json:"type,omitempty"`
}

type VolumeQoSOptions struct {
	IoPriority         string `json:"ioPriority,omitempty"`
	TargetResponseTime uint64 `json:"targetResponseTime,omitempty"`
	MaxIops            uint64 `json:"maxIops,omitempty"`
	MaxThroughtput     uint64 `json:"maxThroughtput,omitempty"`
}

//Patch /rest/v2/storage/block/volumes/_volumes
type VolumeModifyOptions struct {
	VolumeQoSOptions
	Name            string         `json:"name,omitempty"`
	TotalSize       uint64         `json:"totalSize,omitempty"`
	BgIoPriority    string         `json:"bgIoPriority,omitempty"`
	CacheMode       string         `json:"cacheMode,omitempty"`
	EnableReadAhead *bool          `json:"enableReadAhead,omitempty"`
	Tags            Tag            `json:"tags,omitempty"`
	Metadata        VolumeMetadata `json:"metadata,omitempty"`
}

//return value of GET /rest/v2/storage/qos/volumes
//Patch /rest/v2/storage/qos/volumes
type QoSData struct {
	EnableQos bool   `json:"enableQos"`
	QosRule   string `json:"qosRule"`
}

// Patch /rest/v2/backup/snapshot/targets/_volumeID
type SnapshotMutableSetting struct {
	ProtectionGroup string `json:"protectionGroup,omitempty"`
	TotalSize       int    `json:"totalSize,omitempty"`
}

// return value of GET /rest/v2/backup/snapshot/targets/_volumeID
// return value of PATCH /rest/v2/backup/snapshot/targets/_volumeID
type SnaphshotSetting struct {
	Type              string `json:"type"`
	SnapshotMaxPolicy struct {
		MaxLimit uint64 `json:"maxLimit"`
		Policy   string `json:"policy"`
	} `json:"snapshotMaxPolicy"`
	SnapshotMutableSetting
	AvailableSize int `json:"availableSize"`
	MinimumSize   int `json:"minimumSize"`
}

// Patch /rest/v2/backup/snapshot/targets/_volumeID/snapshots/_snapshotID
type SnapExpose struct {
	Enable    bool   `json:"enable"`
	Mode      string `json:"mode"`
	WriteSize uint64 `json:"writeSize"`
}

// Patch /rest/v2/backup/snapshot/targets/_volumeID/snapshots/_snapshotID
type SnapTrash struct {
	InTrash bool `json:"inTrash"`
}

// Patch /rest/v2/backup/snapshot/targets/_volumeID/snapshots/_snapshotID
type SnaphshotOptions struct {
	Expose SnapExpose `json:"expose"`
	Trash  SnapTrash  `json:"trash"`
}

// return value of GET /rest/v2/backup/snapshot/targets/_volumeID/snapshots
// return value of Post /rest/v2/backup/snapshot/targets/_volumeID/snapshots
// return value of Patch /rest/v2/backup/snapshot/targets/_volumeID/snapshots/_snapshotID
type SnaphshotData struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	CreateTime int64      `json:"createTime"`
	UsedSize   uint64     `json:"usedSize"`
	Expose     SnapExpose `json:"expose"`
	Trash      SnapTrash  `json:"trash"`
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
func (v *VolumeOp) ListVolumesByPoolID(ctx context.Context, poolId string) (*[]VolumeData, error) {

	req, err := v.client.NewRequest(ctx, http.MethodGet, "/rest/v2/storage/block/volumes?q=poolId='"+poolId+"'", nil)
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
func (v *VolumeOp) CreateVolume(ctx context.Context, poolId, volname string, volsize uint64, options *VolumeCreateOptions) (*VolumeData, error) {

	options.PoolID = poolId
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

func (v *VolumeOp) SetQoS(ctx context.Context, qosEnable bool, qosRule string) (*QoSData, error) {

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

// Get Volume snapshot settings
// GET /rest/v2/backup/snapshot/targets/_volumeID
func (v *VolumeOp) GetSnapshotSetting(ctx context.Context, volId string) (*SnaphshotSetting, error) {

	req, err := v.client.NewRequest(ctx, http.MethodGet, "/rest/v2/backup/snapshot/targets/"+volId, nil)
	if err != nil {
		return nil, err
	}

	res := SnaphshotSetting{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// Enable snapshot space
// PATCH /rest/v2/backup/snapshot/targets/_volumeID
func (v *VolumeOp) SetSnapshotSetting(ctx context.Context, volId string, options *SnapshotMutableSetting) (*SnaphshotSetting, error) {

	rawdata, _ := json.Marshal(options)
	req, err := v.client.NewRequest(ctx, http.MethodPatch, "/rest/v2/backup/snapshot/targets/"+volId, string(rawdata))
	if err != nil {
		return nil, err
	}

	res := SnaphshotSetting{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// Create volume snapshot
// POST /rest/v2/backup/snapshot/targets/_volumeID/snapshots/_snapshotName
func (v *VolumeOp) CreateSnapshot(ctx context.Context, volId, snapeName string) (*SnaphshotData, error) {

	m := map[string]string{
		"name": snapeName,
	}
	rawdata, _ := json.Marshal(m)
	req, err := v.client.NewRequest(ctx, http.MethodPost, "/rest/v2/backup/snapshot/targets/"+volId+"/snapshots", string(rawdata))
	if err != nil {
		return nil, err
	}

	res := SnaphshotData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// List all volume snapshots
// GET /rest/v2/backup/snapshot/targets/_volumeID/snapshots
func (v *VolumeOp) ListSnapshots(ctx context.Context, volId string) (*[]SnaphshotData, error) {

	req, err := v.client.NewRequest(ctx, http.MethodGet, "/rest/v2/backup/snapshot/targets/"+volId+"/snapshots", nil)
	if err != nil {
		return nil, err
	}

	res := []SnaphshotData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// Get Volume certain snapshot list
// GET /rest/v2/backup/snapshot/targets/_volumeID/snapshots/_snapshotID
func (v *VolumeOp) GetSnapshot(ctx context.Context, volId, snapId string) (*SnaphshotData, error) {

	req, err := v.client.NewRequest(ctx, http.MethodGet, "/rest/v2/backup/snapshot/targets/"+volId+"/snapshots/"+snapId, nil)
	if err != nil {
		return nil, err
	}

	res := SnaphshotData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// Patch certain volume snapshot
// PATCH /rest/v2/backup/snapshot/targets/_volumeID/snapshots/_snapshotID
func (v *VolumeOp) ModifySnapshot(ctx context.Context, volId, snapId string, options *SnaphshotOptions) (*[]SnaphshotData, error) {

	rawdata, _ := json.Marshal(options)
	req, err := v.client.NewRequest(ctx, http.MethodPatch, "/rest/v2/backup/snapshot/targets/"+volId+"/snapshots/"+snapId, string(rawdata))
	if err != nil {
		return nil, err
	}

	res := []SnaphshotData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Rollback to certain volume snapshot
// POST /rest/v2/backup/snapshot/targets/_volumeID/snapshots/_snapshotID
func (v *VolumeOp) RollbackSnapshot(ctx context.Context, volId, snapId string) error {
	req, err := v.client.NewRequest(ctx, http.MethodPost, "/rest/v2/backup/snapshot/targets/"+volId+"/snapshots/"+snapId+"/rollback", nil)
	if err != nil {
		return err
	}

	res := EmptyData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return err
	}

	return nil
}

// Delete all volume snapshots
// DELETE /rest/v2/backup/snapshot/targets/_volumeID/snapshots
func (v *VolumeOp) DeleteAllSnapshots(ctx context.Context, volId string) error {
	req, err := v.client.NewRequest(ctx, http.MethodDelete, "/rest/v2/backup/snapshot/targets/"+volId+"/snapshots", nil)
	if err != nil {
		return err
	}

	res := EmptyData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return err
	}

	return nil
}

// Delete certain volume snapshot
// DELETE /rest/v2/backup/snapshot/targets/_volumeID/snapshots/_snapshotID
func (v *VolumeOp) DeleteSnapshot(ctx context.Context, volId, snapId string) error {
	req, err := v.client.NewRequest(ctx, http.MethodDelete, "/rest/v2/backup/snapshot/targets/"+volId+"/snapshots/"+snapId, nil)
	if err != nil {
		return err
	}

	res := EmptyData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return err
	}

	return nil
}

// Get metadata Timestamp
func (v *VolumeOp) GetTimestamp(ctx context.Context, volId string) (string, error) {

	req, err := v.client.NewRequest(ctx, http.MethodGet, "/rest/v2/storage/block/volumes/"+volId, nil)
	if err != nil {
		return "", err
	}

	res := VolumeData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return "", err
	}
	return res.Metadata.Timestamp, nil
}

// update metadata Timestamp
func (v *VolumeOp) SetTimestamp(ctx context.Context, volId, timestamp string) (string, error) {

	param := &VolumeModifyOptions{
		Metadata: VolumeMetadata{
			Timestamp: timestamp,
		},
	}
	rawdata, _ := json.Marshal(param)

	req, err := v.client.NewRequest(ctx, http.MethodPatch, "/rest/v2/storage/block/volumes/"+volId, string(rawdata))
	if err != nil {
		return "", err
	}

	res := VolumeData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return "", err
	}
	return res.Metadata.Timestamp, nil
}

// Get metadata
func (v *VolumeOp) GetMetadata(ctx context.Context, volId string) (metastatus, metatype string, metacontent []byte, err error) {

	req, err := v.client.NewRequest(ctx, http.MethodGet, "/rest/v2/storage/block/volumes/"+volId, nil)
	if err != nil {
		return "", "", nil, err
	}

	res := VolumeData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return "", "", nil, err
	}

	rawDecodedText, _ := b64.StdEncoding.DecodeString(res.Metadata.Content)

	return res.Metadata.Status, res.Metadata.Type, []byte(rawDecodedText), nil
}

// Update metadata
func (v *VolumeOp) SetMetadata(ctx context.Context, volId, metastatus, metatype string, metacontent []byte) (string, string, []byte, error) {

	metacontent64 := b64.StdEncoding.EncodeToString(metacontent)
	param := &VolumeModifyOptions{
		Metadata: VolumeMetadata{
			Status:  metastatus,
			Type:    metatype,
			Content: metacontent64,
		},
	}
	rawdata, _ := json.Marshal(param)

	req, err := v.client.NewRequest(ctx, http.MethodPatch, "/rest/v2/storage/block/volumes/"+volId, string(rawdata))
	if err != nil {
		return "", "", nil, err
	}

	res := VolumeData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return "", "", nil, err
	}

	rawDecodedText, _ := b64.StdEncoding.DecodeString(res.Metadata.Content)

	return res.Metadata.Status, res.Metadata.Type, []byte(rawDecodedText), nil
}

func (v *VolumeOp) Clone(ctx context.Context, volId, newVolName, poolId string) (*VolumeData, error) {
	m := map[string]string{
		"volumeName": newVolName,
		"poolID":     poolId,
	}
	rawdata, _ := json.Marshal(m)
	req, err := v.client.NewRequest(ctx, http.MethodPost, "/rest/v2/storage/block/volumes/"+volId+"/clone", string(rawdata))
	if err != nil {
		return nil, err
	}

	res := VolumeData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
