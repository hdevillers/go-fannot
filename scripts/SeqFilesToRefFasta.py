#!/usr/bin/env python

import argparse
from Bio import SeqIO
from Bio.SeqRecord import SeqRecord
import glob
from os.path import exists
from os.path import basename
import sys

def fatal(msg):
    print('[ERROR]: '+msg, file=sys.stderr)
    sys.exit(2)

def warn(msg):
    print('[WARNING]: '+msg, file=sys.stderr)

def main(argv):
    # Default values
    allowed_format = ['embl', 'gb', 'genbank']

    # Parse arguments
    parser = argparse.ArgumentParser(prog="SeqFilesToRefFasta.py")
    parser.add_argument('-s', '--seq-files', help="Sequence file(s) (embl or genbank) to write annotations.", required=True)
    parser.add_argument('-f', '--in-format', default="embl", help = "Input sequence file format.")
    parser.add_argument('-o', '--output', default="./", help="Output file names.")
    parser.add_argument('-i', '--id-qualifier', default="locus_tag", help="Target qualifier to get the id of a feature.")
    args = parser.parse_args()

    # Check argument values
    if args.seq_files is None:
        fatal('You must provide an input sequence file.')
    if args.in_format not in allowed_format:
        fatal('Input format ('+args.in_format+') not supported')
    
    # Check input sequence files
    seq_files = []
    if exists(args.seq_files):
        seq_files.append(args.seq_files)
    else:
        for fn in glob.glob(args.seq_files):
            seq_files.append(fn)
    if len(seq_files) == 0:
        fatal("Failed to find input sequence file(s).")

    # Init. array of bioseq object
    proteins = []

    # Scan each record in the input files
    for f in seq_files:
        with open(f) as hdl:
            for record in SeqIO.parse(hdl, args.in_format):
                for feature in record.features:
                    # Only check CDS
                    if feature.type == 'CDS':
                        id = "CDS_%05d" % (len(proteins))
                        gene_name = ""
                        procuct = ""
                        note=""
                        function=""
                        if args.id_qualifier in feature.qualifiers:
                            id = feature.qualifiers[args.id_qualifier][0]
                        if 'note' in feature.qualifiers:
                            note = feature.qualifiers['note'][0]
                        if 'product' in feature.qualifiers:
                            product = feature.qualifiers['product'][0]
                        if 'gene' in feature.qualifiers:
                            gene = feature.qualifiers['gene'][0]
                        if 'function' in feature.qualifiers:
                            function = feature.qualifiers['function'][0]
                        protein = feature.translate(record.seq, cds=False)
                        prot_record = SeqRecord(
                            protein,
                            id = id,
                            description = "%s::%s::%s::::%s" % (product, gene_name, function, note)
                        )
                        proteins.append(prot_record)
            hdl.close()
    
    # Write out proteins
    with open(args.output, 'w') as hdl:
        for rec in proteins:
            SeqIO.write(rec, hdl, 'fasta')
        hdl.close()

if __name__ == "__main__":
    main(sys.argv[1:])


