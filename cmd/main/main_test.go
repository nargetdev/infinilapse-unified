package main

import "testing"

func TestSetTheStage(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "toggle AC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetTheStage()
		})
	}
}
