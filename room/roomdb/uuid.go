package roomdb

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateShortID() string {
	// two letters + two numbers
	// == 26*26*10*10
	// == 369664 combinations
	const alphabet = "abcdefghijklmnopqrstuvwxyz"

	firstLetterIndex := randomNumber(0, len(alphabet)-1)
	secondLetterIndex := randomNumber(0, len(alphabet)-1)
	firstNumber := randomNumber(0, 9)
	secondNumber := randomNumber(0, 9)

	return fmt.Sprintf("%s%s%d%d", string(alphabet[firstLetterIndex]), string(alphabet[secondLetterIndex]), firstNumber, secondNumber)
}

func randomNumber(min int, max int) int {
	i, err := rand.Int(rand.Reader, big.NewInt(int64(max+min)))
	if err != nil {
		panic("failed to generate random number")
	}

	return int(i.Int64()) + min
}
