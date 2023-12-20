package main

import (
	"fmt"
)

// Initial Permutation table for DES
var initialPermutationTable = [64]int{
	58, 50, 42, 34, 26, 18, 10, 2,
	60, 52, 44, 36, 28, 20, 12, 4,
	62, 54, 46, 38, 30, 22, 14, 6,
	64, 56, 48, 40, 32, 24, 16, 8,
	57, 49, 41, 33, 25, 17, 9, 1,
	59, 51, 43, 35, 27, 19, 11, 3,
	61, 53, 45, 37, 29, 21, 13, 5,
	63, 55, 47, 39, 31, 23, 15, 7,
}

// InitialPermutation performs the initial permutation on a 64-bit block
func InitialPermutation(block uint64) uint64 {
	result := uint64(0)
	for i := 0; i < 64; i++ {
		// Get the bit at the i-th position in the block
		bit := (block >> (64 - initialPermutationTable[i])) & 1
		// Set the corresponding bit in the result
		result |= (bit << (63 - i))
	}
	return result
}

// Final Permutation table for DES
var finalPermutationTable = [64]int{
	40, 8, 48, 16, 56, 24, 64, 32,
	39, 7, 47, 15, 55, 23, 63, 31,
	38, 6, 46, 14, 54, 22, 62, 30,
	37, 5, 45, 13, 53, 21, 61, 29,
	36, 4, 44, 12, 52, 20, 60, 28,
	35, 3, 43, 11, 51, 19, 59, 27,
	34, 2, 42, 10, 50, 18, 58, 26,
	33, 1, 41, 9, 49, 17, 57, 25,
}

// FinalPermutation performs the final permutation on a 64-bit block
func FinalPermutation(block uint64) uint64 {
	result := uint64(0)
	for i := 0; i < 64; i++ {
		// Get the bit at the i-th position in the block
		bit := (block >> (64 - finalPermutationTable[i])) & 1
		// Set the corresponding bit in the result
		result |= (bit << (63 - i))
	}
	return result
}

// Expansion permutation table for DES
var expansionTable = [48]int{
	32, 1, 2, 3, 4, 5,
	4, 5, 6, 7, 8, 9,
	8, 9, 10, 11, 12, 13,
	12, 13, 14, 15, 16, 17,
	16, 17, 18, 19, 20, 21,
	20, 21, 22, 23, 24, 25,
	24, 25, 26, 27, 28, 29,
	28, 29, 30, 31, 32, 1,
}

// Example 56-bit DES key
var desKey uint64 = 0x133457799BBCDFF1

func main() {

	// Generate the 16 subkeys from the DES key
	subkeys := GenerateSubKeys(desKey)

	fmt.Printf("Original Key: 0x%016X\n", desKey)
	fmt.Println("Generated Subkeys:")
	for i, subkey := range subkeys {
		fmt.Printf("Key Round %2d: 0x%012X\n", i+1, subkey)
	}

	// input block
	//inputBlock := uint64(0x0123456789ABCDEF)
	inputBlock := uint64(0xAA238F292D2D04A1)

	// Perform the initial permutation
	result := InitialPermutation(inputBlock)

	// Display the result in hexadecimal
	fmt.Printf("Input Block:  0x%016X\n", inputBlock)
	fmt.Printf("Initial Permutated Block: 0x%016X\n", result)

	for i, subkey := range subkeys {
		result = OneRoundDES(result, subkey)
		fmt.Printf("DES Round %2d: 0x%012X\n", i+1, result)
	}

	// Perform the final permutation
	result = FinalPermutation(result)

	// Display the result in hexadecimal
	fmt.Printf("Input Block:  0x%016X\n", inputBlock)
	fmt.Printf("Final Permutated Block: 0x%016X\n", result)

}

func OneRoundDES(data, key uint64) uint64 {
	// Expansion permutation
	expandedData := uint64(0)
	for i := 0; i < 48; i++ {
		bit := (data >> (32 - expansionTable[i])) & 1
		expandedData |= (bit << (47 - i))
	}

	// XOR with round key
	result := expandedData ^ key

	// Example: Apply S-boxes and other operations here

	return result
}

// Permutation table for key permutation choice 1 (PC1)
var pc1Table = [56]int{
	57, 49, 41, 33, 25, 17, 9,
	1, 58, 50, 42, 34, 26, 18,
	10, 2, 59, 51, 43, 35, 27,
	19, 11, 3, 60, 52, 44, 36,
	63, 55, 47, 39, 31, 23, 15,
	7, 62, 54, 46, 38, 30, 22,
	14, 6, 61, 53, 45, 37, 29,
	21, 13, 5, 28, 20, 12, 4,
}

// Permutation table for key permutation choice 2 (PC2)
var pc2Table = [48]int{
	14, 17, 11, 24, 1, 5,
	3, 28, 15, 6, 21, 10,
	23, 19, 12, 4, 26, 8,
	16, 7, 27, 20, 13, 2,
	41, 52, 31, 37, 47, 55,
	30, 40, 51, 45, 33, 48,
	44, 49, 39, 56, 34, 53,
	46, 42, 50, 36, 29, 32,
}

// Left shift schedule for key generation
var leftShifts = [16]int{
	1, 1, 2, 2,
	2, 2, 2, 2,
	1, 2, 2, 2,
	2, 2, 2, 1,
}

// GenerateSubKeys generates 16 subkeys for DES
func GenerateSubKeys(key uint64) [16]uint64 {
	subkeys := [16]uint64{}
	// Perform PC1 permutation on the key
	keyPermuted := uint64(0)
	for i := 0; i < 56; i++ {
		bit := (key >> (64 - pc1Table[i])) & 1
		keyPermuted |= (bit << (55 - i))
	}
	// Split the 56-bit key into two 28-bit halves
	leftHalf := keyPermuted >> 28
	rightHalf := keyPermuted & 0x0FFFFFFF
	for i := 0; i < 16; i++ {
		// Perform left shifts on the halves
		leftHalf = ((leftHalf << leftShifts[i]) | (leftHalf >> (28 - leftShifts[i]))) & 0x0FFFFFFF
		rightHalf = ((rightHalf << leftShifts[i]) | (rightHalf >> (28 - leftShifts[i]))) & 0x0FFFFFFF
		// Combine the halves and perform PC2 permutation
		combinedHalf := (leftHalf << 28) | rightHalf
		subkey := uint64(0)
		for j := 0; j < 48; j++ {
			bit := (combinedHalf >> (56 - pc2Table[j])) & 1
			subkey |= (bit << (47 - j))
		}
		// Store the subkey
		subkeys[i] = subkey
	}
	return subkeys
}
