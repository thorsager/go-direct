#!/usr/bin/env bash
if [[ -z "$REGISTRY_USERNAME" ]] || [[ -z "$REGISTRY_PASSWORD" ]]; then
	echo "both REGISTRY_USER and REGISTRY_PASSWORD are required, please set"
	exit 7	
fi
echo "$REGISTRY_PASSWORD" | docker login -u "$REGISTRY_USERNAME" --password-stdin
