#!/usr/bin/env bash

pushd $(dirname "$0") > /dev/null

VERSION=$(git describe --abbrev=0 --tags)
TAG=$(git describe --tag)

if [[ ${TAG} = ${VERSION} ]]; then
	echo ${VERSION}
else
	# Current commit is not a tag, dev commit
	HASH=$(git rev-parse --short HEAD)
	echo "${VERSION} dev ${HASH}"
fi
popd >/dev/null
