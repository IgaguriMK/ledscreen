#! /bin/bash

set -eu

if [ ! -v FONT ]; then
	FONT=Noto-Sans-Mono-Regular
fi

if [ $# -lt 1 ]; then
	echo "No arg"
	exit 1
fi

if [ $# -eq 1 ]; then
	CHAR=$1
else
	EDITOR=$1
	CHAR=$2
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

if [ -v EDITOR ]; then
	$EDITOR $FILE
fi
