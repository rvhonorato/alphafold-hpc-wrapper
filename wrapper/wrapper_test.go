package main

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestLoadEnv(t *testing.T) {
	maxDate := "2022-01-01"
	outputDir := "/path/to/output"
	dataDir := "/path/to/data"
	os.Setenv("DATA_DIR", dataDir)
	args := loadEnv(maxDate, outputDir)

	if args.Max_template_date != maxDate {
		t.Errorf("AFArguments.Max_template_date = %q, want %q", args.Max_template_date, maxDate)
	}

	if args.Output_dir != outputDir {
		t.Errorf("AFArguments.Output_dir = %q, want %q", args.Output_dir, outputDir)
	}

	if args.Data_dir != dataDir {
		t.Errorf("AFArguments.Data_dir = %q, want %q", args.Data_dir, dataDir)
	}

	if args.Uniref90_database_path != "/path/to/data/uniref90/uniref90.fasta" {
		t.Errorf("AFArguments.Uniref90_database_path = %q, want %q", args.Uniref90_database_path, "/path/to/data/uniref90/uniref90.fasta")
	}

	if args.Mgnify_database_path != "/path/to/data/mgnify/mgy_clusters_2022_05.fa" {
		t.Errorf("AFArguments.Mgnify_database_path = %q, want %q", args.Mgnify_database_path, "/path/to/data/mgnify/mgy_clusters_2022_05.fa")
	}

	if args.Template_mmcif_dir != "/path/to/data/pdb_mmcif/mmcif_files" {
		t.Errorf("AFArguments.Template_mmcif_dir = %q, want %q", args.Template_mmcif_dir, "/path/to/data/pdb_mmcif/mmcif_files")
	}

	if args.Bfd_database_path != "/path/to/data/bfd/bfd_metaclust_clu_complete_id30_c90_final_seq.sorted_opt" {
		t.Errorf("AFArguments.Bfd_database_path = %q, want %q", args.Bfd_database_path, "/path/to/data/bfd/bfd_metaclust_clu_complete_id30_c90_final_seq.sorted_opt")
	}

	if args.Uniref30_database_path != "/path/to/data/uniref30/UniRef30_2021_03" {
		t.Errorf("AFArguments.Uniref30_database_path = %q, want %q", args.Uniref30_database_path, "/path/to/data/uniref30/uniref30.fasta")
	}

	if args.Pdb70_database_path != "/path/to/data/pdb70/pdb70" {
		t.Errorf("AFArguments.Pdb70_database_path = %q, want %q", args.Pdb70_database_path, "/path/to/data/pdb70/pdb70")
	}

	if args.Obsolete_pdbs_path != "/path/to/data/pdb_mmcif/obsolete.dat" {
		t.Errorf("AFArguments.Obsolete_pdbs_path = %q, want %q", args.Obsolete_pdbs_path, "/path/to/data/pdb70/obsolete.dat")
	}

}

func TestPrepareOutputDir(t *testing.T) {
	existingDir := "existing-dir"
	os.MkdirAll(existingDir, 0755)
	defer os.RemoveAll(existingDir)
	nonExistingDir := uuid.New().String()
	defer os.RemoveAll(nonExistingDir)

	type args struct {
		outputDir string
		force     bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test prepare output dir",
			args: args{
				outputDir: nonExistingDir,
				force:     false,
			},
		},
		{
			name: "Test prepare output dir with force",
			args: args{
				outputDir: existingDir,
				force:     true,
			},
		},
		{
			name: "Test prepare output for existing dir without force",
			args: args{
				outputDir: existingDir,
				force:     false,
			},
			wantErr: true,
		},
		{
			name: "Test prepare output for existing dir with force",
			args: args{
				outputDir: existingDir,
				force:     true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputDir := tt.args.outputDir
			force := tt.args.force
			if err := prepareOutputDir(outputDir, force); (err != nil) != tt.wantErr {
				t.Errorf("prepareOutputDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})

	}

}

func TestAFArguments_FormatCmd(t *testing.T) {
	args := &AFArguments{
		Fasta_paths:            "/path/to/fasta",
		Max_template_date:      "2022-01-01",
		Data_dir:               "/path/to/data",
		Output_dir:             "/path/to/output",
		Uniref90_database_path: "/path/to/uniref90",
		Mgnify_database_path:   "/path/to/mgnify",
		Template_mmcif_dir:     "/path/to/template",
		Bfd_database_path:      "/path/to/bfd",
		Uniref30_database_path: "/path/to/uniref30",
		Pdb70_database_path:    "/path/to/pdb70",
		Obsolete_pdbs_path:     "/path/to/obsolete_pdbs",
		Use_gpu_relax:          true,
	}

	expectedCmd := "source /trinity/login/rodrigo/repos/alphafold-wrapper/miniconda3/etc/profile.d/conda.sh && conda activate af2 && python /trinity/login/rodrigo/repos/alphafold-wrapper/alphafold/run_alphafold.py --fasta_paths=/path/to/fasta --max_template_date=2022-01-01 --data_dir=/path/to/data --output_dir=/path/to/output --uniref90_database_path=/path/to/uniref90 --mgnify_database_path=/path/to/mgnify --template_mmcif_dir=/path/to/template --bfd_database_path=/path/to/bfd --uniref30_database_path=/path/to/uniref30 --pdb70_database_path=/path/to/pdb70 --obsolete_pdbs_path=/path/to/obsolete_pdbs --use_gpu_relax=true"

	if gotCmd := args.FormatCmd(); gotCmd != expectedCmd {
		t.Errorf("AFArguments.FormatCmd() = %q, want %q", gotCmd, expectedCmd)
	}
}
