package vinfo

import (
	"testing"
)

func TestInfo_JSONString(t *testing.T) {
	type fields struct {
		GitVersion string
		GitCommit  string
		GitBranch  string
		BuildDate  string
		Platform   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Info{
				GitVersion: tt.fields.GitVersion,
				GitCommit:  tt.fields.GitCommit,
				GitBranch:  tt.fields.GitBranch,
				BuildDate:  tt.fields.BuildDate,
				Platform:   tt.fields.Platform,
			}
			got, err := i.JSONString()
			if (err != nil) != tt.wantErr {
				t.Errorf("Info.JSONString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Info.JSONString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInfo_String(t *testing.T) {
	type fields struct {
		GitVersion string
		GitCommit  string
		GitBranch  string
		BuildDate  string
		Platform   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Info{
				GitVersion: tt.fields.GitVersion,
				GitCommit:  tt.fields.GitCommit,
				GitBranch:  tt.fields.GitBranch,
				BuildDate:  tt.fields.BuildDate,
				Platform:   tt.fields.Platform,
			}
			if got := i.String(); got != tt.want {
				t.Errorf("Info.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
