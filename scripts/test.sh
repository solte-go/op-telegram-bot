#!/bin/bash

# shellcheck disable=SC2039

regex='v?[0-9]+\.[0-9]+\.[0-9]+$'

tag="$(git show -s --format=%B | grep "#tag")"
echo "$tag"
if [[ $tag =~ $regex ]];
then
  echo "${BASH_REMATCH[0]}"
else
  echo "doesn't match" # this could get noisy if there are a lot of non-matching files
  tag='0.0.0'
fi
