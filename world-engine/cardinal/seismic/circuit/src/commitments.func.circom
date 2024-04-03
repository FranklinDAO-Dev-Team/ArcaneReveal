pragma circom 2.1.1;

include "../node_modules/circomlib/circuits/gates.circom";
include "../node_modules/circomlib/circuits/comparators.circom";
include "../node_modules/circomlib/circuits/poseidon.circom";

template CheckAllCommitments(N_STAFFS, N_STAFF_SPELLS) {
    signal input commitments[N_STAFFS][N_STAFF_SPELLS];

    signal input attributes[N_STAFFS][N_STAFF_SPELLS];
    signal input salts[N_STAFFS][N_STAFF_SPELLS];
    signal output out;

    signal commitmentsCorrect[N_STAFFS][N_STAFF_SPELLS];
    signal accumulator[N_STAFFS * N_STAFF_SPELLS + 1];
    accumulator[0] <== 1;

    for (var s = 0; s < N_STAFFS; s++) {
        for (var p = 0; p < N_STAFF_SPELLS; p++) {
            commitmentsCorrect[s][p] <== CheckCommitment()(commitments[s][p], 
                attributes[s][p], salts[s][p]);

            var accIndex = 1 + (s * N_STAFF_SPELLS) + p;
            accumulator[accIndex] <== AND()(accumulator[accIndex - 1],
                commitmentsCorrect[s][p]);
        }
    }

    out <== accumulator[N_STAFFS * N_STAFF_SPELLS];
}

template CheckCommitment() {
    signal input commitment;

    signal input id;
    signal input salt;

    signal output out;

    signal circuitCommitment <== Poseidon(2)([id, salt]);
    out <== IsEqual()([commitment, circuitCommitment]);
}
