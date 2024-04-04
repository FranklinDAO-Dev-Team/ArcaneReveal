pragma circom 2.1.1;

include "../node_modules/circomlib/circuits/gates.circom";
include "../node_modules/circomlib/circuits/comparators.circom";
include "../node_modules/circomlib/circuits/poseidon.circom";
include "permutation.circom";

template UpToN(N) {
    signal output out[N];
    for (var i = 0; i < N; i++) {
        out[i] <== i;
    }
}

template CheckSourcedChoice(N_STAFFS, N_STAFF_SPELLS, N_ALL_SPELLS) {
    signal input playerSource;
    signal input seismicSource;

    signal input attributes[N_STAFFS][N_STAFF_SPELLS];
    signal output out;

    // TODO: can two staffs hold the same spell atst?

    signal matchHash <== Poseidon(2)([playerSource, seismicSource]);

    var spells[N_ALL_SPELLS] = UpToN(N_ALL_SPELLS)();

    signal circuitAttributes[N_ALL_SPELLS] <== RandomPermutate(N_ALL_SPELLS)(
       matchHash, spells);

    var nTotalSpells = N_STAFFS * N_STAFF_SPELLS;
    signal spellsCorrect[nTotalSpells];
    signal accumulator[nTotalSpells + 1];
    accumulator[0] <== 1;

    for (var s = 0; s < N_STAFFS; s++) {
        for (var p = 0; p < N_STAFF_SPELLS; p++) {
            var spellIndex = (s * N_STAFF_SPELLS) + p;
            var accIndex = spellIndex + 1;

            spellsCorrect[spellIndex] <== IsEqual()([attributes[s][p], 
                circuitAttributes[spellIndex]]);

            accumulator[accIndex] <== AND()(accumulator[accIndex - 1], 
                spellsCorrect[spellIndex]);
        }
    }

    out <== accumulator[N_STAFFS * N_STAFF_SPELLS];
}
