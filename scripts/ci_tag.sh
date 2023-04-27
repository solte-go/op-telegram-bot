#!/bin/bash

setOutput() {
    echo "${1}=${2}" >> "${GITHUB_OUTPUT}"
}

tagFmt="^v?[0-9]+\.[0-9]+\.[0-9]+$"
tag="$(git for-each-ref --sort=-v:refname --format '%(refname:lstrip=2)' | grep -E "$tagFmt" | head -n 1)"

setOutput "tag" "$tag"