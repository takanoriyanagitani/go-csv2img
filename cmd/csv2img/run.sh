#!/bin/sh

export ENV_WIDTH=5
export ENV_HEIGHT=7

output=./sample.d/out.png

mkdir -p sample.d

input(){
	echo 0,1,2,3,4
	echo 1,1,2,3,4
	echo 2,1,2,3,4
	echo 3,1,2,3,4
	echo 4,1,2,3,4
	echo 5,1,2,3,4
	echo 6,1,2,3,4
}

input |
	./csv2img |
	dd \
		if=/dev/stdin \
		of="${output}" \
		bs=1048576 \
		status=progress \
		conv=fsync

file "${output}"
