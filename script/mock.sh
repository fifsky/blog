#!/usr/bin/env bash
set -e

APIDIR="proto/gen"

if [ -n "$1" ]; then
  APIDIR=$(find "$APIDIR" -type d -path "*/$1")
fi

gfilelist=$(find ${APIDIR} -type f \( -name "*_grpc.pb.go" -o -name "*_job.pb.go" \))

if [ -z "$gfilelist" ]; then
  echo "[Generate mock file] no grpc file found"
  exit 0
fi

loadVersion && checkMockVersion

filelist=$(find ${APIDIR} -type f \( -name "*_grpc.pb.go" -o -name "*_job.pb.go" \) | grep -v "mock_")

for file in ${filelist}; do
  fname=$(basename ${file})
  dname=$(dirname ${file})
  pname=$(grep "package " ${file} | cut -d" " -f2)
  echo "[Generate mock file]${file}" filename=${fname} dirname=${dname} package=${pname}
  mockgen -source=${file} -destination=${dname}/mock_${fname} -package=${pname} &
done
