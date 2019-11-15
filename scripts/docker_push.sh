#!/usr/bin/env bash
while getopts ":ln:" opt; do
	case "$opt" in 
		l)
			LOGIN=true
			;;
		n)
			NAME="$OPTARG"
			;;
		\?) 
			echo "Unknown option: -$OPTARG" >&2
			exit 7
			;;
        :) 
			echo "Missing option argument for -$OPTARG" >&2
			exit 7
			;;
		*) 
			echo "Unimplemented option: -$OPTARG" >&2 
			exit 7
			;;
	esac
done

if [ -z "$NAME" ]; then
	echo >&2 "imange name is required, use -n"
	exit 7
fi

if [ -n "$LOGIN" ]; then
	if [[ -z "$REGISTRY_USERNAME" ]] || [[ -z "$REGISTRY_PASSWORD" ]]; then
		echo "both REGISTRY_USER and REGISTRY_PASSWORD are required, please set"
		exit 7	
	fi
	echo "==> logging in as $REGISTRY_USERNAME"
	echo "$REGISTRY_PASSWORD" | docker login -u "$REGISTRY_USERNAME" --password-stdin
fi

echo "==> pushing $NAME"
docker push "$NAME"
