# go-FAnnoT: Functional Annotation Transfer tool in Golang

## About

`go-FAnnoT` is functional annotion transfer tool based on protein homology. Our motivations to develop this tools were manyfold:

* __Defining a precise strategy to build reference datasets__. Indeed, most of time, transfer tools consider the annotation of one closely related species annotation as reference, copying possible errors. While it is necessary to adapt reference proteins to the organisms, a more robust strategy is required to ensure the quality of functional annotations.
* __Evaluating homology from global alignment and not from a local alignment__. Most of the existing tools identify matches on a basis of BLAST search. Unfortunatly, measuring homology on BLAST alignment is not sufficient and sequences should be realigned with a global alignment tool.
* __Allowing a flexible thresold setting__. In addition to reference datasets, homology thresholds should depends on the organism to annotate. Hence, for example, it can be necessary to lower threshold for species that does not have closely related species in reference databases.
* __Standardizing functional annotation in sequence files__. This latter aspect is critical to facilitate annotation comparisons.

Hence, `go-FAnnoT` broadly consists in the following steps:

1. __Extracting reference datasets from rich and high quality databases.__ We decided to use `Uniprot` and `TrEMBL`.
2. __Building a hierarchy between the different reference datasets.__
3. __Defining rules (different levels of homolgy) to transfer annotation.__
4. __Process each input proteins iteratively against each datasets until finding a suitable annotation.__
5. __(optional) Complete annotation with InterProScan functional domain prediction.__
6. __Produce standardized functional annotations.__

## Requierments

### Download Uniprot and TrEMBL databases

Our tool has been design to use Uniprot databases (`SwissProt` or `TrEMBL`). The complete `SwissProt` database can be downloaded  [here](https://ftp.uniprot.org/pub/databases/uniprot/current_release/knowledgebase/complete/) (choose the file __uniprot_sprot.dat.gz__)

Concerning the `TrEMBL` data, it is recommanded to download only a subset of the database as the complete one is too loarge. Thus, taxon level subsets are available [here](https://ftp.uniprot.org/pub/databases/uniprot/current_release/knowledgebase/taxonomic_divisions/).

### External tools

To run `go-FAnnoT`, it is necessary to have __NCBI-BLAST+__ tool suite and __NEEDLE__ (from __EMBOSS__ tool suite) in the system `PATH`. To do so, there are several solutions:

* Use a __conda__ environment with these two tools.
* (Or) Install these tools. Binaries are available at the following urls:
    * [NCBI-BLAST+](https://ftp.ncbi.nlm.nih.gov/blast/executables/blast+/LATEST/)
    * [EMBOSS](ftp://emboss.open-bio.org/pub/EMBOSS/emboss-latest.tar.gz)
* (Or, for __linux__ only) Most of the recent distributions have these tools available directly in there repositories: 

```
# Example with Ubuntu
apt-get install ncbi-blast+ emboss
```

## Install `go-FAnnoT`

### Build the project from source (github)

To build the project you will have to install Go (see instructions [here](https://go.dev/doc/install)).

Then clone this repository:

```
git clone https://github.com/hdevillers/go-fannot.git
```

Enter the `go-fannot`directory and build the project with `make` instructions:

```
cd go-fannot
make
make test
```

For __linux__ and __macos__, binary can be installed by running `make install` with administrator rights. The default installation path is `/usr/local/bin/`. It is possible to indicate a different installation path as follow:

```
make install -prefix my/install/path
```

### Download binaries

Precompiled binaries for all platforms will be available soon.


## Licence

[MIT](https://opensource.org/licenses/MIT)