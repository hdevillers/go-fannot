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
    allowed_qualifiers = ['note', 'product', 'gene', 'function']


    # Parse arguments
    parser = argparse.ArgumentParser(prog="AnnotationsToSeqFiles.py")
    parser.add_argument('-a', '--annotations', help="Annotation file generated by go-FAnnoT (tsv).", required=True)
    parser.add_argument('-s', '--seq-files', help="Sequence file(s) (embl or genbank) to write annotations.", required=True)
    parser.add_argument('-f', '--in-format', default="embl", help = "Input sequence file format.")
    parser.add_argument('-o', '--out-dir', default="./", help="Output directory.")
    parser.add_argument('-F', '--out-format', default="embl", help="Output sequence file format.")
    parser.add_argument('-t', '--feature-type', default="CDS", help="Feature type to wrote annotation.")
    parser.add_argument('-k', '--keep-prev-annot', help="Keep previous annotation when editing a feature.", action="store_true")
    parser.add_argument('-i', '--id-qualifier', default="locus_tag", help="Target qualifier to get the id of a feature.")
    parser.add_argument('-q', '--wanted-qualifiers', default="note,product,gene", help="List of qualifier to transfer (coma separator).")
    parser.add_argument('-c', '--copy-gene-status', default=1, type=int, help="Minimal status value to allow gene name transfer.")
    args = parser.parse_args()

    # Check argument values
    if  args.annotations is None:
        fatal('You must provide an annotation file.')
    if args.seq_files is None:
        fatal('You must provide an input sequence file.')
    if args.in_format not in allowed_format:
        fatal('Input format ('+args.in_format+') not supported')
    if args.out_format not in allowed_format:
        fatal('Ouput format ('+args.out_format+') not supported')
    if args.wanted_qualifiers == '':
        fatal("You must provide at least one qualifier type to copy.")
    wq = args.wanted_qualifiers.split(',')
    for q in wq:
        if q not in allowed_qualifiers:
            fatal("The queried qualifier type ("+q+") is not supported.")

    # Load and store annotation data
    try:
        f = open(args.annotations, 'r')
    except OSError:
        fatal('Failed to open/read input annotation file')
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
                fatal('Warning, the gene id '+dt[0]+' seems to be duplicated.')
            # Extract data
            ann_data[dt[0]] = {
                'product': dt[1],
                'note': dt[2],
                'function': dt[3],
                'gene': dt[7],
                'status': int(dt[10]),
                'copied': False
            }

            # Check copy gene status
            if ann_data[dt[0]]['status'] < args.copy_gene_status:
                # Delete gene data to prevent copy
                ann_data[dt[0]]['gene'] = ''

            # Read next line
            line = f.readline()
        f.close()

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
                for feature in record.features:
                    # Check feature type
                    if feature.type == args.feature_type:
                        # Get feature name
                        if args.id_qualifier in feature.qualifiers:
                            # NOTE: test only the first value (should not be multiple)
                            if feature.qualifiers[args.id_qualifier][0] in ann_data:
                                id = feature.qualifiers[args.id_qualifier][0]
                                ann_data[id]['copied'] = True
                                if args.keep_prev_annot:
                                    for q in wq:
                                        if ann_data[id][q] != "":
                                            if q in feature.qualifiers:
                                                feature.qualifiers[q].append(ann_data[id][q])        
                                            else:
                                                feature.qualifiers[q] = [ann_data[id][q]]
                                else:
                                    for q in wq:
                                        if ann_data[id][q] != "":
                                            feature.qualifiers[q] = [ann_data[id][q]]
                # Save the record
                edited_records.append(record)
            # Write the records in the output file
            with open(path_out, 'w') as hdl_out:
                for record in edited_records:
                    SeqIO.write(record, hdl_out, args.out_format)
    
    # Now scan annotation data to identify possible missed copies
    for id in ann_data.keys():
        if not ann_data[id]['copied']:
            warn("The annotation associated to id "+id+" has not been transfered.")


if __name__ == "__main__":
    main(sys.argv[1:])