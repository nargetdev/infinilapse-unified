package webcamMgmt

import (
	"reflect"
	"testing"
)

func TestCaptureFromDevicesList(t *testing.T) {
	type args struct {
		devices []string
	}
	tests := []struct {
		name                 string
		args                 args
		wantImgFilePathsList []string
	}{
		// TODO: Add test cases.
		{
			name: "deeply equal",
			args: args{
				devices: []string{
					"/dev/video0", "/dev/video4", "/dev/video8", "/dev/video2", "/dev/video6",
				},
			},
			wantImgFilePathsList: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotImgFilePathsList := CaptureFromDevicesList(tt.args.devices); !reflect.DeepEqual(gotImgFilePathsList, tt.wantImgFilePathsList) {
				t.Errorf("CaptureFromDevicesList() = %v, want %v", gotImgFilePathsList, tt.wantImgFilePathsList)
			}
		})
	}
}
