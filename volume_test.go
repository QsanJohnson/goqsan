package goqsan

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestVolume(t *testing.T) {
	fmt.Println("------------TestVolume--------------")
	ctx = context.Background()

	//listTest(t)

	//TRUE := true
	FALSE := false

	poolId64, _ := strconv.ParseUint(testConf.poolId, 10, 64)

	paramCVol := VolumeCreateOptions{
		BlockSize:       4096,
		IoPriority:      "HIGH",
		BgIoPriority:    "HIGH",
		CacheMode:       "WRITE_THROUGH",
		EnableReadAhead: &FALSE,
	}

	// options2 := VolumeModifyOptions{
	// 	Name:                 "afterMod",
	// 	IoPriority:           "MEDIUM",
	// 	BgIoPriority: "LOW",
	// 	CacheMode:            "WRITE_BACK",
	// 	EnableReadAhead:      &TRUE,
	// }

	now := time.Now()
	timeStamp := now.Format("20060102150405")
	volName := "gotest-vol-" + timeStamp
	createDeleteVolumeTest(t, uint64(poolId64), 5120, volName, &paramCVol)

	now = time.Now()
	timeStamp = now.Format("20060102150405")
	volName = "gotest-vol-" + timeStamp
	modifyVolumeTest(t, uint64(poolId64), 10240, volName, &paramCVol)

}

func listTest(t *testing.T) {

	vols, err := testConf.volumeOp.ListVolumes(ctx, "")
	if err != nil {
		t.Fatalf("ListVolumes failed: %v", err)
	}
	fmt.Printf("[listVolume] : %+v \n", vols)

	volsP, err := testConf.volumeOp.ListVolumesByPoolID(ctx, testConf.poolId)
	if err != nil {
		t.Fatalf("ListVolumes failed: %v", err)
	}
	fmt.Printf("[listVolume] : %+v \n", volsP)
}

func createDeleteVolumeTest(t *testing.T, poolID, volsize uint64, volname string, options *VolumeCreateOptions) {
	fmt.Printf("createDeleteVolumeTest Enter (volSize: %d,  %+v )\n", options.UsedSize, *options)

	options.Name = volname
	options.UsedSize = volsize
	options.PoolID = poolID

	//create volume
	vol, err := testConf.volumeOp.CreateVolume(ctx, options)
	if err != nil {
		t.Fatalf("createVolume failed: %v", err)
	}
	fmt.Printf("  A volume was created. Id:%s \n", vol.ID)

	//list volume
	vols, err := testConf.volumeOp.ListVolumes(ctx, vol.ID)
	if err != nil {
		t.Fatalf("ListVolumes failed: %v", err)
	}
	if len(*vols) != 1 {
		t.Fatalf("Volume %s not found.", vol.ID)
	}
	fmt.Printf("[listVolume] : %+v \n", vols)

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

func modifyVolumeTest(t *testing.T, poolID, volsize uint64, volname string, options *VolumeCreateOptions) {
	fmt.Println("ModifyVolumeTest Enter")

	options.Name = volname
	options.UsedSize = volsize
	options.PoolID = poolID

	// create volume
	vol, err := testConf.volumeOp.CreateVolume(ctx, options)
	if err != nil {
		t.Fatalf("createVolume failed: %v", err)
	}
	//volID := strconv.Itoa(vol.ID)
	fmt.Printf("  A volume was created. Id:%s \n", vol.ID)

	readahead := true
	param := &VolumeModifyOptions{
		UsedSize:        20480,
		CacheMode:       "WRITE_THROUGH",
		EnableReadAhead: &readahead,
	}
	volMod, err := testConf.volumeOp.ModifyVolume(ctx, vol.ID, param)
	if err != nil {
		t.Fatalf("modifyVolume failed: %v", err)
	}
	fmt.Printf("  A volume was modified. %+v \n", volMod)

	// check volume data after mod
	if volMod.UsedSize != param.UsedSize {
		t.Fatalf("modifyVolume change UsedSize failed. \n")
	}
	if volMod.CacheMode != param.CacheMode {
		t.Fatalf("modifyVolume change CacheMode failed. \n")
	}
	if volMod.EnableReadAhead != *param.EnableReadAhead {
		t.Fatalf("modifyVolume change EnableReadAhead failed. \n")
	}

	readahead = false
	param = &VolumeModifyOptions{
		Name:            "afterModRaw131",
		IoPriority:      "MEDIUM",
		BgIoPriority:    "LOW",
		EnableReadAhead: &readahead,
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
