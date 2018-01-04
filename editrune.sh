#! /bin/bash

set -eu

source loadconf.sh

if [ $# -lt 1 ]; then
	echo "Few arg"
	exit 1
fi

CHAR=$1

FILE=`./runeid "$CHAR"`

$EDITOR $FILE
