package goqsan

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestVolume(t *testing.T) {
	fmt.Println("------------TestVolume--------------")
	ctx = context.Background()

	listTest(t)

	//TRUE := true
	FALSE := false

	//byte test for metadata
	metabyte := []byte{48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64}
	metabyte64 := b64.StdEncoding.EncodeToString(metabyte)
	pmetabyte := []byte{48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 97, 98, 99, 100, 101}

	paramCVol := VolumeCreateOptions{
		BlockSize:       4096,
		IoPriority:      "HIGH",
		BgIoPriority:    "HIGH",
		CacheMode:       "WRITE_THROUGH",
		EnableReadAhead: &FALSE,
		Metadata: MetadataStruct{
			Status:  "VALID",
			Type:    "CSI Driver",
			Content: metabyte64,
		},
	}

	//patch QoS settings
	//IoPriority needs to be "HIGH",in order to make TargetResponseTime value apply to the machine.
	paramPVolQoS := VolumeModifyOptions{
		VolumeQoSOptions: VolumeQoSOptions{
			IoPriority:         "HIGH",
			TargetResponseTime: 123,
			MaxIops:            1234,
			MaxThroughtput:     1234,
		},
	}

	// options2 := VolumeModifyOptions{
	// 	Name:            "afterMod",
	// 	IoPriority:      "MEDIUM",
	// 	BgIoPriority:    "LOW",
	// 	CacheMode:       "WRITE_BACK",
	// 	EnableReadAhead: &TRUE,
	// }

	// create snapshot name
	paramSnapName := VolumeSnapshotName{
		Name: "kyle_snap1",
	}
	// assign snapshot space
	paramSnapSpace := VolumeSnapshotPatchSetting{
		TotalSize: 20480,
	}

	//create, list, delete
	now := time.Now()
	timeStamp := now.Format("20060102150405")
	volName := "gotest-vol-" + timeStamp
	createDeleteVolumeTest(t, testConf.poolId, volName, 5120, &paramCVol)

	//create, list, get metadata, patch metadata, get timestamp, list timestamp, delete
	now = time.Now()
	timeStamp = now.Format("20060102150405")
	volName = "gotest-vol-" + timeStamp
	metaDataTest(t, testConf.poolId, volName, 5120, "1670571795", pmetabyte, &paramCVol)

	//create, patch volume, delete
	now = time.Now()
	timeStamp = now.Format("20060102150405")
	volName = "gotest-vol-" + timeStamp
	modifyVolumeTest(t, testConf.poolId, volName, 10240, &paramCVol)

	modifyQoSTest(t, testConf.poolId, volName, 10240, &paramCVol, &paramPVolQoS)

	qosTest(t, true, "IO_PRIORITY")
	qosTest(t, false, "IO_PRIORITY")
	qosTest(t, true, "MAX_IOPS_THROUGHPUT")

	now = time.Now()
	timeStamp = now.Format("20060102150405")
	volName = "gotest-vol-" + timeStamp
	snapshotTest(t, testConf.poolId, volName, 10240, &paramCVol, &paramSnapSpace, &paramSnapName)

}

func listTest(t *testing.T) {
	fmt.Println("listTest Enter")

	vols, err := testConf.volumeOp.ListVolumes(ctx)
	if err != nil {
		t.Fatalf("ListVolumes failed: %v", err)
	}
	fmt.Printf("ListVolumes: cnt=%d\n%+v \n", len(*vols), vols)

	if len(*vols) >= 1 {
		vol, err := testConf.volumeOp.ListVolumeByID(ctx, (*vols)[0].ID)
		if err != nil {
			t.Fatalf("ListVolumeByID with exist ID failed: %v", err)
		}
		fmt.Printf("ListVolumeByID with exist ID: %+v \n", vol)
	}

	_, err = testConf.volumeOp.ListVolumeByID(ctx, "2222222222")
	if err != nil {
		resterr, ok := err.(*RestError)
		if ok {
			fmt.Printf("ListVolumeByID with non-exist ID, StatusCode=%d ErrResp=%+v\n", resterr.StatusCode, resterr.ErrResp)
			if resterr.StatusCode != http.StatusNotFound {
				t.Fatalf("ListVolumeByID with non-exist ID failed: StatusCode=%d ErrResp=%+v\n", resterr.StatusCode, resterr.ErrResp)
			}
		} else {
			t.Fatalf("ListVolumeByID with non-exist ID failed: %v\n", resterr)
		}
	}
	fmt.Printf("ListVolumeByID with non-exist ID PASS\n")

	vols, err = testConf.volumeOp.ListVolumesByPoolID(ctx, testConf.poolId)
	if err != nil {
		t.Fatalf("ListVolumesByPoolID failed: %v", err)
	}
	fmt.Printf("ListVolumesByPoolID: cnt=%d\n%+v \n", len(*vols), vols)

	fmt.Println("listTest Leave")
}

