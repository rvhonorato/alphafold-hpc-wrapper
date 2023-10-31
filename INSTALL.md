# Installation

The following instructions are based on the [alphafold repository](https://github.com/google-deepmind/alphafold) and assumes you are using a Linux system with a CUDA-enabled GPU, it was tested on:

```bash
$ cat /etc/centos-release
CentOS Linux release 7.9.2009 (Core)

$ nvidia-smi --query-gpu=name --format=csv,noheader
NVIDIA GeForce GTX 1080 Ti
NVIDIA GeForce GTX 1080 Ti
NVIDIA GeForce GTX 1080 Ti
NVIDIA GeForce GTX 1080 Ti
```

1. [Clone this repository](#clone-this-repository)
2. [Install `aria2c`](#aria2c)
3. [Create a Python environment with Miniconda](#python-environment)
4. [Download the data needed for Alphafold](#download-data)
5. [Setup Alphafold](#setup)
6. [Make a prediction](#make-a-prediction)

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
```

## Alphafold

### Download data

Note that the uncompressed data size is ~550 GB so send it to the background with `nohup`. This might take a **LONG** time, so be patient and check the log file if needed.

```bash
$ cd $ALPHAFOLD_WRAPPER
$ mkdir -p $ALPHAFOLD_WRAPPER/data
$ nohup bash $ALPHAFOLD_WRAPPER/alphafold/scripts/download_all_data.sh $ALPHAFOLD_WRAPPER/data &

# inspect the log if needed
$ tail nohup.out
```

### Download `stereo_chemical_props.txt`

```bash
wget https://git.scicore.unibas.ch/schwede/openstructure/-/raw/7102c63615b64735c4941278d92b554ec94415f8/modules/mol/alg/src/stereo_chemical_props.txt -O $ALPHAFOLD_WRAPPER/alphafold/alphafold/common/stereo_chemical_props.txt
```

### Setup

Find out what is your CUDA version and set the `CUDA_VERSION` variable

```bash
$ nvidia-smi | grep "CUDA Version"
| NVIDIA-SMI 510.47.03    Driver Version: 510.47.03    CUDA Version: 11.6     |
$ export CUDA_VERSION=11.6
```

Configure the Python environment and its dependencies

```bash
$ source $ALPHAFOLD_WRAPPER/miniconda3/etc/profile.d/conda.sh

# (Optional) Configure the base env with libmamba for faster dependency solution
$ conda update -n base conda
$ conda install -y -n base conda-libmamba-solver
$ conda config --set solver libmamba

# Create and activate the alphafold environment
$ conda create -y -n af2 python=3.10 && conda activate af2

# Install conda packages
(af2) $ conda install -y -c conda-forge -c bioconda \
  tensorflow-gpu \
  hmmer==3.3.2 \
  kalign3 \
  hhsuite==3.3.0 \
  openmm=7.7.0 \
  cudatoolkit==${CUDA_VERSION} \
  absl-py==1.0.0 \
  biopython==1.79 \
  chex==0.0.7 \
  dm-haiku==0.0.10 \
  dm-tree==0.1.8 \
  immutabledict==2.0.0 \
  numpy==1.24.3 \
  pandas==2.0.3 \
  scipy==1.11.1 \
  ml_dtypes==0.2.0 \
  pdbfixer \
  pip

# Install the PIP packages
## mind the cuda version here!
(af2) $ pip install jax==0.4.14 jaxlib==0.4.14+cuda11.cudnn86 -f https://storage.googleapis.com/jax-releases/jax_cuda_releases.html
(af2) $ pip install ml-collections==0.1.0
```

```bash
# Make sure the GPU is available
(af2) $ python -c "import tensorflow as tf; print(tf.config.list_physical_devices('GPU'))"
[PhysicalDevice(name='/physical_device:GPU:0', device_type='GPU'), PhysicalDevice(name='/physical_device:GPU:1', device_type='GPU'), PhysicalDevice(name='/physical_device:GPU:2', device_type='GPU'), PhysicalDevice(name='/physical_device:GPU:3', device_type='GPU')]
```

Test the command, you should not see any errors or warnings!

```bash
(af2) $ python $ALPHAFOLD_WRAPPER/alphafold/run_alphafold.py --helpshort
Full AlphaFold protein structure prediction script.
flags:

/trinity/login/rodrigo/repos/alphafold-wrapper/alphafold/run_alphafold.py:
  --[no]benchmark: Run multiple JAX model evaluations to obtain a timing that excludes the compilation time, which should be more
    indicative of the time required for inferencing many proteins.
    (default: 'false')
  --bfd_database_path: Path to the BFD database for use by HHblits.
  # ... etc etc
```
