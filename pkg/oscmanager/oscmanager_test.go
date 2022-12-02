package oscmanager

import "testing"

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			OscListen(tt.args.connectionString)
		})
	}
}
