package gcpMgmt

import (
	"testing"
)

func TestUploadFile(t *testing.T) {
	type args struct {
		filePath string
		bucket   string
		object   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "good path",
			args: args{
				filePath: "./getToDaBucket/clobber.jpg",
				bucket:   "infinilapse-sousveillant-0",
				object:   "unittest/clobber.jpg",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UploadFile(tt.args.filePath, tt.args.object, tt.args.bucket); (err != nil) != tt.wantErr {
				t.Errorf("UploadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStoreFileToBucket(t *testing.T) {
	type args struct {
		fullPath string
		objPath  string
		bucket   string
	}
	tests := []struct {
		name       string
		args       args
		wantObjUrl string
		wantErr    bool
	}{
		{
			name: "good path",
			args: args{
				fullPath: "./getToDaBucket/clobber.jpg",
				bucket:   "gcb-site-pub",
				objPath:  "clobber.jpg",
			},
			wantErr:    false,
			wantObjUrl: "https://storage.googleapis.com/gcb-site-pub/clobber.jpg",
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotObjUrl, err := StoreFileToBucket(tt.args.fullPath, tt.args.objPath, tt.args.bucket)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreFileToBucket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotObjUrl != tt.wantObjUrl {
				t.Errorf("StoreFileToBucket() gotObjUrl = %v, want %v", gotObjUrl, tt.wantObjUrl)
			}
		})
	}
}
