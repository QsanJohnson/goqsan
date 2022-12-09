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

// type Metadata struct {
// 	Status  string `json:"status"`
// 	Type    string `json:"type"`
// 	Content string `json:"content"`
// }

type VolumeData struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	PoolID                string `json:"poolId"`
	LunID                 string `json:"lunId"`
	TargetID              string `json:"targetId"`
	Online                bool   `json:"online"`
	State                 string `json:"state"`
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
	VDMMetadata struct {
		Status    string `json:"status"`
		Type      string `json:"type"`
		Content   string `json:"content"`
		Timestamp string `json:"timestamp"`
	} `json:"metadata"`
}

type MetadataStruct struct {
	Status    string `json:"status,omitempty"`
	Type      string `json:"type,omitempty"`
	Content   string `json:"content,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
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
	Metadata        MetadataStruct `json:"metadata,omitempty"`
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
	Metadata        MetadataStruct `json:"metadata,omitempty"`
}

// type VolumeModifyOptions struct {
// 	Name            string `json:"name,omitempty"`
// 	TotalSize       uint64 `json:"totalSize,omitempty"`
// 	IoPriority      string `json:"ioPriority,omitempty"`
// 	BgIoPriority    string `json:"bgIoPriority,omitempty"`
// 	CacheMode       string `json:"cacheMode,omitempty"`
// 	EnableReadAhead *bool  `json:"enableReadAhead,omitempty"`
// }

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

// return value of GET /rest/v2/backup/snapshot/targets/_volumeID
// return value of PATCH /rest/v2/backup/snapshot/targets/_volumeID
type VolSnaphshotSetting struct {
	Type              string `json:"type"`
	SnapshotMaxPolicy struct {
		MaxLimit uint64 `json:"maxLimit"`
		Policy   string `json:"policy"`
	} `json:"snapshotMaxPolicy"`
	ProtectionGroup string `json:"protectionGroup"`
	TotalSize       uint64 `json:"totalSize"`
}

// return value of GET /rest/v2/backup/snapshot/targets/_volumeID/snapshots
// return value of Post /rest/v2/backup/snapshot/targets/_volumeID/snapshots
// return value of Patch /rest/v2/backup/snapshot/targets/_volumeID/snapshots/_snapshotID
type VolSnaphshotLists struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CreateTime string `json:"createTime"`
	UsedSize   uint64 `json:"usedSize"`
	Expose     struct {
		Enable    bool        `json:"enable"`
		Mode      interface{} `json:"mode"`
		WriteSize uint64      `json:"writeSize"`
	} `json:"expose"`
	Trash struct {
		InTrash    bool        `json:"inTrash"`
		DeleteTime interface{} `json:"deleteTime"`
	} `json:"trash"`
}

// Post /rest/v2/backup/snapshot/targets/_volumeID/snapshots
type VolumeSnapshotName struct {
	Name string `json:"name,omitempty"`
}

// Patch /rest/v2/backup/snapshot/targets/_volumeID
type VolumeSnapshotPatchSetting struct {
	ProtectionGroup string `json:"protectionGroup,omitempty"`
	TotalSize       int    `json:"totalSize,omitempty"`
}

// Patch /rest/v2/backup/snapshot/targets/_volumeID/snapshots/_snapshotID
type ExposeStruct struct {
	Enable    bool   `json:"enable"`
	Mode      string `json:"mode"`
	WriteSize uint64 `json:"writeSize"`
}

// Patch /rest/v2/backup/snapshot/targets/_volumeID/snapshots/_snapshotID
type TrashStruct struct {
	InTrash bool `json:"inTrash"`
}

// Patch /rest/v2/backup/snapshot/targets/_volumeID/snapshots/_snapshotID
type VolSnaphshotOptions struct {
	Expose ExposeStruct `json:"expose"`
	Trash  TrashStruct  `json:"trash"`
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

// Get Volume snapshot settings
// GET /rest/v2/backup/snapshot/targets/_volumeID
func (v *VolumeOp) GetVolumeSnapshotSetting(ctx context.Context, volId string) (*VolSnaphshotSetting, error) {

	req, err := v.client.NewRequest(ctx, http.MethodGet, "/rest/v2/backup/snapshot/targets/"+volId, nil)
	if err != nil {
		return nil, err
	}

	res := VolSnaphshotSetting{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// Enable snapshot space
// PATCH /rest/v2/backup/snapshot/targets/_volumeID
func (v *VolumeOp) PatchVolumeSnapshotSetting(ctx context.Context, volId string, options *VolumeSnapshotPatchSetting) (*VolSnaphshotSetting, error) {

	rawdata, _ := json.Marshal(options)
	req, err := v.client.NewRequest(ctx, http.MethodPatch, "/rest/v2/backup/snapshot/targets/"+volId, string(rawdata))
	if err != nil {
		return nil, err
	}

	res := VolSnaphshotSetting{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// Create volume snapshot
// POST /rest/v2/backup/snapshot/targets/_volumeID/snapshots/_snapshotName
func (v *VolumeOp) CreateVolumeSnapshotLists(ctx context.Context, volId string, options *VolumeSnapshotName) (*VolSnaphshotLists, error) {

	rawdata, _ := json.Marshal(options)
	req, err := v.client.NewRequest(ctx, http.MethodPost, "/rest/v2/backup/snapshot/targets/"+volId+"/snapshots", string(rawdata))
	if err != nil {
		return nil, err
	}

	res := VolSnaphshotLists{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// Get Volume snapshot lists
// GET /rest/v2/backup/snapshot/targets/_volumeID/snapshots
func (v *VolumeOp) GetVolumeSnapshotLists(ctx context.Context, volId string) (*[]VolSnaphshotLists, error) {

	req, err := v.client.NewRequest(ctx, http.MethodGet, "/rest/v2/backup/snapshot/targets/"+volId+"/snapshots", nil)
	if err != nil {
		return nil, err
	}

	res := []VolSnaphshotLists{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// Get Volume certain snapshot list
// GET /rest/v2/backup/snapshot/targets/_volumeID/snapshots/_snapshotID
func (v *VolumeOp) GetVolumeSnapshotList(ctx context.Context, volId, snapId string) (*VolSnaphshotLists, error) {

	req, err := v.client.NewRequest(ctx, http.MethodGet, "/rest/v2/backup/snapshot/targets/"+volId+"/snapshots/"+snapId, nil)
	if err != nil {
		return nil, err
	}

	res := VolSnaphshotLists{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// Patch certain volume snapshot
// PATCH /rest/v2/backup/snapshot/targets/_volumeID/snapshots/_snapshotID
func (v *VolumeOp) PatchVolumeSnapshot(ctx context.Context, volId, snapId string, options *VolSnaphshotOptions) (*[]VolSnaphshotLists, error) {

	rawdata, _ := json.Marshal(options)
	req, err := v.client.NewRequest(ctx, http.MethodPatch, "/rest/v2/backup/snapshot/targets/"+volId+"/snapshots/"+snapId, string(rawdata))
	if err != nil {
		return nil, err
	}

	res := []VolSnaphshotLists{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Rollback to certain volume snapshot
// POST /rest/v2/backup/snapshot/targets/_volumeID/snapshots/_snapshotID
func (v *VolumeOp) RollbackVolumeSnapshot(ctx context.Context, volId, snapId string) error {
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
func (v *VolumeOp) DeleteVolumeSnapshots(ctx context.Context, volId string) error {
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
func (v *VolumeOp) DeleteVolumeSnapshot(ctx context.Context, volId, snapId string) error {
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
func (v *VolumeOp) GetMetadataTimestamp(ctx context.Context, volId string) (string, error) {

	req, err := v.client.NewRequest(ctx, http.MethodGet, "/rest/v2/storage/block/volumes/"+volId, nil)
	if err != nil {
		return "", err
	}

	res := VolumeData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return "", err
	}
	return res.VDMMetadata.Timestamp, nil
}

// update metadata Timestamp
func (v *VolumeOp) PatchMetadataTimestamp(ctx context.Context, volId, timestamp string) (string, error) {

	param := &VolumeModifyOptions{
		Metadata: MetadataStruct{
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
	return res.VDMMetadata.Timestamp, nil
}

// Get metadata
func (v *VolumeOp) GetMetadata(ctx context.Context, volId string) (string, string, []byte, error) {

	req, err := v.client.NewRequest(ctx, http.MethodGet, "/rest/v2/storage/block/volumes/"+volId, nil)
	if err != nil {
		return "", "", nil, err
	}

	res := VolumeData{}
	if err := v.client.SendRequest(ctx, req, &res); err != nil {
		return "", "", nil, err
	}

	rawDecodedText, _ := b64.StdEncoding.DecodeString(res.VDMMetadata.Content)

	return res.VDMMetadata.Status, res.VDMMetadata.Type, []byte(rawDecodedText), nil
}

// Update metadata
func (v *VolumeOp) PatchMetadata(ctx context.Context, volId, metastatus, metatype string, metacontent []byte) (string, string, []byte, error) {

	metacontent64 := b64.StdEncoding.EncodeToString(metacontent)
	param := &VolumeModifyOptions{
		Metadata: MetadataStruct{
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

	rawDecodedText, _ := b64.StdEncoding.DecodeString(res.VDMMetadata.Content)

	return res.VDMMetadata.Status, res.VDMMetadata.Type, []byte(rawDecodedText), nil
}
