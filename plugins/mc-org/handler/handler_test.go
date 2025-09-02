package handler

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestOrgHandler_HandleCobraCommand(t *testing.T) {
	type fields struct {
		bearerToken string
	}
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// {
		// 	name:    "No file provided",
		// 	fields:  fields{""},
		// 	args:    args{cmd: &cobra.Command{}, args: []string{}},
		// 	wantErr: true,
		// },
		// {
		// 	name:    "File not found",
		// 	fields:  fields{""},
		// 	args:    args{cmd: &cobra.Command{}, args: []string{"nonexistentfile.yaml"}},
		// 	wantErr: true,
		// },
		// {
		// 	name:    "Invalid file format",
		// 	fields:  fields{""},
		// 	args:    args{cmd: &cobra.Command{}, args: []string{"testdata/invalid-format.txt"}},
		// 	wantErr: true,
		// },
		{
			name:    "Valid JSON file",
			fields:  fields{""},
			args:    args{cmd: &cobra.Command{}, args: []string{"testdata/org.json"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrgHandler{
				bearerToken: tt.fields.bearerToken,
			}
			if err := o.HandleCobraCommand(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("OrgHandler.HandleCobraCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
