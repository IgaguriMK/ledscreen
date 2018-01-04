#! /bin/bash

if [ $# -lt 2 ]; then
	echo "Few args"
	exit 1
fi

ffmpeg -r 30 -i "${1}%06d.png" -vcodec libx264 -pix_fmt yuv420p ${2}.mp4
