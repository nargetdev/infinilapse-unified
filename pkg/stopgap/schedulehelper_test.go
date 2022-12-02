package stopgap

import "testing"

func TestSTOPGAP_getFadeByTime(t *testing.T) {
	tests := []struct {
		name string
		want float64
	}{
		{
			name: "day check",
			want: 1.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := STOPGAP_getFadeByTime(); got != tt.want {
				t.Errorf("STOPGAP_getFadeByTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
