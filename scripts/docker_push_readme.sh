#!/usr/bin/env bash
while getopts ":hf:" opt; do
	case "$opt" in
		f)
			README_FILE="$OPTARG"
			;;
		r)
			REPO_NAME="$OPTARG"
			;;
		h)
			echo "USAGE: $(basename $0) [-h] -f <readme-file>  -r <full-repo-name>" >&2
			echo "" >&2
			echo "Push readme file to docker registry repo" >&2
			echo "" >&2
			echo "Flags:" >&2
			echo "  -h     Show this this help text" >&2
			echo "" >&2
			echo "Options:" >&2
			echo "  -f <readme-file>    File to be pushed to registry-repo as README.md" >&2
			echo "  -r <full-repo-name> Name of te repo on docker hub (ex. thorsger/go-direct)" >&2
			echo "" >&2
			exit 1
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

if [ -z "$README_FILE" ]; then
	echo "readme-file is requried, use -f <readme-file>" >&2
	exit 7
fi

if [ -z "$REPO_NAME" ]; then
	echo "repo-name is requried, use -r <repo-name>" >&2
	exit 7
else 
	R_PREFIX=$(echo $REPO_NAME | cut -d '/' -f1)
	R_NAME=$(echo $REPO_NAME | cut -d '/' -f2) 
fi

if [[ -z "$R_PREFIX" ]] || [[ -z "$R_NAME" ]]; then
	echo "invalid repo-name, must be '<org>/<repo>'"
	exit 7
fi

if [[ -z "$REGISTRY_USERNAME" ]] || [[ -z "$REGISTRY_PASSWORD" ]]; then
	echo "both REGISTRY_USER and REGISTRY_PASSWORD are required, please set"
	exit 7
fi

echo "==> pushing $README_FILE to $REPO_NAME"
docker run --rm \
    -v $README_FILE:/data/README.md \
    -e DOCKERHUB_USERNAME=${REGISTRY_USERNAME} \
    -e DOCKERHUB_PASSWORD=${REGISTRY_PASSWORD} \
    -e DOCKERHUB_REPO_PREFIX=$R_PREFIX \
    -e DOCKERHUB_REPO_NAME=${R_NAME} \
     sheogorath/readme-to-dockerhub
