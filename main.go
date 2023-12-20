package main

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"fmt"
)

func main() {
	// Define the 8-byte key and plaintext
	key := []byte("12345678")
	plaintext := []byte("Hello123")

	// Create a new DES block cipher with the provided key
	block, err := des.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Initial Permutation (IP) table
	ipTable := [64]byte{
		58, 50, 42, 34, 26, 18, 10, 2,
		60, 52, 44, 36, 28, 20, 12, 4,
		62, 54, 46, 38, 30, 22, 14, 6,
		64, 56, 48, 40, 32, 24, 16, 8,
		57, 49, 41, 33, 25, 17, 9, 1,
		59, 51, 43, 35, 27, 19, 11, 3,
		61, 53, 45, 37, 29, 21, 13, 5,
		63, 55, 47, 39, 31, 23, 15, 7,
	}

	// Permute the plaintext using the initial permutation (IP) table
	initialPermutation(plaintext, ipTable)

	// Create an initialization vector (IV) for Cipher Block Chaining (CBC) mode
	iv := make([]byte, block.BlockSize())

	// Create a CBC mode encrypter
	encrypter := cipher.NewCBCEncrypter(block, iv)

	// Encrypt the plaintext
	ciphertext := make([]byte, len(plaintext))
	encrypter.CryptBlocks(ciphertext, plaintext)

	// Print the ciphertext as hex
	fmt.Printf("Ciphertext (hex): %s\n", hex.EncodeToString(ciphertext))
}

// initialPermutation permutes the input data according to the given permutation table
func initialPermutation(data []byte, table [64]byte) {
	result := make([]byte, len(data))
	for i, bit := range table {
		result[i] = data[bit-1]
	}
	copy(data, result)
}