func createDeleteVolumeTest(t *testing.T, poolID, volname string, volsize uint64, options *VolumeCreateOptions) {
	fmt.Printf("createDeleteVolumeTest Enter (volSize: %d,  %+v )\n", options.TotalSize, *options)

	//create volume
	vol, err := testConf.volumeOp.CreateVolume(ctx, poolID, volname, volsize, options)
	if err != nil {
		t.Fatalf("createVolume failed: %v", err)
	}
	fmt.Printf("  A volume was created. Id:%s \n", vol.ID)

	//list volume
	v, err := testConf.volumeOp.ListVolumeByID(ctx, vol.ID)
	if err != nil {
		resterr, ok := err.(*RestError)
		if ok {
			if resterr.StatusCode == http.StatusNotFound {
				t.Fatalf("Volume %s not found.", vol.ID)
			}
			fmt.Printf("[ListVolumeByID] StatusCode=%d ErrResp=%+v\n", resterr.StatusCode, resterr.ErrResp)
		}

		t.Fatalf("ListVolumeByID failed: %v", err)
	}
	fmt.Printf("ListVolumeByID : %+v \n", v)

	fmt.Printf("  Sleep 5 seconds\n")
	time.Sleep(5 * time.Second)

	//delete volume
	fmt.Println("start delete")
	err = testConf.volumeOp.DeleteVolume(ctx, vol.ID)
	if err != nil {
		t.Fatalf("DeleteVolume failed: %v", err)
	}
	fmt.Printf("  A volume was deleted. Id:%s\n", vol.ID)

	fmt.Println("createDeleteVolumeTest Leave")
}

func metaDataTest(t *testing.T, poolID, volname string, volsize uint64, ptimestamp string, pMetabyte []byte, options *VolumeCreateOptions) {
	fmt.Printf("createDeleteVolumeTest Enter (volSize: %d,  %+v )\n", options.TotalSize, *options)

	//create volume
	vol, err := testConf.volumeOp.CreateVolume(ctx, poolID, volname, volsize, options)
	if err != nil {
		t.Fatalf("createVolume failed: %v", err)
	}
	fmt.Printf("  A volume was created. Id:%s \n", vol.ID)

	//list volume
	v, err := testConf.volumeOp.ListVolumeByID(ctx, vol.ID)
	if err != nil {
		resterr, ok := err.(*RestError)
		if ok {
			if resterr.StatusCode == http.StatusNotFound {
				t.Fatalf("Volume %s not found.", vol.ID)
			}
			fmt.Printf("[ListVolumeByID] StatusCode=%d ErrResp=%+v\n", resterr.StatusCode, resterr.ErrResp)
		}

		t.Fatalf("ListVolumeByID failed: %v", err)
	}
	fmt.Printf("ListVolumeByID : %+v \n", v)

	//Get metadata
	metaStatus, metaType, metaContent, err := testConf.volumeOp.GetMetadata(ctx, vol.ID)
	if err != nil {
		t.Fatalf("Get Metadata failed: %v", err)
	}
	fmt.Printf(" metadata Status is :%s \n", metaStatus)
	fmt.Printf(" metadata Type is :%s \n", metaType)
	fmt.Printf(" metadata Content is :%b \n", metaContent)

	metaStatus, metaType, metaContent, err = testConf.volumeOp.PatchMetadata(ctx, vol.ID, metaStatus, metaType, pMetabyte)
	if err != nil {
		t.Fatalf("Update Metadata failed: %v", err)
	}
	fmt.Printf(" New metadata Status is :%s \n", metaStatus)
	fmt.Printf(" New metadata Type is :%s \n", metaType)
	fmt.Printf(" New metadata Content is :%b \n", metaContent)

	//Get metadata timestamp
	tstamp, err := testConf.volumeOp.GetMetadataTimestamp(ctx, vol.ID)
	if err != nil {
		t.Fatalf("Get Timestamp failed: %v", err)
	}
	fmt.Printf(" metadata timestamp is :%s \n", tstamp)

	//update metadata timestamp
	tstamp, err = testConf.volumeOp.PatchMetadataTimestamp(ctx, vol.ID, ptimestamp)
	if err != nil {
		t.Fatalf("Update Timestamp failed: %v", err)
	}
	fmt.Printf("Updated metadata timestamp is :%s \n", tstamp)

	fmt.Printf("  Sleep 5 seconds\n")
	time.Sleep(5 * time.Second)

	//delete volume
	fmt.Println("start delete")
	err = testConf.volumeOp.DeleteVolume(ctx, vol.ID)
	if err != nil {
		t.Fatalf("DeleteVolume failed: %v", err)
	}
	fmt.Printf("  A volume was deleted. Id:%s\n", vol.ID)

	fmt.Println("createDeleteVolumeTest Leave")
}

