// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package report

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-configfs-tsm/configfs/configfsi"
	"github.com/google/go-configfs-tsm/configfs/faketsm"
	lt "github.com/google/go-configfs-tsm/configfs/linuxtsm"
)


func GetRawQuote(reportData [64]byte) ([]uint8, error) {
	req := &Request{
		InBlob:     reportData[:],
		GetAuxBlob: false,
	}
	fmt.Printf("getting raw quote through tsm")
	client, err := lt.MakeClient()
	r, err := Create(client, req)
	fmt.Printf("created report")
	
	if err != nil {
		panic("")
	}
	response, err := r.Get()
	
	return response.OutBlob, nil
}

func TestConfigfs(t *testing.T) {
	var arr [64]byte
	for i := range arr {
		arr[i] = 3
	}
	fmt.Printf("in tsm client")
	var err error
	if err != nil {
		panic("")
	}
	resp, err := GetRawQuote(arr)
	
	fmt.Printf("got response %x", resp)
}

func TestGet(t *testing.T) {
	c := &faketsm.Client{Subsystems: map[string]configfsi.Client{"report": faketsm.ReportV7(0)}}
	req := &Request{
		InBlob:     []byte("lessthan64bytesok"),
		GetAuxBlob: true,
	}
	resp, err := Get(c, req)
	if err != nil {
		t.Fatalf("Get(%+v) = %+v, %v, want nil", req, resp, err)
	}
	wantOut := "privlevel: 0\ninblob: 6c6573737468616e363462797465736f6b"
	if !bytes.Equal(resp.OutBlob, []byte(wantOut)) {
		t.Errorf("OutBlob %v is not %v", string(resp.OutBlob), wantOut)
	}
	wantProvider := "fake\n"
	if resp.Provider != wantProvider {
		t.Errorf("provider = %q, want %q", resp.Provider, wantProvider)
	}
	if !bytes.Equal(resp.AuxBlob, []byte(`auxblob`)) {
		t.Errorf("auxblob = %v, want %v", resp.AuxBlob, []byte(`auxblob`))
	}
}

func TestGetErr(t *testing.T) {
	tcs := []struct {
		name    string
		req     *Request
		floor   uint
		wantErr string
	}{
		{
			name: "inblob too big",
			req: &Request{
				InBlob: make([]byte, 4096),
			},
			wantErr: "invalid argument",
		},
		{
			name: "privlevel too high",
			req: &Request{
				InBlob:    make([]byte, 64),
				Privilege: &Privilege{Level: 300},
			},
			wantErr: "privlevel must be 0-3",
		},
		{
			name:    "missing inblob",
			req:     &Request{},
			wantErr: "invalid argument",
		},
		{
			name: "privlevel too low",
			req: &Request{
				InBlob:    make([]byte, 64),
				Privilege: &Privilege{Level: 0},
			},
			floor:   1,
			wantErr: "privlevel 0 cannot be less than 1",
		},
		{
			name: "non-guid",
			req: &Request{
				InBlob:      make([]byte, 64),
				ServiceGuid: "00000000-0000-0000-0000-00000000000g",
			},
			wantErr: "invalid UUID format",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			c := &faketsm.Client{
				Subsystems: map[string]configfsi.Client{"report": faketsm.Report611(tc.floor)}}
			resp, err := Get(c, tc.req)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("Get(%+v) = %+v, %v, want %q", tc.req, resp, err, tc.wantErr)
			}
		})
	}
}
