package goqsan

import (
	"context"
	"fmt"
	"testing"
)

var ctx context.Context

func TestSystem(t *testing.T) {
	fmt.Println("------------TestSystem--------------")

	ctx = context.Background()

	getAboutTest(t)
}

func getAboutTest(t *testing.T) {
	fmt.Println("getAboutTest Enter")

	_, err := testConf.systemOp.GetAbout(ctx)
	if err != nil {
		t.Fatalf("getAbout failed: %v", err)
	}

	fmt.Println("getAboutTest Leave")
}
