package goqsan

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestTarget(t *testing.T) {
	fmt.Println("------------TestTarget--------------")

	ctx = context.Background()

	// create volume paramater
	now := time.Now()
	timeStamp := now.Format("20060102150405")
	volName := "gotest-vol-" + timeStamp
	paramV := VolumeCreateOptions{
		Name:         volName,
		BlockSize:    4096,
		TotalSize:    12288,
		PoolID:       testConf.poolId,
		IoPriority:   "HIGH",
		BgIoPriority: "HIGH",
		CacheMode:    "WRITE_THROUGH",
	}

	//create volume
	vol, err := testConf.volumeOp.CreateVolume(ctx, paramV.PoolID, paramV.Name, paramV.TotalSize, &paramV)
	if err != nil {
		t.Fatalf("createVolume failed: %v", err)
	}
	//volID := strconv.Itoa(vol.ID)
	fmt.Printf("  A volume was created. Id:%s \n", vol.ID)

	// create iSCSI parameter
	paramCSCSI := CreateTargetParam{
		Name: "kyle_test_groupCSCSI",
		Type: "iSCSI",
		Iscsis: []Iscsi{
			{
				Eths: []string{"c0e5", "c0e6"},
			},
		},
	}

	// create FCP parameter
	paramCFCP := CreateTargetParam{
		Name: "kyle_test_groupCFCP",
		Type: "FCP",
	}

	// patch iSCSI parameter
	paramPSCSI := TargetParam{
		Name: "kyle_goqsm_PSCSI",
		Type: "iSCSI",
		Iscsis: []Iscsi{
			{
				Name: "2",
				Eths: []string{"c0e1", "c0e2"},
			},
		},
	}

	// patch FCP parameter
	paramPFCP := TargetParam{
		Name: "kyle_goqsm_PFCP",
		Type: "FCP",
	}

	// lun mapping parameter
	paramMLun := LunMapParam{
		Name:     "10",   // Lun number, can choose from 0 to 254
		VolumeID: vol.ID, //2074967409
		Hosts: []Host{
			{
				Name: "*", // iqn/WWN(iSCSI/FCP)
			},
		},
	}

	// lun patch parameter
	paramPLun := LunParam{
		Name: "5", // Lun number, can choose from 0 to 254
		Hosts: []Host{
			{
				Name: "*", // iqn/WWN(iSCSI/FCP)
			},
		},
	}

	// listTargetTest(t)
	// createDeleteTargetTest(t, &paramCSCSI)
	// createDeleteTargetTest(t, &paramCFCP)

	// createTarget, listTarget, patchTarget, deleteTarget
	createDLPTargetTest(t, &paramCSCSI, &paramPSCSI)
	createDLPTargetTest(t, &paramCFCP, &paramPFCP)

	// createTarget, mapLun, listLun, patchLun, unmapLun, deleteTarget
	createTargetMapLunTest(t, &paramCSCSI, &paramMLun, &paramPLun)
	createTargetMapLunTest(t, &paramCFCP, &paramMLun, &paramPLun)

	//delete volume
	err = testConf.volumeOp.DeleteVolume(ctx, vol.ID)
	if err != nil {
		t.Fatalf("DeleteVolume failed: %v", err)
	}
	fmt.Printf("  A volume was deleted. Id:%s\n", vol.ID)

	//list, patch
	listPatchFCTest(t)
}

func listTargetTest(t *testing.T) {

	tgts, err := testConf.targetOp.ListTargets(ctx, "SCOTT")
	if err != nil {
		t.Fatalf("ListTargets failed: %v", err)
	}
	fmt.Printf("[listTarget] : %+v", tgts)

}

func createDeleteTargetTest(t *testing.T, optionsT *CreateTargetParam) {
	fmt.Println("createTargetTest Enter")

	//create Target
	tgt, err := testConf.targetOp.CreateTarget(ctx, optionsT.Name, optionsT.Type, optionsT)
	if err != nil {
		t.Fatalf("CreateTarget failed: %v", err)
	}
	fmt.Printf("  A Target was created. %+v\n", tgt)

	//delete Target
	err = testConf.targetOp.DeleteTarget(ctx, tgt.ID)
	if err != nil {
		t.Fatalf("DeleteTarget failed: %v", err)
	}
	fmt.Printf("  A Target was deleted. \n")

	fmt.Println("createDeleteTargetTest Leave")
}