func modifyVolumeTest(t *testing.T, poolID, volname string, volsize uint64, options *VolumeCreateOptions) {
	fmt.Println("ModifyVolumeTest Enter")

	// create volume
	vol, err := testConf.volumeOp.CreateVolume(ctx, poolID, volname, volsize, options)
	if err != nil {
		t.Fatalf("createVolume failed: %v", err)
	}
	//volID := strconv.Itoa(vol.ID)
	fmt.Printf("  A volume was created. Id:%s \n", vol.ID)

	readahead := true
	param := &VolumeModifyOptions{
		TotalSize:       20480,
		CacheMode:       "WRITE_THROUGH",
		EnableReadAhead: &readahead,
	}
	volMod, err := testConf.volumeOp.ModifyVolume(ctx, vol.ID, param)
	if err != nil {
		t.Fatalf("modifyVolume failed: %v", err)
	}
	fmt.Printf("  A volume was modified. %+v \n", volMod)

	// check volume data after mod
	if volMod.TotalSize != param.TotalSize {
		t.Fatalf("modifyVolume change TotalSize failed. \n")
	}
	if volMod.CacheMode != param.CacheMode {
		t.Fatalf("modifyVolume change CacheMode failed. \n")
	}
	if volMod.EnableReadAhead != *param.EnableReadAhead {
		t.Fatalf("modifyVolume change EnableReadAhead failed. \n")
	}

	readahead = false
	param = &VolumeModifyOptions{
		Name:             "afterModRaw131",
		VolumeQoSOptions: VolumeQoSOptions{IoPriority: "MEDIUM"},
		BgIoPriority:     "LOW",
		EnableReadAhead:  &readahead,
	}
	volMod, err = testConf.volumeOp.ModifyVolume(ctx, vol.ID, param)
	if err != nil {
		t.Fatalf("modifyVolume failed: %v", err)
	}
	fmt.Printf("  A volume was modified. %+v \n", volMod)

	// check volume data after mod
	if volMod.Name != param.Name {
		t.Fatalf("modifyVolume change Name failed. \n")
	}
	if volMod.IoPriority != param.IoPriority {
		t.Fatalf("modifyVolume change IoPriority failed. \n")
	}
	if volMod.BgIoPriority != param.BgIoPriority {
		t.Fatalf("modifyVolume change BgIoPriority failed. \n")
	}

	if volMod.EnableReadAhead != *param.EnableReadAhead {
		t.Fatalf("modifyVolume change EnableReadAhead failed. \n")
	}

	fmt.Printf("  Sleep 5 seconds\n")
	time.Sleep(5 * time.Second)

	//delete volume
	err = testConf.volumeOp.DeleteVolume(ctx, vol.ID)
	if err != nil {
		t.Fatalf("DeleteVolume failed: %v", err)
	}
	fmt.Printf("  A volume was deleted. Id:%s\n", vol.ID)

	fmt.Println("ModifyVolumeTest Leave")
}

