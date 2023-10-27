package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"
)

var fasta_paths string
var output_dir string
var max_template_date string
var force bool

func init() {
	flag.StringVar(&fasta_paths, "f", "", "Comma separated list of fasta files to be used as input\n example: seq1.fasta,seq2.fasta")
	flag.StringVar(&output_dir, "o", "", "The directory where the results will be stored")
	flag.StringVar(&max_template_date, "max_template_date", "2022-01-01", "The maximum template release date to consider:\n format: YYYY-MM-DD")
}

func main() {
	_ = flag.Set("logtostderr", "true")
	_ = flag.Set("stderrthreshold", "DEBUG")
	_ = flag.Set("v", "2")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:  %s -f seq1.fasta -o /path/to/output\n", os.Args[0])
	}

	flag.Parse()

	if fasta_paths == "" {
		glog.Error("-f argument not set")
		flag.Usage()
		os.Exit(1)
	}

	if output_dir == "" {
		glog.Error("-o argument not set")
		flag.Usage()
		os.Exit(1)
	}

	glog.Info("######################################################")
	glog.Info("AlphaFold wrapper")
	glog.Info("######################################################")

	err := prepareOutputDir(output_dir, force)
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
	glog.Info("Output directory: ", output_dir)

	args := loadEnv(max_template_date, output_dir)

	args.Fasta_paths = fasta_paths
	args.Output_dir = output_dir
	args.Max_template_date = max_template_date

	cmd := args.FormatCmd()

	out, err := sbatch(cmd, args.Partition)
	if err != nil {
		glog.Info("Error running sbatch")
		glog.Error(err)
		glog.Error(out)
		os.Exit(1)
	}
	glog.Info(out)

	os.Exit(0)

}
