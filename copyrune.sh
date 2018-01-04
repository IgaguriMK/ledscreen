#! /bin/bash

set -eu

if [ $# -lt 2 ]; then
	echo "few args."
	exit 1
fi

FROM=`./runeid $1`
TO=`./runeid $2`
cp $FROM $TO
