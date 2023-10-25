package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang/glog"
)

func main() {
	var fasta_paths string
	var output_dir string
	var max_template_date string
	var help bool
	var force bool

	flag.StringVar(&fasta_paths, "fasta_paths", "", "Comma separated list of fasta files to be used as input\n example: seq1.fasta,seq2.fasta")
	flag.StringVar(&output_dir, "output_dir", "alphafold-prediction", "The directory where the results will be stored")
	flag.StringVar(&max_template_date, "max_template_date", "2022-01-01", "The maximum template release date to consider:\n format: YYYY-MM-DD")
	flag.BoolVar(&force, "f", false, "Overwrite the output directory if it exists")
	flag.BoolVar(&help, "h", false, "Prints the help message")

	flag.Parse()

	if help || (fasta_paths == "") {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if fasta_paths == "" {
		log.Fatal("fasta_paths argument not set")
	}

	err := prepareOutputDir(output_dir, force)
	if err != nil {
		glog.Fatal(err)
	}

	args := loadEnv(max_template_date, output_dir)

	args.Fasta_paths = fasta_paths
	args.Output_dir = output_dir
	args.Max_template_date = max_template_date

	cmd := args.FormatCmd()

	out, err := remoteRun(cmd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)

}