func modifyQoSTest(t *testing.T, poolID, volname string, volsize uint64, optionsV *VolumeCreateOptions, optionsQ *VolumeModifyOptions) {
	fmt.Println("ModifyQoSTest Enter")

	// create volume
	vol, err := testConf.volumeOp.CreateVolume(ctx, poolID, volname, volsize, optionsV)
	if err != nil {
		t.Fatalf("createVolume failed: %v", err)
	}
	fmt.Printf("  A volume was created. Id:%s \n", vol.ID)

	volMod, err := testConf.volumeOp.ModifyVolume(ctx, vol.ID, optionsQ)
	if err != nil {
		t.Fatalf("modifyQoS failed: %v", err)
	}
	fmt.Printf("  A volume's QoS was modified. %+v \n", volMod)

	// check volume QoS after mod
	if volMod.IoPriority != optionsQ.IoPriority {
		t.Fatalf("modifyQoS change IoPriority failed. \n")
	}
	if volMod.TargetResponseTime != optionsQ.TargetResponseTime {
		t.Fatalf("modifyQoS change TargetResponseTime failed. \n")
	}
	if volMod.MaxIops != optionsQ.MaxIops {
		t.Fatalf("modifyQoS change MaxIops failed. \n")
	}
	if volMod.MaxThroughtput != optionsQ.MaxThroughtput {
		t.Fatalf("modifyQoS change MaxThroughtput failed. \n")
	}

	fmt.Printf("  Sleep 5 seconds\n")
	time.Sleep(5 * time.Second)

	//delete volume
	err = testConf.volumeOp.DeleteVolume(ctx, vol.ID)
	if err != nil {
		t.Fatalf("DeleteVolume failed: %v", err)
	}
	fmt.Printf("  A volume was deleted. Id:%s\n", vol.ID)

	fmt.Println("ModifyQoSTest Leave")
}

func qosTest(t *testing.T, qosEnable bool, qosRule string) {
	fmt.Println("QoSTest Enter")

	fmt.Println("Get QoS now.")
	qosdata, err := testConf.volumeOp.GetQoS(ctx)
	if err != nil {
		t.Fatalf("Getqos failed: %v", err)
	}
	fmt.Printf("QoS: %+v \n", qosdata)

	fmt.Println("Patch QoS now.")
	qosdata, err = testConf.volumeOp.PatchQoS(ctx, qosEnable, qosRule)
	if err != nil {
		t.Fatalf("PatchQoS failed: %v", err)
	}
	fmt.Printf("Patched QoS: %+v \n", qosdata)

	//check if Patch QoS is working
	if qosEnable == false {
		if qosdata.EnableQos != qosEnable {
			t.Fatalf("Patch enableQos failed! Input is %t,and QoS enableQos is %t. \n", qosEnable, qosdata.EnableQos)
		}
		if qosdata.QosRule != "NONE" {
			t.Fatalf("Patch QosRule failed.QoS qosRule should be \"NONE\",but it return %s.\n", qosdata.QosRule)
		}
	} else if qosEnable == true {
		if qosdata.EnableQos != qosEnable {
			t.Fatalf("Patch enableQos failed! Input is %t,and QoS enableQos is %t. \n", qosEnable, qosdata.EnableQos)
		}
		if qosdata.QosRule != qosRule {
			t.Fatalf("Patch enableQos failed! Input is %s,and QoS qosRule is %s. \n", qosRule, qosdata.QosRule)
		}
	}

	fmt.Println("QoSTest Leave")
}

