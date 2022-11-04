package goqsan

import (
	"context"
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
	createDeleteVolumeTest(t, testConf.poolId, volName, 5120, &paramCVol)

	now = time.Now()
	timeStamp = now.Format("20060102150405")
	volName = "gotest-vol-" + timeStamp
	modifyVolumeTest(t, testConf.poolId, volName, 10240, &paramCVol)

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
