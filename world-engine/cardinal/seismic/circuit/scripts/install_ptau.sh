#! /bin/bash

PTAU_URL=https://storage.googleapis.com/zkevm/ptau/powersOfTau28_hez_final_17.ptau
PTAU_FILE=powersOfTau28_hez_final_17.ptau


if [ ! -f artifacts/$PTAU_FILE ]; then
  mkdir -p artifacts
  curl -o artifacts/$PTAU_FILE $PTAU_URL
  echo " == Installed dev ptau file: $PTAU_FILE"
else
  echo " == Dev ptau file already installed"
fi
