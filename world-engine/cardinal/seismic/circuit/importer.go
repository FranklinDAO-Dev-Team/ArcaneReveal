package circuit

import _ "embed"

//go:embed artifacts/circuit.wasm
var CircuitWasm []byte

//go:embed artifacts/circuit.zkey
var CircuitZkey []byte

//go:embed artifacts/verification_key.json
var VerificationKey []byte
