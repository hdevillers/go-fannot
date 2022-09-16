#!/usr/bin/env python

import sys
import getopt

def main(argv):
    # Default values
    in_annot = ''
    in_seq = ''
    in_format = 'embl'
    out_seq = ''
    out_format = 'embl'

    # Parse arguments
    try:
        opts, args = getopt.getopt(argv, "ha:s:f:o:F:", ["help", "annotations=", "seq-files=", "in-format=", "output=", "out-format="])
    except getopt.GetopError:
        print 'An error occured while parsing input arguments.'
        sys.exit(2)
    for opt, arg in opts:
        if opt == '-h':
            print 'Help'
            sys.exit(1)
        elif opt in ('-a', '--annotations'):
            in_annot = arg

if __name__ == "__main__":
    main(sys.argv[1:])