#! /bin/bash

PIC=output/
MOV=output

if [ $# -ge 1 ]; then
	PIC=$1
fi

if [ $# -ge 2 ]; then
	MOV=$2
fi

ffmpeg -r 30 -i "${PIC}%06d.png" -vcodec libx264 -pix_fmt yuv420p ${MOV}.mp4
