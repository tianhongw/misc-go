package util

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"math/rand"
)

func init() {
	// Ensure we use a high-entropy seed for the pseudo-random generator
	rand.Seed(newSeed())
}

// returns an int64 from a crypto random source
// can be used to seed a source for a math/rand.
func newSeed() int64 {
	r, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		panic(fmt.Errorf("failed to read random bytes: %v", err))
	}
	return r.Int64()
}

// GenerateUUID is used to generate a random UUID.
func GenerateUUID() string {
	buf := make([]byte, 16)
	if _, err := crand.Read(buf); err != nil {
		panic(fmt.Errorf("failed to read random bytes: %v", err))
	}

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		buf[0:4],
		buf[4:6],
		buf[6:8],
		buf[8:10],
		buf[10:16])
}

// PoorManUUID generate a uint64 uuid
func PoorManUUID(client bool) (result uint64) {
	result = PoorManUUID2()

	if client {
		result |= 1 //odd for client
	} else {
		result &= math.MaxUint64 - 1 //even for server
	}
	return
}

// PoorManUUID2 doesn't care whether client/server side
func PoorManUUID2() (result uint64) {
	buf := make([]byte, 8)
	rand.Read(buf)

	result = binary.LittleEndian.Uint64(buf)

	if result == 0 {
		result = math.MaxUint64
	}

	return
}
