package parser

import (
	"reflect"
	"testing"
)

func Test_parseGphotoAutoDetect(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		// TODO: Add test cases.
		{
			name:  "mycooltest",
			input: "Canon EOS M50                  usb:001,010",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			println("THIS TEST: " + tt.name + "\n")
			ParseOneLineGphotoAutoDetect(tt.input)
		})
	}
}

func TestParseOneLineGphotoAutoDetect(t *testing.T) {
	type args struct {
		sampleInput string
	}
	tests := []struct {
		name      string
		args      args
		wantModel string
		wantPort  string
	}{
		// TODO: Add test cases.
		{
			name:      "mycooltest",
			args:      args{sampleInput: "Canon EOS M50                  usb:001,010"},
			wantModel: "Canon EOS M50",
			wantPort:  "usb:001,010",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotModel, gotPort := ParseOneLineGphotoAutoDetect(tt.args.sampleInput)
			if gotModel != tt.wantModel {
				t.Errorf("ParseOneLineGphotoAutoDetect() gotModel = %v, want %v", gotModel, tt.wantModel)
			}
			if gotPort != tt.wantPort {
				t.Errorf("ParseOneLineGphotoAutoDetect() gotPort = %v, want %v", gotPort, tt.wantPort)
			}
		})
	}
}

func TestNamesAndPortsFromMultiLineAutoDetect(t *testing.T) {
	type args struct {
		autoDetectMultilineString string
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		// TODO: Add test cases.
		{
			name: "Test multiline input from gphoto --auto-detect",
			args: args{
				autoDetectMultilineString: `
Canon EOS M50                  usb:001,010
Canon EOS 6D                   usb:001,007`,
			},
			want: [][]string{
				[]string{"Canon EOS M50", "usb:001,010"},
				[]string{"Canon EOS 6D", "usb:001,007"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NamesAndPortsFromMultiLineAutoDetect(tt.args.autoDetectMultilineString); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NamesAndPortsFromMultiLineAutoDetect() = %v, want %v", got, tt.want)
			}
		})
	}
}
