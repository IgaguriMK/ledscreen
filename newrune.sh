#! /bin/bash

FONT=./fonts/00-jp.otf

set -eu

source loadconf.sh

if [ $# -lt 1 ]; then
	echo "No arg"
	exit 1
fi

if [ $# -eq 1 ]; then
	CHAR=$1
else
	CHAR=$2
	FONT=$1
fi


FILE=`./runeid "$CHAR"`

if [ -e $FILE ]; then
	echo "Exist file"
	exit 1
fi

if ./runewidth "$CHAR"; then
	SIZE=12
else
	SIZE=16
fi

convert -size ${SIZE}x16 -font $FONT -gravity Center -fill white -background black -pointsize 18 "label:$CHAR" $FILE

$EDITOR $FILE
