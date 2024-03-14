package main

import (
	"testing"

	gogesrtp "github.com/hoilung/go-gesrtp"
)

func TestOpen(t *testing.T) {

	client := gogesrtp.NewGeSrtp("127.0.0.1", 18245, 3000)
	gen := client.Open()
	if !gen {
		defer client.Close()
		t.Fatalf("Open failed")
	}
	val, _ := client.ReadBoolean("M7")
	if !val {
		t.Fatalf("Read bool failed")
	}
	defer client.Close()

}
