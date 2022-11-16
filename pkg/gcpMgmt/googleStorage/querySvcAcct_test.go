package googleStorage

import (
	"bytes"
	"google.golang.org/api/iam/v1"
	"reflect"
	"testing"
)

func Test_listServiceAccounts(t *testing.T) {
	type args struct {
		projectID string
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		want    []*iam.ServiceAccount
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "example from google",
			args: args{
				projectID: "gcb-cube-base",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got, err := listServiceAccounts(w, tt.args.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("listServiceAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("listServiceAccounts() gotW = %v, want %v", gotW, tt.wantW)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listServiceAccounts() got = %v, want %v", got, tt.want)
			}
		})
	}
}
