#! /bin/bash

PTAU=artifacts/powersOfTau28_hez_final_17.ptau
OUT_DIR=artifacts

# Compile circuit
yarn circom2 src/init.circom --r1cs --wasm
if [ $? -ne 0 ]; then
    echo "Circuit compilation failed"
    exit 1
fi

# Generate proving key
yarn run snarkjs groth16 setup init.r1cs \
    ${PTAU} \
    circuit.zkey

# Generate verifying key
yarn run snarkjs zkey export verificationkey circuit.zkey \
    verification_key.json

# Compute witness, used as smoke test for circuit
node init_js/generate_witness.js \
     init_js/init.wasm \
     test/init.smoke.json \
     witness.wtns
rm -rf witness.wtns

# Clean up and save ZK files
mkdir -p ${OUT_DIR}
mv circuit.zkey verification_key.json ${OUT_DIR}
mv init_js/init.wasm ${OUT_DIR}/circuit.wasm

rm -rf init_js init.r1cs
