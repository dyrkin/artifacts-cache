#!/bin/sh -ex

rm -rf files
./artifacts pull 'files/*.txt'
count=$(ls files | wc -l)

if [[ "$count" -ne 2 ]]; then
  echo "invalid file count $count"
  exit 1
fi
file1=$(cat files/file1.txt)
if [[ "$file1" != "Content 1" ]]; then
  echo "invalid file1 content $file1"
  exit 1
fi
file2=$(cat files/file2.txt)
if [[ "$file2" != "Content 2" ]]; then
  echo "invalid file2 content $file2"
  exit 1
fi
