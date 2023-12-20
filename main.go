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

	// input data block
	inputBlock := uint64(0x0123456789ABCDEF)
	//inputBlock := uint64(0xA88028A888800820)

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
	xorResult := expandedData ^ key

	// Apply the S-boxes
	sBoxResult := uint64(0)
	for i := 0; i < 8; i++ {
		// Extract 6 bits from xorResult for each S-box
		sixBits := (xorResult >> (42 - i*6)) & 0x3F
		// Get the row and column from the six bits
		row := ((sixBits >> 4) & 0x2) | (sixBits & 0x1)
		col := (sixBits >> 1) & 0xF
		// Look up the value in the S-box and append it to the result
		sBoxResult |= uint64(sBoxes[i][row][col]) << (28 - i*4)
	}

	// Perform the P permutation
	pPermutation := uint64(0)
	for i := 0; i < 32; i++ {
		bit := (sBoxResult >> (32 - i)) & 1
		pPermutation |= (bit << (31 - i))
	}

	// XOR the result with the left half of the input data
	result := data>>32 ^ pPermutation

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

// S-boxes for DES
var sBoxes = [8][4][16]int{
	{
		// S-box 1
		{14, 4, 13, 1, 2, 15, 11, 8, 3, 10, 6, 12, 5, 9, 0, 7},
		{0, 15, 7, 4, 14, 2, 13, 1, 10, 6, 12, 11, 9, 5, 3, 8},
		{4, 1, 14, 8, 13, 6, 2, 11, 15, 12, 9, 7, 3, 10, 5, 0},
		{15, 12, 8, 2, 4, 9, 1, 7, 5, 11, 3, 14, 10, 0, 6, 13},
	},
	{
		// S-box 2
		{15, 1, 8, 14, 6, 11, 3, 4, 9, 7, 2, 13, 12, 0, 5, 10},
		{3, 13, 4, 7, 15, 2, 8, 14, 12, 0, 1, 10, 6, 9, 11, 5},
		{0, 14, 7, 11, 10, 4, 13, 1, 5, 8, 12, 6, 9, 3, 2, 15},
		{13, 8, 10, 1, 3, 15, 4, 2, 11, 6, 7, 12, 0, 5, 14, 9},
	},
	{
		// S-box 3
		{10, 0, 9, 14, 6, 3, 15, 5, 1, 13, 12, 7, 11, 4, 2, 8},
		{13, 7, 0, 9, 3, 4, 6, 10, 2, 8, 5, 14, 12, 11, 15, 1},
		{13, 6, 4, 9, 8, 15, 3, 0, 11, 1, 2, 12, 5, 10, 14, 7},
		{1, 10, 13, 0, 6, 9, 8, 7, 4, 15, 14, 3, 11, 5, 2, 12},
	},
	{
		// S-box 4
		{7, 13, 14, 3, 0, 6, 9, 10, 1, 2, 8, 5, 11, 12, 4, 15},
		{13, 8, 11, 5, 6, 15, 0, 3, 4, 7, 2, 12, 1, 10, 14, 9},
		{10, 6, 9, 0, 12, 11, 7, 13, 15, 1, 3, 14, 5, 2, 8, 4},
		{3, 15, 0, 6, 10, 1, 13, 8, 9, 4, 5, 11, 12, 7, 2, 14},
	},
	{
		// S-box 5
		{2, 12, 4, 1, 7, 10, 11, 6, 8, 5, 3, 15, 13, 0, 14, 9},
		{14, 11, 2, 12, 4, 7, 13, 1, 5, 0, 15, 10, 3, 9, 8, 6},
		{4, 2, 1, 11, 10, 13, 7, 8, 15, 9, 12, 5, 6, 3, 0, 14},
		{11, 8, 12, 7, 1, 14, 2, 13, 6, 15, 0, 9, 10, 4, 5, 3},
	},
	{
		// S-box 6
		{12, 1, 10, 15, 9, 2, 6, 8, 0, 13, 3, 4, 14, 7, 5, 11},
		{10, 15, 4, 2, 7, 12, 9, 5, 6, 1, 13, 14, 0, 11, 3, 8},
		{9, 14, 15, 5, 2, 8, 12, 3, 7, 0, 4, 10, 1, 13, 11, 6},
		{4, 3, 2, 12, 9, 5, 15, 10, 11, 14, 1, 7, 6, 0, 8, 13},
	},
	{
		// S-box 7
		{4, 11, 2, 14, 15, 0, 8, 13, 3, 12, 9, 7, 5, 10, 6, 1},
		{13, 0, 11, 7, 4, 9, 1, 10, 14, 3, 5, 12, 2, 15, 8, 6},
		{1, 4, 11, 13, 12, 3, 7, 14, 10, 15, 6, 8, 0, 5, 9, 2},
		{6, 11, 13, 8, 1, 4, 10, 7, 9, 5, 0, 15, 14, 2, 3, 12},
	},
	{
		// S-box 8
		{13, 2, 8, 4, 6, 15, 11, 1, 10, 9, 3, 14, 5, 0, 12, 7},
		{1, 15, 13, 8, 10, 3, 7, 4, 12, 5, 6, 11, 0, 14, 9, 2},
		{7, 11, 4, 1, 9, 12, 14, 2, 0, 6, 10, 13, 15, 3, 5, 8},
		{2, 1, 14, 7, 4, 10, 8, 13, 15, 12, 9, 0, 3, 5, 6, 11},
	},
}
