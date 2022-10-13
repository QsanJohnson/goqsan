package goqsan

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
	//"time"
)

func TestVolume(t *testing.T) {
	fmt.Println("------------TestVolume--------------")
	ctx = context.Background()

	listTest(t)

	//TRUE := true
	//FALSE := false

	scId64, _ := strconv.ParseUint(testConf.scId, 10, 64)

	options1 := VolumeCreateOptions{
		BlockSize:    4096,
		PoolId:       uint64(scId64),
		IoPriority:   "HIGH",
		BgIoPriority: "HIGH",
		CacheMode:    "WRITE_THROUGH",
	}

	// options2 := VolumeModifyOptions{
	// 	Name:                 "afterMod",
	// 	IoPriority:           "MEDIUM",
	// 	BgIoPriority: "LOW",
	// 	CacheMode:            "WRITE_BACK",
	// 	EnableReadAhead:      &TRUE,
	// }

	// createDeleteVolumeTest(t, 5120, &options1)

	modifyVolumeTest(t, 10240, &options1)

}

func listTest(t *testing.T) {

	vols, err := testConf.volumeOp.ListVolumes(ctx, "")
	if err != nil {
		t.Fatalf("ListVolumes failed: %v", err)
	}
	fmt.Printf("[listVolume] : %+v \n", vols)

	volsP, err := testConf.volumeOp.ListVolumesByPoolID(ctx, testConf.scId)
	if err != nil {
		t.Fatalf("ListVolumes failed: %v", err)
	}
	fmt.Printf("[listVolume] : %+v \n", volsP)
}

func createDeleteVolumeTest(t *testing.T, volSize uint64, options *VolumeCreateOptions) {
	fmt.Printf("createDeleteVolumeTest Enter (volSize: %d,  %+v )\n", volSize, *options)

	now := time.Now()
	timeStamp := now.Format("20060102150405")
	volName := "gotest-vol-" + timeStamp

	//create volume
	vol, err := testConf.volumeOp.CreateVolume(ctx, volName, volSize, options)
	if err != nil {
		t.Fatalf("createVolume failed: %v", err)
	}
	//volId := strconv.Itoa(vol.ID)
	fmt.Printf("  A volume was created. Id:%s \n", vol.ID)

	//list volume
	vols, err := testConf.volumeOp.ListVolumes(ctx, vol.ID)
	if err != nil {
		t.Fatalf("ListVolumes failed: %v", err)
	}
	if len(*vols) != 1 {
		t.Fatalf("Volume %s not found.", vol.ID)
	}
	fmt.Printf("[listVolume] : %+v", vols)

	//delete volume
	fmt.Println("start delete")
	err = testConf.volumeOp.DeleteVolume(ctx, vol.ID)
	if err != nil {
		t.Fatalf("DeleteVolume failed: %v", err)
	}
	fmt.Printf("  A volume was deleted. Id:%s\n", vol.ID)

	fmt.Println("createDeleteVolumeTest Leave")
}

func modifyVolumeTest(t *testing.T, volSize uint64, options *VolumeCreateOptions) {
	fmt.Println("ModifyVolumeTest Enter")
	now := time.Now()
	timeStamp := now.Format("20060102150405")
	volName := "gotest-vol-" + timeStamp

	// create volume
	vol, err := testConf.volumeOp.CreateVolume(ctx, volName, volSize, options)
	if err != nil {
		t.Fatalf("createVolume failed: %v", err)
	}
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

	// delete volume
	err = testConf.volumeOp.DeleteVolume(ctx, vol.ID)
	if err != nil {
		t.Fatalf("DeleteVolume failed: %v", err)
	}
	fmt.Printf("  A volume was deleted. Id:%s\n", vol.ID)

	fmt.Println("ModifyVolumeTest Leave")
}
