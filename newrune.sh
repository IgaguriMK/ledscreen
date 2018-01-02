#! /bin/bash

if [ "$FONT" == "" ]; then
	FONT=Noto-Sans-Mono-CJK-JP-Regular
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

convert -size 16x16 -font $FONT -gravity Center -fill white -background black -pointsize 18 "label:$CHAR" $FILE

if [ "$EDITOR" != "" ]; then
	$EDITOR $FILE
fi
