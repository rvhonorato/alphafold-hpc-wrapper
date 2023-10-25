package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
}

var (
	GPU_DEVICE  string
	INSTALL_DIR string
)

func init() {
	gpu := os.Getenv("GPU_DEVICE")
	if gpu == "" {
		GPU_DEVICE = "gpu001"
	} else {
		GPU_DEVICE = gpu
	}
	installationDir := os.Getenv("INSTALL_DIR")
	if installationDir == "" {
		INSTALL_DIR = "/trinity/login/rodrigo/repos/alphafold-wrapper"
	} else {
		INSTALL_DIR = installationDir
	}
}

// Load the environment variables into the arguments
func loadEnv(maxDate, outputDir string) AFArguments {

	dataDir := os.Getenv("DATA_DIR")
	// if dataDir == "" {
	// 	log.Fatal("DATA_DIR environment variable not set")
	// }

	args := AFArguments{}
	args.Fasta_paths = ""
	args.Max_template_date = maxDate
	args.Data_dir = dataDir
	args.Output_dir = outputDir
	args.Uniref90_database_path = filepath.Join(dataDir, "uniref90/uniref90.fasta")
	args.Mgnify_database_path = filepath.Join(dataDir, "mgnify/mgy_clusters_2022_05.fa")
	args.Template_mmcif_dir = filepath.Join(dataDir, "pdb_mmcif/mmcif_files")
	args.Bfd_database_path = filepath.Join(dataDir, "bfd/bfd_metaclust_clu_complete_id30_c90_final_seq.sorted_opt")
	args.Uniref30_database_path = filepath.Join(dataDir, "uniref30/UniRef30_2021_03")
	args.Pdb70_database_path = filepath.Join(dataDir, "pdb70/pdb70")
	args.Obsolete_pdbs_path = filepath.Join(dataDir, "pdb_mmcif/obsolete.dat")
	args.Use_gpu_relax = true

	return args

}

// FormatCmd formats the arguments into the proper command line arguments
func (args *AFArguments) FormatCmd() string {
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

	condaCMD := "source " + filepath.Join(INSTALL_DIR, "/miniconda3/etc/profile.d/conda.sh")
	condaCMD += " && conda activate af2"

	return condaCMD + " && " + afCmd

}

// prepareOutputDir prepares the output directory, if it exists and the force flag is not set, it will exit
func prepareOutputDir(output_dir string, force bool) error {
	_, err := os.Stat(output_dir)
	if !os.IsNotExist(err) && !force {
		return errors.New("output directory `" + output_dir + "` exists, run with -f to overwrite")
	} else if !os.IsNotExist(err) && force {
		os.RemoveAll(output_dir)

	}
	os.MkdirAll(output_dir, 0755)

	return nil

}

// remoteRun runs the command in the remote machine
func remoteRun(c string) (string, error) {
	// fmt.Println("ssh " + GPU_DEVICE + " \"" + c + "\"")
	out, err := exec.Command("ssh", GPU_DEVICE, "\""+c+"\"").CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		return string(out), err
	}
	return string(out), nil

}
