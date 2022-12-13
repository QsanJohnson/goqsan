package goqsan

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestVolume(t *testing.T) {
	fmt.Println("------------TestVolume--------------")
	ctx = context.Background()

	listTest(t)

	//TRUE := true
	FALSE := false

	paramCVol := VolumeCreateOptions{
		BlockSize:       4096,
		IoPriority:      "HIGH",
		BgIoPriority:    "HIGH",
		CacheMode:       "WRITE_THROUGH",
		EnableReadAhead: &FALSE,
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

	//create, list, delete
	now := time.Now()
	timeStamp := now.Format("20060102150405")
	volName := "gotest-vol-" + timeStamp
	createDeleteVolumeTest(t, testConf.poolId, volName, 5120, &paramCVol)

	//create, list, get metadata, patch metadata, get timestamp, list timestamp, delete
	now = time.Now()
	timeStamp = now.Format("20060102150405")
	volName = "gotest-vol-" + timeStamp
	metaDataTest(t, testConf.poolId, volName, 5120, &paramCVol)

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
	snapshotTest(t, testConf.poolId, volName, 10240, &paramCVol)

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

func metaDataTest(t *testing.T, poolID, volname string, volsize uint64, options *VolumeCreateOptions) {
	fmt.Printf("metaDataTest Enter (volname: %s)\n", volname)

	metabyte := []byte{48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64}
	metabyte64 := b64.StdEncoding.EncodeToString(metabyte)

	options.Metadata = VolumeMetadata{
		Status:  "VALID",
		Type:    "CSI Driver",
		Content: metabyte64,
	}
	fmt.Printf("  options: %+v\n", options)

	//create volume
	vol, err := testConf.volumeOp.CreateVolume(ctx, poolID, volname, volsize, options)
	if err != nil {
		t.Fatalf("createVolume failed: %v", err)
	}
	fmt.Printf("  A volume was created. Id: %s, metabyte: %v\n", vol.ID, metabyte)

	_, _, metaContent, err := testConf.volumeOp.GetMetadata(ctx, vol.ID)
	if equal := bytes.HasPrefix(metaContent, metabyte); !equal {
		fmt.Printf("  %v vs %v \n", metabyte, metaContent)
		t.Fatalf("metaDataTest failed. metadata content is not equal.")
	}

	// ASCII 0 ~ 255 test
	buf := make([]byte, 22)
	i := 0
	for ascii := 0; ascii <= 255; ascii++ {
		buf[i] = byte(ascii)
		if i == len(buf)-1 {
			i = 0
			err = testMetaData(vol.ID, buf)
			if err != nil {
				t.Fatalf("metaDataTest failed. err: %v", err)
			}
			buf = make([]byte, 22)
			continue
		}
		i++
	}
	err = testMetaData(vol.ID, buf)
	if err != nil {
		t.Fatalf("metaDataTest failed. err: %v", err)
	}

	tstamp, err := testConf.volumeOp.SetTimestamp(ctx, vol.ID, "AUTO")
	if err != nil {
		t.Fatalf("Update Timestamp failed: %v", err)
	}
	fmt.Printf("  Current timestamp is :%s \n", tstamp)

	// var sleepSec time.Duration = 5
	var sleepSec uint64 = 5
	fmt.Printf("  Sleep %d seconds\n", sleepSec)
	time.Sleep(time.Duration(sleepSec) * time.Second)

	tstamp2, err := testConf.volumeOp.SetTimestamp(ctx, vol.ID, "AUTO")
	if err != nil {
		t.Fatalf("Update Timestamp failed: %v", err)
	}
	fmt.Printf("  Current timestamp is :%s\n", tstamp2)
	t1, _ := strconv.ParseUint(tstamp, 10, 64)
	t2, _ := strconv.ParseUint(tstamp2, 10, 64)
	if (t2 - t1) < sleepSec {
		t.Fatalf("Update timestamp function failed, diff time < %d sec", sleepSec)
	}
	fmt.Printf("  timestamp function OK\n")

	//delete volume
	fmt.Println("start delete")
	err = testConf.volumeOp.DeleteVolume(ctx, vol.ID)
	if err != nil {
		t.Fatalf("DeleteVolume failed: %v", err)
	}
	fmt.Printf("  A volume was deleted. Id:%s\n", vol.ID)

	fmt.Println("metaDataTest Leave")
}

func testMetaData(volId string, buf []byte) error {
	testStatus := "VALID"
	testType := "CSI Driver"

	fmt.Printf("  testMetaData with buff: %v\n", buf)
	metaStatus, metaType, metaContent, err := testConf.volumeOp.SetMetadata(ctx, volId, testStatus, testType, buf)
	if err != nil {
		return fmt.Errorf("testMetaData failed on SetMetadata(%s, %s, %v), err: %v", testStatus, testType, buf, err)
	}
	metaStatus2, metaType2, metaContent2, err := testConf.volumeOp.GetMetadata(ctx, volId)
	if testStatus != metaStatus || metaStatus != metaStatus2 {
		return fmt.Errorf("testMetaData failed. metadata Status is not equal (%s vs %s vs %s)", testStatus, metaStatus, metaStatus2)
	}
	if testType != metaType || metaType != metaType2 {
		return fmt.Errorf("testMetaData failed. metadata Type is not equal (%s vs %s vs %s)", testType, metaType, metaType2)
	}
	if equal, equal2 := bytes.HasPrefix(metaContent, buf), bytes.HasPrefix(metaContent2, buf); !equal || !equal2 {
		fmt.Printf("           buf: %v\n", buf)
		fmt.Printf("   metaContent: %v\n", metaContent)
		fmt.Printf("  metaContent2: %v\n", metaContent2)
		return fmt.Errorf("testMetaData failed. metadata content is not equal.")
	}
	return nil
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

	fmt.Println("Set QoS now.")
	qosdata, err = testConf.volumeOp.SetQoS(ctx, qosEnable, qosRule)
	if err != nil {
		t.Fatalf("SetQoS failed: %v", err)
	}
	fmt.Printf("Set QoS: %+v \n", qosdata)

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

func snapshotTest(t *testing.T, poolID, volname string, volsize uint64, optionsV *VolumeCreateOptions) {
	fmt.Printf("snapshotTest Enter \n")

	//create volume
	vol, err := testConf.volumeOp.CreateVolume(ctx, poolID, volname, volsize, optionsV)
	if err != nil {
		t.Fatalf("createVolume failed: %v", err)
	}
	fmt.Printf("  A volume was created. Id:%s \n", vol.ID)

	fmt.Printf("  Sleep 10 seconds\n")
	time.Sleep(10 * time.Second)

	//get volume snapshot settings
	snapSet, err := testConf.volumeOp.GetSnapshotSetting(ctx, vol.ID)
	if err != nil {
		t.Fatalf("Get volume snapshot setting failed: %v", err)
	}
	fmt.Printf("Volume snapshot settings: %v \n", snapSet)

	//enable snapshot center and assign space
	//patch volume snapshot settings
	optionsSP := &SnapshotMutableSetting{
		TotalSize: 20480,
	}
	snapPat, err := testConf.volumeOp.SetSnapshotSetting(ctx, vol.ID, optionsSP)
	if err != nil {
		t.Fatalf("Enable snapshot center failed: %v", err)
	}
	fmt.Printf("Snapshot center enabled, with space: %d \n", snapPat.TotalSize)

	fmt.Printf("  Sleep 5 seconds\n")
	time.Sleep(5 * time.Second)

	//create volume snapshot
	snapName := "kyle_snap1"
	snapC, err := testConf.volumeOp.CreateSnapshot(ctx, vol.ID, snapName) //need to be changed to body input
	if err != nil {
		t.Fatalf("createVolume snapshot failed: %v", err)
	}
	fmt.Printf("  A volume snapshot was created. Snapshot Id:%s \n", snapC.ID)

	//get volume snapshot lists
	snaplist, err := testConf.volumeOp.ListSnapshots(ctx, vol.ID)
	if err != nil {
		t.Fatalf("Get volume snapshot lists failed: %v", err)
	}
	fmt.Printf("Volume snapshot lists: %v \n", snaplist)

	fmt.Printf("  Sleep 3 seconds\n")
	time.Sleep(3 * time.Second)

	//get certain volume snapshot lists
	snapName = "kyle_snap2"
	snapC2, err := testConf.volumeOp.CreateSnapshot(ctx, vol.ID, snapName) //need to be changed to body input
	if err != nil {
		t.Fatalf("createVolume snapshot failed: %v", err)
	}
	fmt.Printf("  A volume snapshot was created. Snapshot Id:%s \n", snapC2.ID)

	snaplist2, err := testConf.volumeOp.GetSnapshot(ctx, vol.ID, snapC2.ID)
	if err != nil {
		t.Fatalf("Get certain volume snapshot list failed: %v", err)
	}
	fmt.Printf("snapshot ID: %s 's list: %v \n", snapC2.ID, snaplist2)

	fmt.Printf("  Sleep 5 seconds\n")
	time.Sleep(5 * time.Second)

	//Rollback to the first snapshot
	err = testConf.volumeOp.RollbackSnapshot(ctx, vol.ID, snapC.ID)
	if err != nil {
		t.Fatalf("Rollback to first snapshot failed: %v", err)
	}
	//check if rollback to first snapshot will delete the rest snapshot
	snaplist, err = testConf.volumeOp.ListSnapshots(ctx, vol.ID)
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
	err = testConf.volumeOp.DeleteAllSnapshots(ctx, vol.ID)
	if err != nil {
		t.Fatalf("Delete Volume snapshots failed: %v", err)
	}
	fmt.Printf("  All snapshots were deleted. \n")

	//disable snapshot center
	optionsDisable := SnapshotMutableSetting{
		TotalSize: 0,
	}
	snapPat, err = testConf.volumeOp.SetSnapshotSetting(ctx, vol.ID, &optionsDisable)
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
