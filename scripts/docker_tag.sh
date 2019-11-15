#!/usr/bin/env bash

while getopts ":s:t:" opt; do
	case "$opt" in 
		s) 
			SOURCE="$OPTARG"
			;;
		t) 
			TARGET="$OPTARG"
			;;
		[?])
			echo >&2 "USAGE: $(basename $0) -s <source> -t <tag>"
			echo >&2 ""
			echo >&2 "Add additional tag to existing docker image"
			echo >&2 ""
			echo >&2 "Options:"
			echo >&2 "  -s <source>  Source image (and tag) to which the additional"
			echo >&2 "               tag is to be added"
			echo >&2 "  -t <tag>     Name of the tag to be added"
			echo >&2 ""
			exit 7
	esac
done
if [[ -z "$SOURCE" ]] || [[ -z "$TARGET" ]]; then
	echo >&2 "both source and tag are requred, use -s and -t"
	exit 7
fi

echo "==> tagging $SOURCE as $TARGET"
docker tag "$SOURCE" "$TARGET" 
