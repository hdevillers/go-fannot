#!/usr/bin/env python

import argparse
from Bio import SeqIO
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
    parser = argparse.ArgumentParser(prog="AddCommentsToSeqFiles.py")
    parser.add_argument('-s', '--seq-files', help="Sequence file(s) (embl or genbank) to write annotations.", required=True)
    parser.add_argument('-f', '--in-format', default="embl", help = "Input sequence file format.")
    parser.add_argument('-o', '--out-dir', default="./", help="Output directory.")
    parser.add_argument('-F', '--out-format', default="", help="Output sequence file format.")
    parser.add_argument('-g', '--go-fannot-version', default="1.1.0", help="Version of go-fannot used.")
    parser.add_argument('-u', '--uniprot-version', default="", help="Version of Uniprot used.")
    parser.add_argument('-i', '--interproscan-version', default="", help="Version of InterProScan used.")
    args = parser.parse_args()

    # Check argument values
    if args.seq_files is None:
        fatal('You must provide an input sequence file.')
    if args.in_format not in allowed_format:
        fatal('Input format ('+args.in_format+') not supported')
    if args.out_format == '':
        args.out_format = args.in_format
    if args.out_format not in allowed_format:
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
    
    # Prepare output files
    seq_out = []
    ask_owt = False
    for fn in seq_files:
        bn = basename(fn)
        so = args.out_dir + '/' + bn
        if exists(so):
            ask_owt = True
            warn("The output sequence file (" + so + ") already exists.")
        seq_out.append(so)
    if ask_owt:
        asw = input("Overwrite existing output file(s)? [No/yes]: ")
        if asw == "" or asw.lower() == "no" or asw.lower() == "n":
            fatal("Please, check output arguments.")
        elif asw.lower() == "yes" or asw.lower() == "y":
            pass
        else:
            fatal(asw + " is not an appropriate answer.")
    
    # Scan each sequence file(s)
    for i in range(len(seq_files)):
        path_in = seq_files[i]
        path_out = seq_out[i]
        with open(path_in) as hdl_in:
            edited_records = []
            for record in SeqIO.parse(hdl_in, args.in_format):
                # Get the current comment
                cc = record.annotations.get("comment")
                if cc == None:
                    cc = ""
                else:
                    cc = cc + "\n\n"
                
                # Add information
                cc = cc + "Functional annotation performed with Go-Fannot v"+args.go_fannot_version
                if args.uniprot_version != "" or args.interproscan_version != "":
                    cc = cc + "\nReference annotations were based on:"
                    if args.uniprot_version != "":
                        cc = cc + "\n    - UniProt v" + args.uniprot_version
                    if args.interproscan_version != "":
                        cc = cc + "\n    - InterProScan v" + args.interproscan_version
                
                # Set comment
                record.annotations["comment"] = cc
                edited_records.append(record)
            # Save edited record(s)
            with open(path_out, 'w') as hdl_out:
                for record in edited_records:
                    SeqIO.write(record, hdl_out, args.out_format)
    
if __name__ == "__main__":
    main(sys.argv[1:])
