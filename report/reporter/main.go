package main

import (
	"fmt"

	lt "github.com/google/go-configfs-tsm/configfs/linuxtsm"
	report "github.com/google/go-configfs-tsm/report"
)


func TGetRawQuote(reportData [64]byte) ([]uint8, error) {
	req := &report.Request{
		InBlob:     reportData[:],
		GetAuxBlob: false,
	}
	fmt.Printf("getting raw quote through tsm")
	client, err := lt.MakeClient()
	r, err := report.Create(client, req)
	fmt.Printf("created report")
	
	if err != nil {
		panic("")
	}
	response, err := r.Get()
	
	return response.OutBlob, nil
}

func main() {
	var arr [64]byte
	for i := range arr {
		arr[i] = 3
	}
	fmt.Printf("in tsm client")
	var err error
	if err != nil {
		panic("")
	}
	resp, err := TGetRawQuote(arr)
	
	fmt.Printf("got response %x", resp)
}