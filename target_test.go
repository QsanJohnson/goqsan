package goqsan

import (
	"context"
	"fmt"

	//"strconv"
	"testing"
)

func TestTarget(t *testing.T) {
	fmt.Println("------------TestTarget--------------")

	ctx = context.Background()

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
	// 	Name: "kyle_goqsan_PSCSI",
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
	// 	Name: "kyle_goqsan_PFCP",
	// }

	// lun mapping parameter
	paramMLun := LunMapParam{
		Name:     "10", // Lun number, can choose from 0 to 254
		VolumeID: "2076050576",
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
	createTargetMapLunTest(t, &paramCSCSI, &paramMLun, &paramPLun)
	createTargetMapLunTest(t, &paramCFCP, &paramMLun, &paramPLun)

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

func createDeleteTargetTest(t *testing.T, options *CreateTargetParam) {
	fmt.Println("createTargetTest Enter")

	//create Target
	tgt, err := testConf.targetOp.CreateTarget(ctx, options)
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
func createDLPTargetTest(t *testing.T, optionsC *CreateTargetParam, optionsP *PatchTargetParam) {
	fmt.Println("createTargetTest Enter")

	//create Target
	tgt, err := testConf.targetOp.CreateTarget(ctx, optionsC)
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

func createTargetMapLunTest(t *testing.T, optionsC *CreateTargetParam, optionsL *LunMapParam, optionsP *LunPatchParam) {
	fmt.Println("createTargetMapLunTest Enter")

	//create Target
	tgt, err := testConf.targetOp.CreateTarget(ctx, optionsC)
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
	lunD, err = testConf.targetOp.ListTargetLun(ctx, tgt.ID, lunD.ID, optionsL)
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
	lunAllD, err := testConf.targetOp.ListAllLuns(ctx, tgt.ID, optionsL)
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
