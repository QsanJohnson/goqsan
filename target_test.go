package goqsan

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestTarget(t *testing.T) {
	fmt.Println("------------TestTarget--------------")

	ctx = context.Background()

	poolId64, _ := strconv.ParseUint(testConf.poolId, 10, 64)

	// create volume paramater
	now := time.Now()
	timeStamp := now.Format("20060102150405")
	volName := "gotest-vol-" + timeStamp
	paramV := VolumeCreateOptions{
		Name:         volName,
		BlockSize:    4096,
		UsedSize:     10240,
		PoolID:       uint64(poolId64),
		IoPriority:   "HIGH",
		BgIoPriority: "HIGH",
		CacheMode:    "WRITE_THROUGH",
	}

	//create volume
	vol, err := testConf.volumeOp.CreateVolume(ctx, &paramV)
	if err != nil {
		t.Fatalf("createVolume failed: %v", err)
	}
	//volID := strconv.Itoa(vol.ID)
	fmt.Printf("  A volume was created. Id:%s \n", vol.ID)

	// create iSCSI parameter
	paramCSCSI := CreateTargetParam{
		Name: "kyle_test_groupSCSI",
		Type: "iSCSI",
		Iscsis: []Iscsi{
			{
				Eths: []string{"c0e1", "c0e2"},
			},
		},
	}

	// create FCP parameter
	paramCFCP := CreateTargetParam{
		Name: "test_groupFCP",
		Type: "FCP",
	}

	// // patch iSCSI parameter
	// paramPSCSI := PatchTargetParam{
	// 	Name: "kyle_goqsm_PSCSI",
	// 	Type: "iSCSI",
	// 	Iscsis: []Iscsi{
	// 		{
	// 			Name: "2",
	// 			Eths: []string{"c0e1", "c0e2"},
	// 		},
	// 	},
	// }

	// // patch FCP parameter
	// paramPFCP := PatchTargetParam{
	// 	Name: "kyle_goqsm_PFCP",
	// }

	// lun mapping parameter
	paramMLun := LunMapParam{
		Name:     "10", // Lun number, can choose from 0 to 254
		VolumeID: "",   //2074967409
		Hosts: []Host{
			{
				Name: []string{"*"}, // iqn/WWN(iSCSI/FCP)
			},
		},
	}

	// lun patch parameter
	paramPLun := LunPatchParam{
		Name: "5", // Lun number, can choose from 0 to 254
		Hosts: []Host{
			{
				Name: []string{"*"}, // iqn/WWN(iSCSI/FCP)
			},
		},
	}

	// listTargetTest(t)
	// createDeleteTargetTest(t, &paramCSCSI)
	// createDeleteTargetTest(t, &paramCFCP)

	// // createTarget, listTarget, patchTarget, deleteTarget
	// createDLPTargetTest(t, &paramCSCSI, &paramPSCSI)
	// createDLPTargetTest(t, &paramCFCP, &paramPFCP)

	// createTarget, mapLun, listLun, patchLun, unmapLun, deleteTarget
	createTargetMapLunTest(t, vol.ID, &paramCSCSI, &paramMLun, &paramPLun)
	createTargetMapLunTest(t, vol.ID, &paramCFCP, &paramMLun, &paramPLun)

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
	tgt, err := testConf.targetOp.CreateTarget(ctx, optionsT)
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
func createDLPTargetTest(t *testing.T, optionsT *CreateTargetParam, optionsP *PatchTargetParam) {
	fmt.Println("createTargetTest Enter")

	//create Target
	tgt, err := testConf.targetOp.CreateTarget(ctx, optionsT)
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
	tgt, err = testConf.targetOp.PatchTarget(ctx, tgt.ID, optionsP)
	if err != nil {
		t.Fatalf("PatchTarget failed: %v", err)
	}
	fmt.Printf("  A Target has been patched. %+v\n", tgt)

	// //delete Target
	// err = testConf.targetOp.DeleteTarget(ctx, tgt.ID)
	// if err != nil {
	// 	t.Fatalf("DeleteTarget failed: %v", err)
	// }
	// fmt.Printf("  A Target was deleted. \n")

	// fmt.Println("createDeleteTargetTest Leave")
}

func createTargetMapLunTest(t *testing.T, volID string, optionsT *CreateTargetParam, optionsL *LunMapParam, optionsP *LunPatchParam) {
	fmt.Println("createTargetMapLunTest Enter")

	//assign volID to lun map parameter
	optionsL.VolumeID = volID

	//create Target
	tgt, err := testConf.targetOp.CreateTarget(ctx, optionsT)
	if err != nil {
		t.Fatalf("CreateTarget failed: %v", err)
	}
	fmt.Printf("  A Target was created. %+v\n", tgt)

	//map Lun
	lunD, err := testConf.targetOp.MapLun(ctx, tgt.ID, optionsL)
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
	lunTgtP, err := testConf.targetOp.PatchTargetLun(ctx, tgt.ID, lunD.ID, optionsP)
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
