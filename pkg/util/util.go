package util

import (
	"fmt"
	"math/rand"

	"golang.org/x/exp/constraints"
)

func CheckError(err error) {

	if err != nil {
		panic(err)
	}

}

var Check = CheckError

// Removes an item from a slice and keeps the order
func SliceRemove[T string | byte](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
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

// Generates random number
func RandomNumber[T uint8 | int | uint32](min T, max T) T {

	return T(rand.Intn(int(max-min+1)) + int(min))

}

func AllEqualToX[T uint32 | byte](splice []T, checkValue T) bool {

	for _, v := range splice {

		if v != checkValue {
			return false
		}

	}

	return true

}

func ConvertHex[T int | uint32 | uint64](i T) string {

	return fmt.Sprintf("0x"+"%08X", i)

}

func CombineMap[Map ~map[K]V, K comparable, V any](a map[K]V, b map[K]V) (result map[K]V) {
	result = a

	for k, v := range b {
		result[k] = v
	}

	return
}

type Number interface {
	constraints.Integer | constraints.Float
}

func Clamp[T Number](number T, min T, max T) T {
	if number < min {
		return min
	} else if number > max {
		return max
	} else {
		return number
	}
}
