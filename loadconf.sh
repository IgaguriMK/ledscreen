#! /bin/bash

if [ ! -r conf.sh ]; then
	echo "Not found: conf.sh"
	exit 1
fi

source conf.sh
