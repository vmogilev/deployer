package deployer

import (
	"errors"
	"reflect"
	"testing"
)

func TestRunListFile_Parse(t *testing.T) {
	testCases := []struct {
		name string
		f    *RunListFile
		want *RunListFile
		err  error
	}{
		{
			name: "ok1",
			f: &RunListFile{
				Name: "1_d1.json",
			},
			want: &RunListFile{
				Name:      "1_d1.json",
				id:        1,
				directive: "d1",
			},
		},
		{
			name: "invalid-id",
			f: &RunListFile{
				Name: "x_d1.json",
			},
			err: errors.New("doesn't matter"),
		},
		{
			name: "invalid-ext",
			f: &RunListFile{
				Name: "0_d1.txt",
			},
			err: errors.New("doesn't matter"),
		},
		{
			name: "invalid-separator",
			f: &RunListFile{
				Name: "1-d1.json",
			},
			err: errors.New("doesn't matter"),
		},
		{
			name: "missing-ext",
			f: &RunListFile{
				Name: "1_d1",
			},
			err: errors.New("doesn't matter"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.f.Parse()
			if tc.err != nil && err == nil {
				t.Fatalf("RunListFile.Parse() shoud have failed")
			}
			if tc.err == nil && err != nil {
				t.Fatalf("RunListFile.Parse() shoud NOT have failed")
			}
			if tc.err != nil && err != nil {
				return
			}
			if !reflect.DeepEqual(tc.f, tc.want) {
				t.Fatalf("RunListFile.Parse() failed\nwant:%+v\ngot:%+v",
					tc.want,
					tc.f,
				)
			}
		})
	}
}
