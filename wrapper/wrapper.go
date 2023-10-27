package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type AFArguments struct {
	Fasta_paths            string
	Max_template_date      string
	Data_dir               string
	Output_dir             string
	Uniref90_database_path string
	Mgnify_database_path   string
	Template_mmcif_dir     string
	Bfd_database_path      string
	Uniref30_database_path string
	Pdb70_database_path    string
	Obsolete_pdbs_path     string
	Use_gpu_relax          bool
	Install_dir            string
	Partition              string
}

var (
	PARTITION   string
	INSTALL_DIR string
	DATA_DIR    string
)

func init() {
	PARTITION = os.Getenv("PARTITION")
	if PARTITION == "" {
		PARTITION = "gpu"
	}
	INSTALL_DIR = os.Getenv("INSTALL_DIR")
	if INSTALL_DIR == "" {
		INSTALL_DIR = "/trinity/login/rodrigo/repos/alphafold-wrapper"
	}
	DATA_DIR = os.Getenv("DATA_DIR")
	if DATA_DIR == "" {
		DATA_DIR = "/trinity/login/rodrigo/repos/alphafold-wrapper/data"
	}
}

// Load the environment variables into the arguments
func loadEnv(maxDate, outputDir string) AFArguments {

	args := AFArguments{}
	args.Fasta_paths = ""
	args.Max_template_date = maxDate
	args.Data_dir = DATA_DIR
	args.Partition = PARTITION
	args.Install_dir = INSTALL_DIR
	args.Output_dir = outputDir
	args.Uniref90_database_path = filepath.Join(DATA_DIR, "uniref90/uniref90.fasta")
	args.Mgnify_database_path = filepath.Join(DATA_DIR, "mgnify/mgy_clusters_2022_05.fa")
	args.Template_mmcif_dir = filepath.Join(DATA_DIR, "pdb_mmcif/mmcif_files")
	args.Bfd_database_path = filepath.Join(DATA_DIR, "bfd/bfd_metaclust_clu_complete_id30_c90_final_seq.sorted_opt")
	args.Uniref30_database_path = filepath.Join(DATA_DIR, "uniref30/UniRef30_2021_03")
	args.Pdb70_database_path = filepath.Join(DATA_DIR, "pdb70/pdb70")
	args.Obsolete_pdbs_path = filepath.Join(DATA_DIR, "pdb_mmcif/obsolete.dat")
	args.Use_gpu_relax = true

	return args

}

// FormatCmd formats the arguments into the proper command line arguments
func (args *AFArguments) FormatCmd() string {

	wd, _ := os.Getwd()
	cdCmd := "cd " + wd

	afCmd := "python " + filepath.Join(INSTALL_DIR, "alphafold/run_alphafold.py")
	afCmd += " --fasta_paths=" + args.Fasta_paths
	afCmd += " --max_template_date=" + args.Max_template_date
	afCmd += " --data_dir=" + args.Data_dir
	afCmd += " --output_dir=" + args.Output_dir
	afCmd += " --uniref90_database_path=" + args.Uniref90_database_path
	afCmd += " --mgnify_database_path=" + args.Mgnify_database_path
	afCmd += " --template_mmcif_dir=" + args.Template_mmcif_dir
	afCmd += " --bfd_database_path=" + args.Bfd_database_path
	afCmd += " --uniref30_database_path=" + args.Uniref30_database_path
	afCmd += " --pdb70_database_path=" + args.Pdb70_database_path
	afCmd += " --obsolete_pdbs_path=" + args.Obsolete_pdbs_path
	afCmd += " --use_gpu_relax=" + fmt.Sprintf("%t", args.Use_gpu_relax)

	condaCMD := "source " + filepath.Join(INSTALL_DIR, "/miniconda3/etc/profile.d/conda.sh") + "\n"
	condaCMD += "conda activate af2"

	return condaCMD + "\n\n" + cdCmd + "\n\n" + afCmd

}

// prepareOutputDir prepares the output directory, if it exists and the force flag is not set, it will exit
func prepareOutputDir(output_dir string, force bool) error {
	_, err := os.Stat(output_dir)
	if !os.IsNotExist(err) && !force {
		return errors.New("output directory `" + output_dir + "` exists, erase it or define a new one")
	} else if !os.IsNotExist(err) && force {
		os.RemoveAll(output_dir)

	}
	os.MkdirAll(output_dir, 0755)

	return nil

}

// prepareJobFile prepares the job file
func prepareJobFile(c, partition string) string {

	header := "#!/bin/bash\n"
	header += "#SBATCH --job-name=alphafold\n"
	header += "#SBATCH --nodes=1\n"
	header += "#SBATCH --ntasks-per-node=1\n"
	header += "#SBATCH --cpus-per-task=1\n"
	// header += "#SBATCH --mem=0\n"
	// header += "#SBATCH --time=24:00:00\n"
	header += "#SBATCH --partition=" + partition + "\n"
	header += "#SBATCH --gres=gpu:1\n"
	header += "#SBATCH --output=alphafold-%j.out\n"
	header += "#SBATCH --error=alphafold-%j.err\n"

	body := c

	return header + body

}

func sbatch(c, partition string) (string, error) {
	jobFile := prepareJobFile(c, partition)
	cmd := exec.Command("sbatch")
	cmd.Stdin = strings.NewReader(jobFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}
