package compiler

import (
	"reflect"
	"testing"
)

func Test_listCameras(t *testing.T) {
	tests := []struct {
		name           string
		inputDir       string
		wantCameraList []string
	}{
		// TODO: Add test cases.
		{
			name:           "simple list",
			inputDir:       "/data/img/dslr/",
			wantCameraList: []string{"6D", "M50"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCameraList := listCameras(tt.inputDir); !reflect.DeepEqual(gotCameraList, tt.wantCameraList) {
				t.Errorf("listCameras() = %v, want %v", gotCameraList, tt.wantCameraList)
			}
		})
	}
}

func Test_yesterdayDateString(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "wut"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := yesterdayDateString(); got != tt.want {
				t.Errorf("yesterdayDateString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: Add test cases.
//{
//	name: "compilePaths",
//	args: args{
//		inputDir:   "/some/path/to/camera",
//		dateOffset: -1,
//	},
//	wantOutFile:          "hrm",
//	wantOutMp4PathString: "/hrrmrmmm/ffff",
//},

func TestIndexChunk(t *testing.T) {
	type args struct {
		photoFilePath string
		cameraName    string
		bucket        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "what you do",
			args: args{
				photoFilePath: "./data/out/camera/2022-11-05.mp4",
				cameraName:    "camera",
				bucket:        "gcb-site-",
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := IndexChunk(tt.args.photoFilePath, tt.args.cameraName, tt.args.bucket); (err != nil) != tt.wantErr {
				t.Errorf("IndexChunk() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPrintCyan(t *testing.T) {
	type args struct {
		say string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "triCyan", args: args{say: "wutwut inda"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintCyan(tt.args.say)
			PrintCyanBold(tt.args.say)
			PrintMagenta("hao")
		})
	}
}

func TestPrintMagenta(t *testing.T) {
	type args struct {
		say string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "blah magenta",
			args: args{
				say: "ohaie",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintMagenta(tt.args.say)
		})
	}
}

func TestChunkCompiler(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "floopachunky compiler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ChunkCompiler()
		})
	}
}

func TestListAvailableDates(t *testing.T) {
	type args struct {
		camDir string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "braveList",
			args: args{
				camDir: "./data/img/dslr/6D",
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListAvailableDates(tt.args.camDir); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListAvailableDates() = %v, want %v", got, tt.want)
			}
		})
	}
}