func snapshotTest(t *testing.T, poolID, volname string, volsize uint64, optionsV *VolumeCreateOptions, optionsSP *VolumeSnapshotPatchSetting, optionsSN *VolumeSnapshotName) {
	fmt.Printf("snapshotTest Enter \n")
	fmt.Printf("createVolume (volSize: %d,  %+v )\n", optionsV.TotalSize, *optionsV)

	//create volume
	vol, err := testConf.volumeOp.CreateVolume(ctx, poolID, volname, volsize, optionsV)
	if err != nil {
		t.Fatalf("createVolume failed: %v", err)
	}
	fmt.Printf("  A volume was created. Id:%s \n", vol.ID)

	fmt.Printf("  Sleep 10 seconds\n")
	time.Sleep(10 * time.Second)

	//get volume snapshot settings
	snapSet, err := testConf.volumeOp.GetVolumeSnapshotSetting(ctx, vol.ID)
	if err != nil {
		t.Fatalf("Get volume snapshot setting failed: %v", err)
	}
	fmt.Printf("Volume snapshot settings: %v \n", snapSet)

	//enable snapshot center and assign space
	//patch volume snapshot settings
	snapPat, err := testConf.volumeOp.PatchVolumeSnapshotSetting(ctx, vol.ID, optionsSP)
	if err != nil {
		t.Fatalf("Enable snapshot center failed: %v", err)
	}
	fmt.Printf("Snapshot center enabled, with space: %d \n", snapPat.TotalSize)

	fmt.Printf("  Sleep 5 seconds\n")
	time.Sleep(5 * time.Second)

	//create volume snapshot
	snapC, err := testConf.volumeOp.CreateVolumeSnapshotLists(ctx, vol.ID, optionsSN) //need to be changed to body input
	if err != nil {
		t.Fatalf("createVolume snapshot failed: %v", err)
	}
	fmt.Printf("  A volume snapshot was created. Snapshot Id:%s \n", snapC.ID)

	//get volume snapshot lists
	snaplist, err := testConf.volumeOp.GetVolumeSnapshotLists(ctx, vol.ID)
	if err != nil {
		t.Fatalf("Get volume snapshot lists failed: %v", err)
	}
	fmt.Printf("Volume snapshot lists: %v \n", snaplist)

	fmt.Printf("  Sleep 3 seconds\n")
	time.Sleep(3 * time.Second)

	//get certain volume snapshot lists
	snapName2 := VolumeSnapshotName{
		Name: "kyle_snap2",
	}
	snapC2, err := testConf.volumeOp.CreateVolumeSnapshotLists(ctx, vol.ID, &snapName2) //need to be changed to body input
	if err != nil {
		t.Fatalf("createVolume snapshot failed: %v", err)
	}
	fmt.Printf("  A volume snapshot was created. Snapshot Id:%s \n", snapC2.ID)

	snaplist2, err := testConf.volumeOp.GetVolumeSnapshotList(ctx, vol.ID, snapC2.ID)
	if err != nil {
		t.Fatalf("Get certain volume snapshot list failed: %v", err)
	}
	fmt.Printf("snapshot ID: %s 's list: %v \n", snapC2.ID, snaplist2)

	fmt.Printf("  Sleep 5 seconds\n")
	time.Sleep(5 * time.Second)

	//Rollback to the first snapshot
	err = testConf.volumeOp.RollbackVolumeSnapshot(ctx, vol.ID, snapC.ID)
	if err != nil {
		t.Fatalf("Rollback to first snapshot failed: %v", err)
	}
	//check if rollback to first snapshot will delete the rest snapshot
	snaplist, err = testConf.volumeOp.GetVolumeSnapshotLists(ctx, vol.ID)
	if err != nil {
		t.Fatalf("Get volume snapshot lists failed: %v", err)
	}
	if len(*snaplist) != 1 {
		t.Fatalf("Rollback to first snapshot applied, but there are still more than one snapshot left.")
	}
	fmt.Printf("Rollback to the first snapshot, ID: %s . \n", snapC.ID)

	fmt.Printf("  Sleep 5 seconds\n")
	time.Sleep(5 * time.Second)

	//delete snapshot
	fmt.Println("start delete volume snapshot")
	err = testConf.volumeOp.DeleteVolumeSnapshots(ctx, vol.ID)
	if err != nil {
		t.Fatalf("Delete Volume snapshots failed: %v", err)
	}
	fmt.Printf("  All snapshots were deleted. \n")

	//disable snapshot center
	optionsDisable := VolumeSnapshotPatchSetting{
		TotalSize: 0,
	}
	snapPat, err = testConf.volumeOp.PatchVolumeSnapshotSetting(ctx, vol.ID, &optionsDisable)
	if err != nil {
		t.Fatalf("Disable snapshot center failed: %v", err)
	}
	fmt.Printf("Snapshot center disabled, space now is: %d \n", snapPat.TotalSize)

	fmt.Printf("  Sleep 5 seconds\n")
	time.Sleep(5 * time.Second)

	//delete volume
	fmt.Println("start delete volume")
	err = testConf.volumeOp.DeleteVolume(ctx, vol.ID)
	if err != nil {
		t.Fatalf("DeleteVolume failed: %v", err)
	}
	fmt.Printf("  A volume was deleted. Id:%s\n", vol.ID)

	fmt.Println("snapshotTest Leave")
}
