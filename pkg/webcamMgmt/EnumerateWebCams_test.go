package webcamMgmt

import (
	"reflect"
	"testing"
)

func TestDevicesStringListFromListDevices(t *testing.T) {
	type args struct {
		rawString string
	}
	tests := []struct {
		name              string
		args              args
		wantWebcamDevices []string
	}{
		// TODO: Add test cases.
		{
			name: "runme",
			args: args{
				rawString: `
				bcm2835-codec-decode (platform:bcm2835-codec):
				/dev/video10
				/dev/video11
				/dev/video12
				/dev/video18
				/dev/video31
				/dev/media5

				bcm2835-isp (platform:bcm2835-isp):
				/dev/video13
				/dev/video14
				/dev/video15
				/dev/video16
				/dev/video20
				/dev/video21
				/dev/video22
				/dev/video23
				/dev/media6
				/dev/media7

				unicam (platform:fe801000.csi):
				/dev/video17
				/dev/video19
				/dev/media8

				rpivid (platform:rpivid):
				/dev/video24
				/dev/media9

				Web Camera: Web Camera (usb-0000:01:00.0-1.1):
				/dev/video0
				/dev/video1
				/dev/media0

				USB Camera: USB Camera (usb-0000:01:00.0-1.2.1.1):
				/dev/video4
				/dev/video5
				/dev/media2

				USB Camera: USB Camera (usb-0000:01:00.0-1.2.1.3):
				/dev/video8
				/dev/video9
				/dev/media4

				Streaming Webcam: Streaming Web (usb-0000:01:00.0-1.2.2):
				/dev/video2
				/dev/video3
				/dev/media1

				USB 2.0 Camera: USB Camera (usb-0000:01:00.0-1.2.4):
				/dev/video6
				/dev/video7
				/dev/media3
`,
			},
			wantWebcamDevices: []string{
				"/dev/video0", "/dev/video4", "/dev/video8", "/dev/video2", "/dev/video6",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWebcamDevices := DevicesStringListFromListDevices(tt.args.rawString); !reflect.DeepEqual(gotWebcamDevices, tt.wantWebcamDevices) {
				t.Errorf("DevicesStringListFromListDevices() = %v, want %v", gotWebcamDevices, tt.wantWebcamDevices)
			}
		})
	}
}

func TestDevicesStringListFromListDevices1(t *testing.T) {
	type args struct {
		rawString string
	}
	tests := []struct {
		name              string
		args              args
		wantWebcamDevices []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWebcamDevices := DevicesStringListFromListDevices(tt.args.rawString); !reflect.DeepEqual(gotWebcamDevices, tt.wantWebcamDevices) {
				t.Errorf("DevicesStringListFromListDevices() = %v, want %v", gotWebcamDevices, tt.wantWebcamDevices)
			}
		})
	}
}

func TestEnumerateUsbWebCamDevices(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EnumerateUsbWebCamDevices(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnumerateUsbWebCamDevices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSubstring(t *testing.T) {
	type args struct {
		s       string
		indices []int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSubstring(tt.args.s, tt.args.indices); got != tt.want {
				t.Errorf("getSubstring() = %v, want %v", got, tt.want)
			}
		})
	}
}
