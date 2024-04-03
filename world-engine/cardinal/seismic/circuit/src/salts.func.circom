pragma circom 2.1.1;

include "../node_modules/circomlib/circuits/gates.circom";
include "../node_modules/circomlib/circuits/comparators.circom";
include "../node_modules/circomlib/circuits/poseidon.circom";

template CheckAllSalts(N_STAFFS, N_STAFF_SPELLS) {
    signal input seismicSource;
    signal input salts[N_STAFFS][N_STAFF_SPELLS];
    
    signal output out;

    signal saltsCorrect[N_STAFFS][N_STAFF_SPELLS];
    signal accumulator[N_STAFFS * N_STAFF_SPELLS + 1];
    accumulator[0] <== 1;

    for (var s = 0; s < N_STAFFS; s++) {
        for (var p = 0; p < N_STAFF_SPELLS; p++) {
            saltsCorrect[s][p] <== CheckSalt()(seismicSource, salts[s][p], 
                s * N_STAFF_SPELLS + p);

            var accIndex = 1 + (s * N_STAFF_SPELLS) + p;
            accumulator[accIndex] <== AND()(accumulator[accIndex - 1],
                saltsCorrect[s][p]);
        }
    }

    out <== accumulator[N_STAFFS * N_STAFF_SPELLS];
}

template CheckSalt() {
    signal input seismicSource;
    signal input salt;

    signal input nonce;

    signal output out;

    signal circuitSalt <== Poseidon(2)([nonce, seismicSource]);
    out <== IsEqual()([salt, circuitSalt]);
}