// create list patch delete
func createDLPTargetTest(t *testing.T, optionsT *CreateTargetParam, optionsP *TargetParam) {
	fmt.Println("createTargetTest Enter")

	//create Target
	tgt, err := testConf.targetOp.CreateTarget(ctx, optionsT.Name, optionsT.Type, optionsT)
	if err != nil {
		t.Fatalf("CreateTarget failed: %v", err)
	}
	fmt.Printf("  A Target was created. %+v\n", tgt)

	//list target by target ID
	tgt, err = testConf.targetOp.ListTargetByID(ctx, tgt.ID)
	if err != nil {
		t.Fatalf("listTarget by ID failed: %v", err)
	}
	fmt.Printf("  Target ID %v information : %+v\n", tgt.ID, tgt)

	//patch target
	// tgt, err = testConf.targetOp.ModifyTarget(ctx, tgt.ID, optionsP)
	// if err != nil {
	// 	t.Fatalf("ModifyTarget failed: %v", err)
	// }
	// fmt.Printf("  A Target has been patched. %+v\n", tgt)

	fmt.Printf("  Sleep 5 seconds\n")
	time.Sleep(5 * time.Second)

	//delete Target
	err = testConf.targetOp.DeleteTarget(ctx, tgt.ID)
	if err != nil {
		t.Fatalf("DeleteTarget failed: %v", err)
	}
	fmt.Printf("  A Target was deleted. \n")

	fmt.Printf("  Sleep 5 seconds\n")
	time.Sleep(5 * time.Second)

	fmt.Println("createDeleteTargetTest Leave")
}

func createTargetMapLunTest(t *testing.T, optionsT *CreateTargetParam, optionsL *LunMapParam, optionsP *LunParam) {
	fmt.Println("createTargetMapLunTest Enter")

	//create Target
	tgt, err := testConf.targetOp.CreateTarget(ctx, optionsT.Name, optionsT.Type, optionsT)
	if err != nil {
		t.Fatalf("CreateTarget failed: %v", err)
	}
	fmt.Printf("  A Target was created. %+v\n", tgt)

	//map Lun
	lunD, err := testConf.targetOp.MapLun(ctx, tgt.ID, optionsL.VolumeID, optionsL)
	if err != nil {
		t.Fatalf("MapLun failed: %v", err)
	}
	fmt.Printf("  A Lun was mapped. %+v\n", lunD)

	//list target lun
	lunD, err = testConf.targetOp.ListTargetLun(ctx, tgt.ID, lunD.ID)
	if err != nil {
		t.Fatalf("ListLun failed: %v", err)
	}
	fmt.Printf("  listed lun: %+v\n", lunD)

	//Patch target lun
	lunTgtP, err := testConf.targetOp.ModifyTargetLun(ctx, tgt.ID, lunD.ID, optionsP)
	if err != nil {
		t.Fatalf("PatchLun failed: %v", err)
	}
	fmt.Printf("  Patched lun:  %+v\n", lunTgtP)

	//list all luns under given targetID
	lunAllD, err := testConf.targetOp.ListAllLuns(ctx, tgt.ID)
	if err != nil {
		t.Fatalf("ListAllLuns failed: %v", err)
	}
	fmt.Printf("  listed luns:  %+v\n", lunAllD)

	//unmap Lun
	err = testConf.targetOp.UnmapLun(ctx, tgt.ID, lunD.ID)
	// lunID := strconv.Itoa(lunD.ID)
	// err = testConf.targetOp.UnmapLun(ctx, tgt.ID, lunID)
	if err != nil {
		t.Fatalf("UnmapLun failed: %v", err)
	}
	fmt.Printf("  A Lun was unmapped. \n")

	//delete Target
	err = testConf.targetOp.DeleteTarget(ctx, tgt.ID)
	if err != nil {
		t.Fatalf("DeleteTarget failed: %v", err)
	}
	fmt.Printf("  A Target was deleted. \n")

	fmt.Println("createTargetMapLunTest Leave")
}

//TODO patch FC wait for document
//func listPatchFCTest(t *testing.T, optionsP *FCPatchParam) {
func listPatchFCTest(t *testing.T) {
	fmt.Println("listPatchFCTest Enter")

	//list fibre channel
	fc, err := testConf.targetOp.ListFC(ctx)
	if err != nil {
		t.Fatalf("ListFC failed: %v", err)
	}
	fmt.Printf("  listed fibre channel: %+v\n", fc)

	// //patch fibre channel
	// fc, err = testConf.targetOp.PatchFC(ctx, fc.ID, optionsP)
	// if err != nil {
	// 	t.Fatalf("PatchFC failed: %v", err)
	// }
	// fmt.Printf("  patched fibre channel:  %+v\n", fc)

	// fmt.Println("listPatchFCTest Leave")
}
