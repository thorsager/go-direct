#!/usr/bin/env bash
while getopts ":n:" opt; do
	case "$opt" in 
		n)
			IMGNAME="$OPTARG"
			;;
		[?])
			echo >&2 "USAGE: $(basename $0) -n <name>"
			echo >&2 ""
			echo >&2 "Build a docker image using '.' as context"
			echo >&2 ""
			echo >&2 "Options:"
			echo >&2 "  -n <name>   Name of the new image, can containg :<tag>"
			echo >&2 ""
	esac
done

if [ -z "$IMGNAME" ]; then
	echo >&2 "Tag is required, use -t"
	exit 7
fi

echo "==> building $IMGNAME"
docker build -t ${IMGNAME} .
