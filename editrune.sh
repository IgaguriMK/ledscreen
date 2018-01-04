#! /bin/bash

set -eu

if [ $# -lt 2 ]; then
	echo "Few arg"
	exit 1
fi

EDITOR=$1
CHAR=$2

FILE=`./runeid "$CHAR"`

$EDITOR $FILE
