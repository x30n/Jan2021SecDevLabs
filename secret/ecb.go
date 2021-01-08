package main

import (
	"strings"
)

func FindPrefixSize(oracle Oracle) int {
	blockSize := FindBlockSize(oracle)
	baselineSize := len(oracle("A"))
	dups := oracle(strings.Repeat("A", baselineSize*2))
	var dblock string
	blocks := ChunkStr(string(dups), blockSize)
	controlledIdx := 0
	for idx, block := range blocks {
		if block == dblock {
			controlledIdx = idx - 1
			testBlocks := blocks[controlledIdx : (baselineSize/blockSize)+1]
			// CountDuplicateBlocks ensures that the dup blocks we've found
			// are actually the padding we supplied in case the prefix also
			// contains duplicate blocks.
			// NOTE - If the prefix is made up of the same bytes as our
			// padding ("A"), this will fail.
			// TODO - Could run through FindPrefixSize twice with different
			// padding bytes to account for this.
			if CountDuplicateBlocks(testBlocks) == len(testBlocks) {
				break
			}
		}
		dblock = block
	}
	// No Prefix
	if controlledIdx == 0 {
		return 0
	}
	for i := 0; i < blockSize; i++ {
		enc := ChunkStr(string(oracle(strings.Repeat("A", (blockSize*2)+i))), blockSize)
		if enc[controlledIdx] == enc[controlledIdx+1] {
			return (controlledIdx * blockSize) - i
		}
	}
	return controlledIdx
}

func ECBOracleDecryption(oracle Oracle) []byte {
	var decrypted string
	prefixLen := FindPrefixSize(oracle)
	baselineSize := len(oracle("A")) // Get baseline length feeding 1 byte into oracle
	blockSize := FindBlockSize(oracle)
	maxPadLen := baselineSize - prefixLen
	offsetStart := baselineSize - blockSize
	offsetEnd := baselineSize
	for i := 1; i < maxPadLen; i++ {
		pad := strings.Repeat("A", maxPadLen-i)
		savedPattern := oracle(pad)[offsetStart:offsetEnd]
		patternMap := make(map[string]string)
		for j := 0; j < 256; j++ {
			tempString := pad + decrypted + string(byte(j))
			mapKey := string(oracle(tempString)[offsetStart:offsetEnd]) // Save each iteration of encrypted PAD + bytes 0-255 for lookup
			mapVal := string(byte(j))                                   // Map values 0 - 255 to encrypted equiv for lookup
			patternMap[mapKey] = mapVal
		}
		decrypted += patternMap[string(savedPattern)]
	}
	return []byte(decrypted)
}
