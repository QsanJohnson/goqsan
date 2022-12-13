package goqsan

import (
	"context"
	"fmt"
	"testing"
)

func TestPool(t *testing.T) {
	fmt.Println("------------TestPool--------------")
	ctx = context.Background()

	listPoolTest(t)
}

func listPoolTest(t *testing.T) {
	fmt.Println("listPoolTest Enter")

	pools, err := testConf.poolOp.ListPools(ctx)
	if err != nil {
		t.Fatalf("ListPools failed: %v", err)
	}
	fmt.Printf("  ListPools cnt: %d\n", len(*pools))

	if len(*pools) >= 1 {
		vol, err := testConf.poolOp.ListPoolByID(ctx, (*pools)[0].ID)
		if err != nil {
			t.Fatalf("ListPoolByID with first exist ID(%s) failed: %v", (*pools)[0].ID, err)
		}
		fmt.Printf("  ListPoolByID with first exist ID(%s)\n  %+v\n", (*pools)[0].ID, vol)
	}

	fmt.Println("listPoolTest Leave")
}
