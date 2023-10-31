# AlphaFold-Wrapper for small HPC systems

[![unittests](https://github.com/rvhonorato/alphafold-hpc-wrapper/actions/workflows/unittests.yml/badge.svg)](https://github.com/rvhonorato/alphafold-hpc-wrapper/actions/workflows/unittests.yml)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/a7bc51bb94d748d5ad4e4fc237cb1982)](https://app.codacy.com/gh/rvhonorato/alphafold-hpc-wrapper/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/a7bc51bb94d748d5ad4e4fc237cb1982)](https://app.codacy.com/gh/rvhonorato/alphafold-hpc-wrapper/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Unlicense](https://img.shields.io/badge/License-Unlicense-blue.svg)](https://opensource.org/license/unlicense/)

This repository contains a set of instructions on how-to-install alphafold on a HPC cluster. It is based on the [alphafold repository](https://github.com/google-deepmind/alphafold).

There they go over the instructions to install alphafold on a single machine with Docker - however you might want to install it in your local cluster. This repository contains a step-by-step installation procedure as well as a custom-built wrapper CLI to facilitate the execution of alphafold in a _small_ HPC environment.

Note that this is **not** an in-depth guide and **not** affiliated with alphafold - please check their documentation for an updated and official instructions. This repo is rather a set of instructions of a case that worked for me. I hope it can be useful for others as well.

- [Installation](#installation)
- [Wrapper](#wrapper)
  - [Build](#build)
  - [Configuration](#configuration)
  - [Execution](#execution)
- [TODO](#todo)
- [Support](#support)

## Installation

Please refer to the [INSTALL.md](INSTALL.md) file for a step-by-step installation procedure. These will show how to setup this repository, install the dependencies, configure a self-contained python environment and install alphafold.

## Wrapper

Alphafold's CLI is not very user-friendly and is not integrated with the HPC's scheduling system. This repository contains a wrapper CLI that facilitates the execution of alphafold in a small HPC environment.

It aims to expose less parameters and automates the execution of Alphafold via SLURM. It also provides a simple way to configure the execution via environment variables.

It assumes that you followed the [installation procedure](INSTALL.md); it will create a SLURM job file on the fly (in memory) with the correct initialization of the miniconda alphafold environment, the arguments needed by alphafold and the correct paths to the data.


### Build

[Install go](https://go.dev/doc/install) and build from source:

```
$ cd wrapper
$ go build -o alphafold-wrapper
```

### Configuration

The wrapper will take the configuration from the system variables, currently it supports the following variables:

- `PARTITION`: The SLURM partition to use.
- `INSTALL_DIR`: The installation directory of this repository.
- `DATA_DIR`: The path to the data directory.

Either define them in your environment or directly from the command line:

```bash
$ export PARTITION=gpu
$ export INSTALL_DIR=/trinity/login/rodrigo/repos/alphafold-wrapper
$ export DATA_DIR=/trinity/login/rodrigo/repos/alphafold-wrapper/data
```

### Execution

The build produces a binary called `alphafold-wrapper`. It has the following usage:

```bash
$ ./alphafold-wrapper -h
Usage:  ./alphafold-wrapper -f target.fasta -p <monomer|monomer_casp14|monomer_ptm|multimer> -o /path/to/output

### To fold a multimer, the target.fasta must be a multi-fasta file
### Ex: `https://www.rcsb.org/fasta/entry/2OOB/display`

```


## TODO

- Add more customization options to the SLURM job file
- Allow to re-use MSAs
- Add automated tests to the wrapper
- Add a build pipeline

## Support

If you have any questions or suggestions, please open an issue.
