package oscmanager

import (
	"testing"
	"time"
)

func TestFadeMaster(t *testing.T) {
	type args struct {
		brightness float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "fade to 90",
			args: args{brightness: 0.99},
		},
		{
			name: "fade to 100",
			args: args{brightness: 1.0},
		},
		{
			name: "fade to zero",
			args: args{brightness: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FadeMaster(tt.args.brightness)
		})
	}
}

func TestOscListen(t *testing.T) {
	type args struct {
		connectionString string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "listen to lx",
			args: args{connectionString: "0.0.0.0:3131"},
		},
		{
			name: "listen to blender",
			args: args{connectionString: "0.0.0.0:9002"},
		},
		{
			name: "listen to drivenByMoss",
			args: args{connectionString: "0.0.0.0:9000"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			OscListen(tt.args.connectionString)
		})
	}
}

func TestStartBlenderDailyCycle(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "start blender cycle",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartBlenderDailyCycle()
		})
	}
}

func Test_daySeconds(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "try it out",
			args: args{t: time.Now()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := daySeconds(tt.args.t); got != tt.want {
				t.Errorf("daySeconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetBlenderFrame(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "set the frame",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetBlenderFrame()
		})
	}
}

func Test_sendOscByPathAndValue(t *testing.T) {
	type args struct {
		path  string
		value interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "sanity",
			args: args{
				path:  "/blender/test",
				value: 27,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//sendOscTrue(tt.args.path)
			sendOscByPathAndValue(tt.args.path, tt.args.value)
		})
	}
}

func Test_sendOscBool(t *testing.T) {
	type args struct {
		path string
		tf   bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "false",
			args: args{
				path: "/blender/wut",
				tf:   false,
			},
		},
		{
			name: "true",
			args: args{
				path: "/blender/offsetframe",
				tf:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendOscBool(tt.args.path, tt.args.tf)
		})
	}
}

func Test_frameJumpString(t *testing.T) {
	type args struct {
		path        string
		deltaString string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "frame jump",
			args: args{
				path:        "/blender/offsetframee",
				deltaString: "delta=27",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frameJumpString(tt.args.path, tt.args.deltaString)
		})
	}
}

func TestSynchronizeBlenderPlayhead(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "groovey",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SynchronizeBlenderPlayhead()
		})
	}
}

func Test_getCurrentFrameOfDay(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "frameiam",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCurrentFrameOfDay(); got != tt.want {
				t.Errorf("getCurrentFrameOfDay() = %v, want %v", got, tt.want)
			}
		})
	}
}
