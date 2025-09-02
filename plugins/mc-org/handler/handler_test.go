package handler

import (
	"testing"

	"github.com/Mattilsynet/map-cli/internal/config"
	"github.com/spf13/cobra"
)

func TestOrgHandler_HandleCobraCommand(t *testing.T) {
	if t.Skipped() {
		t.Skip("Skipping test as it requires a valid bearer token and external dependencies.")
	}
	bearerToken := config.CurrentConfig.Zitadel.BearerToken
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
			fields:  fields{bearerToken: bearerToken},
			args:    args{cmd: &cobra.Command{}, args: []string{"testdata/org.json"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrgHandler{
				bearerToken: tt.fields.bearerToken,
			}
			err := o.HandleCobraCommand(tt.args.cmd, tt.args.args)
			if err != nil {
				t.Errorf("OrgHandler.HandleCobraCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
