#!/usr/bin/env python

from email.mime import base
from Bio import SeqIO
import getopt
import glob
from os.path import exists
from os.path import basename
import sys

def fatal(msg):
    print(msg, file=sys.stderr)
    sys.exit(2)

def warn(msg):
    print(msg, file=sys.stderr)

def main(argv):
    # Default values
    in_annot = ''
    in_seq = ''
    in_format = 'embl'
    out_dir = './'
    out_format = 'embl'

    # Parse arguments
    try:
        opts, args = getopt.getopt(argv, "ha:s:f:o:F:", ["help", "annotations=", "seq-files=", "in-format=", "output=", "out-format="])
    except getopt.GetopError:
        print('An error occured while parsing input arguments.', file=sys.stderr)
        sys.exit(2)
    for opt, arg in opts:
        if opt == '-h':
            print('Help')
            sys.exit(1)
        elif opt in ('-a', '--annotations'):
            in_annot = arg
        elif opt in ('-s', '--seq-files'):
            in_seq = arg
        elif opt in ('-f', '--in-format'):
            in_format = arg
        elif opt in ('-o', '--output'):
            out_dir = arg
        elif opt in ('-F', '--out-format'):
            out_format = arg
        else:
            print('Unknonwn argument: '+opt+'.', file=sys.stderr)
            sys.exit(2)

    # Check argument values
    if in_annot == '':
        print('You must provide an annotation file.', file=sys.stderr)
    if in_seq == '':
        print('You must provide an input sequence file.', file=sys.stderr)
    if in_format not in ('embl', 'bg', 'genbank'):
        print('Input format ('+in_format+') not supported', file=sys.stderr)
    if out_format not in ('embl', 'bg', 'genbank'):
        print('Ouput format ('+out_format+') not supported', file=sys.stderr)

    # Load and store annotation data
    try:
        f = open(in_annot, 'r')
    except OSError:
        print('Failed to open/read input annotation file', file=sys.stderr)
        sys.exit(2)
    # Init. annotation data hash
    ann_data = {}
    with f:
        # Skip header line
        line = f.readline()
        line = f.readline()
        while line:
            # split line
            dt = line.split("\t")
            # Check duplicated keys
            if dt[0] in ann_data:
                print('Warning, the gene id '+dt[0]+' seems to be duplicated.', file=sys.stderr)
                sys.exit(2)
            # Extract data
            ann_data[dt[0]] = {
                'Product': dt[1],
                'Note': dt[2],
                'Function': dt[3],
                'GeneName': dt[7],
                'Status': dt[10],
                'Copied': False
            }
            # Read next line
            line = f.readline()
        f.close()

    # Check input sequence files
    seq_files = []
    if exists(in_seq):
        seq_files.append(in_seq)
    else:
        for fn in glob.glob(in_seq):
            seq_files.append(fn)
    if len(seq_files) == 0:
        fatal("Failed to find input sequence file(s).")
    
    # Prepare output files
    seq_out = []
    ask_owt = False
    for fn in seq_files:
        bn = basename(fn)
        so = out_dir + '/' + bn
        if exists(so):
            ask_owt = True
            warn("The output sequence file (" + so + ") already exists.")
        seq_out.append(so)
    if ask_owt:
        asw = input("Overwrite existing output file(s)? [No/yes]: ")
        if asw == "" or asw.lower() == "no" or asw.lower == "n":
            fatal("Please, check output arguments.")
        elif asw.lower() == "yes" or asw.lower == "y":
            pass
        else:
            fatal(asw + " is not an appropriate answer.")

    # Scan each sequence file(s)
    


if __name__ == "__main__":
    main(sys.argv[1:])