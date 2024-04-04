pragma circom 2.1.1;

include "salts.func.circom";
include "choice.func.circom";
include "commitments.func.circom";

template Init() {
    var N_STAFFS = 2;
    var N_STAFF_SPELLS = 1;
    var N_ALL_SPELLS = 2;

    assert(N_STAFFS * N_STAFF_SPELLS <= N_ALL_SPELLS);

    signal input playerSource;
    signal input commitments[N_STAFFS][N_STAFF_SPELLS];
    
    signal input seismicSource;
    signal input attributes[N_STAFFS][N_STAFF_SPELLS];
    signal input salts[N_STAFFS][N_STAFF_SPELLS];

    signal saltsCorrect <== CheckAllSalts(N_STAFFS, N_STAFF_SPELLS)(seismicSource, salts);
    saltsCorrect === 1;

    signal sourcedChoice <== CheckSourcedChoice(N_STAFFS, N_STAFF_SPELLS, 
        N_ALL_SPELLS)(playerSource, seismicSource, attributes);
    sourcedChoice === 1;

    signal commitmentsCorrect <== CheckAllCommitments(N_STAFFS, N_STAFF_SPELLS)(
        commitments, attributes, salts);
    commitmentsCorrect === 1;
}

component main { public [ playerSource, commitments ] } = Init();
