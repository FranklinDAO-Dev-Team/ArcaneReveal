package component

type WandCore struct {
	Number int
	// Abilities [client.NumAbilities]int // Array of 5 integers
	// Revealed  [client.NumAbilities]int // Slice of integers
}

func (WandCore) Name() string {
	return "Wand"
}

// func NewRandomWandCore() WandCore {
// 	w := WandCore{}

// 	// Set Revealed to all -1
// 	for i := range w.Revealed {
// 		w.Revealed[i] = -1
// 	}

// 	// Generate unique random numbers for Abilities
// 	uniqueNumbers := make(map[int]bool)
// 	for i := 0; i < client.NumAbilities; {
// 		num, err := cryptoRandInt(1, client.TotalAbilities) // Random number between 1 and 50
// 		if err != nil {
// 			panic(err)
// 		}

// 		if !uniqueNumbers[num] {
// 			uniqueNumbers[num] = true
// 			w.Abilities[i] = num
// 			i++
// 		}
// 	}

// 	return w
// }

// // cryptoRandInt generates a random integer between min and max using crypto/rand.
// func cryptoRandInt(min, max int) (int, error) {
// 	if max <= min {
// 		return 0, fmt.Errorf("max must be greater than min")
// 	}

// 	var b [8]byte
// 	_, err := rand.Read(b[:])
// 	if err != nil {
// 		return 0, err
// 	}

// 	// Convert the byte slice to an unsigned 64-bit integer
// 	randUint := binary.BigEndian.Uint64(b[:])

// 	// Scale the value to the desired range
// 	return int(randUint%uint64(max-min+1)) + min, nil
// }
