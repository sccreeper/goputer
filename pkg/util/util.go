package util

import (
	"math/rand"
	"sccreeper/govm/pkg/constants"
)

func CheckError(err error) {

	if err != nil {
		panic(err)
	}

}

var Check = CheckError

func SliceRemove[T string | byte](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

// Sees if a splice contains a specified X
func SliceContains[T string | constants.Instruction | constants.Interrupt](splice []T, search_value T) bool {

	for _, v := range splice {

		if v == search_value {
			return true
		}

	}

	return false

}

// Splits a slice into chunks
func SliceChunks[T any](slice []T, chunk_size int) [][]T {

	chunks := make([][]T, 0)

	for i := 0; i < len(slice); i += chunk_size {

		end := i + chunk_size

		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])

	}

	return chunks

}

func RandomNumber[T uint8 | int | uint32](min T, max T) T {

	return T(rand.Intn(int(max-min+1)) + int(min))

}
