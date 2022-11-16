package gqlMgmt

import (
	"testing"
)

func TestIndexToGraph(t *testing.T) {
	type args struct {
		objUrl     string
		bucket     string
		cameraName string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "shloobly test",
			args: args{
				objUrl:     "https://miro.medium.com/max/1400/1*XjVyq2W07x-ZdS4XT1tPjA.jpeg",
				bucket:     "gcb-site-cams",
				cameraName: "wut",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			IndexToGraph(tt.args.objUrl, tt.args.bucket, tt.args.cameraName)
		})
	}
}
