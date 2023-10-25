# AlphaFold-wrapper for HPC systems

This repository contains a set of instructions on how-to-install alphafold on a HPC cluster. It is based on the [alphafold repository](https://github.com/google-deepmind/alphafold).

Both the `alphafold` and `aria2c` repositories are linked here as submodules.

## Clone this repository

```bash
# Clone the repository
$ git clone https://github.com/rvhonorato/alphafold-wrapper.git && cd alphafold-wrapper

# Download the submodules
$ git submodule update --init --recursive

# Set a variable to the current directory to be used later
$ export ALPHAFOLD_WRAPPER=$(pwd)
```

## aria2c

`aria2c` is used to download the alphafold database; either install it with the package manager or  build it from source.

```bash
$ cd $ALPHAFOLD_WRAPPER/aria2/
$ autoreconf -i
$ ./configure ARIA2_STATIC=yes
$ make

# Add the binary to the path
$ export PATH=$PATH:$ALPHAFOLD_WRAPPER/aria2/src
```

## Python environment

Alphafold need a python environment to be executed, so we will use `miniconda` to create one (as suggested by the alphafold repository).

```bash
# install miniconda3
$ cd $ALPHAFOLD_WRAPPER
$ aria2c https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh
$ chmod +x ./Miniconda3-latest-Linux-x86_64.sh
$ ./Miniconda3-latest-Linux-x86_64.sh -b -p $ALPHAFOLD_WRAPPER/miniconda3
$ rm Miniconda3-latest-Linux-x86_64.sh

# create an environment to be used later
$ $ALPHAFOLD_WRAPPER/miniconda3/bin/conda install -y -n af2 python=3.10
```

## Alphafold

### Download data

Note that the uncompressed data size is ~550 GB so send it to the background with `nohup`

```bash
$ cd $ALPHAFOLD_WRAPPER
$ mkdir -p $ALPHAFOLD_WRAPPER/data
$ nohup bash $ALPHAFOLD_WRAPPER/alphafold/scripts/download_all_data.sh $ALPHAFOLD_WRAPPER/data &

# inspect the log if needed
$ tail nohup.out
```

### Install

Find out what is your CUDA version and set the `CUDA_VERSION` variable

```bash
$ nvidia-smi | grep "CUDA Version"
| NVIDIA-SMI 510.47.03    Driver Version: 510.47.03    CUDA Version: 11.6     |
$ export CUDA_VERSION=11.6
```

```bash
$ source $ALPHAFOLD_WRAPPER/miniconda3/etc/profile.d/conda.sh

# (Optional) Configure the base env with libmamba for faster dependency solver
$ conda update -n base conda
$ conda install -n base conda-libmamba-solver
$ conda config --set solver libmamba

# Create and activate the alphafold environment
$ conda create -y -n af2 python=3.10 && conda activate af2

# Install conda packages
(af2) $ conda install -y -c conda-forge openmm=7.7.0 cudatoolkit==${CUDA_VERSION} pdbfixer pip
(af2) $ conda install -y -c bioconda kalign3
(af2) $ conda install -y -c conda-forge -c bioconda hhsuite==3.3.0
(af2) $ conda install -y -c bioconda hmmer==3.3.2

# Install pip packages
(af2) $ pip install --upgrade pip
(af2) $ pip install -r $ALPHAFOLD_WRAPPER/alphafold/requirements.txt

## mind the cuda version here!
(af2) $ pip install jax==0.3.25 jaxlib==0.3.25+cuda11.cudnn805 \
  -f https://storage.googleapis.com/jax-releases/jax_cuda_releases.html
```

## Make a prediction

```bash
# Define the location of the databases
export DATA_DIR=$ALPHAFOLD_WRAPPER/data
export UNIREF90_DATABASE_PATH=$DATA_DIR/uniref90/uniref90.fasta
export MGNIFY_DATABASE_PATH=$DATA_DIR/mgnify/mgy_clusters_2022_05.fa
export BFD_DATABASE_PATH=$DATA_DIR/bfd/bfd_metaclust_clu_complete_id30_c90_final_seq.sorted_opt
export UNIREF50_DATABASE_PATH=$DATA_DIR/uniref50/uniref50.fasta
export PDB70_DATABASE_PATH=$DATA_DIR/pdb70/pdb70
export TEMPLATE_MMCIF_PATH=$DATA_DIR/pdb_mmcif/mmcif_files
export BFD_DATABASE_PATH=$DATA_DIR/bfd/bfd_metaclust_clu_complete_id30_c90_final_seq.sorted_opt
export UNIREF30_DATABASE_PATH=$DATA_DIR/uniref30/UniRef30_2021_03
export OBSOLETE_PDB_PATH=$DATA_DIR/pdb_mmcif/obsolete.dat

# Define the Input sequence
export INPUT_FASTA=$ALPHAFOLD_WRAPPER/example_data/1crn.fasta
export OUTPUT_DIR=$ALPHAFOLD_WRAPPER/example_data
export MAX_TEMPLATE_DATE=2022-01-01

(af2) $ python $ALPHAFOLD_WRAPPER/alphafold/run_alphafold.py \
  --fasta_paths=$INPUT_FASTA \
  --max_template_date=$MAX_TEMPLATE_DATE \
  --data_dir=$DATA_DIR \
  --output_dir=$OUTPUT_DIR \
  --uniref90_database_path=$UNIREF90_DATABASE_PATH \
  --mgnify_database_path=$MGNIFY_DATABASE_PATH \
  --template_mmcif_dir=$TEMPLATE_MMCIF_PATH \
  --bfd_database_path=$BFD_DATABASE_PATH \
  --uniref30_database_path=$UNIREF30_DATABASE_PATH \
  --pdb70_database_path=$PDB70_DATABASE_PATH \
  --obsolete_pdbs_path=$OBSOLETE_PDB_PATH \
  --use_gpu_relax=True
```
