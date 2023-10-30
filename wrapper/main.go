package main

import (
	"flag"
	"fmt"
	"haddocking/alphafold-wrapper/overlay"
	"os"

	"github.com/golang/glog"
)

var fasta_paths string
var output_dir string
var max_template_date string
var preset string

func init() {
	flag.StringVar(&fasta_paths, "f", "", "Comma separated list of fasta files to be used as input\n example: seq1.fasta,seq2.fasta")
	flag.StringVar(&output_dir, "o", "", "The directory where the results will be stored")
	flag.StringVar(&max_template_date, "max_template_date", "2022-01-01", "The maximum template release date to consider:\n format: YYYY-MM-DD")
	flag.StringVar(&preset, "p", "", "<monomer|monomer_casp14|monomer_ptm|multimer>")
}

func main() {
	_ = flag.Set("logtostderr", "true")
	_ = flag.Set("stderrthreshold", "DEBUG")
	_ = flag.Set("v", "2")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:  %s -f target.fasta -p <monomer|monomer_casp14|monomer_ptm|multimer> -o /path/to/output\n\n### To fold a multimer, the target.fasta must be a multi-fasta file\n### Ex: `https://www.rcsb.org/fasta/entry/2OOB/display`\n\n", os.Args[0])
	}

	flag.Parse()

	if fasta_paths == "" {
		glog.Error("-f argument not set")
		flag.Usage()
		os.Exit(1)
	}

	_, err := os.Stat(fasta_paths)
	if os.IsNotExist(err) {
		glog.Error("File ", fasta_paths, " does not exist")
		os.Exit(1)
	}

	if output_dir == "" {
		glog.Error("-o argument not set")
		flag.Usage()
		os.Exit(1)
	}

	if preset == "" {
		glog.Error("-preset argument not set")
		flag.Usage()
		os.Exit(1)
	}

	// Check if preset is one of the allowed values
	allowed_presets := []string{"monomer", "monomer_casp14", "monomer_ptm", "multimer"}
	found := false
	for _, p := range allowed_presets {
		if p == preset {
			found = true
			break
		}
	}
	if !found {
		glog.Error("Invalid preset value: ", preset)
		flag.Usage()
		os.Exit(1)
	}

	glog.Info("######################################################")
	glog.Info("AlphaFold wrapper")
	glog.Info("######################################################")

	err = overlay.PrepareOutputDir(output_dir)
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
	glog.Info("Output directory: ", output_dir)

	args := overlay.LoadEnv(max_template_date, output_dir)

	args.Fasta_paths = fasta_paths
	args.Output_dir = output_dir
	args.Max_template_date = max_template_date
	args.Preset = preset

	cmd := args.FormatCmd()

	out, err := overlay.RunCommand(cmd, args.Partition, "sbatch")
	if err != nil {
		glog.Info("Error running sbatch")
		glog.Error(err)
		glog.Error(out)
		os.Exit(1)
	}
	glog.Info(out)

	os.Exit(0)

}
