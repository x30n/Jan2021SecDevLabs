package main

type PaddingOracle func([]byte) bool

func PaddingOracleDecryptBlock(oracle PaddingOracle, blocks []byte, blockSize int) string {
	if len(blocks)/2 != blockSize {
		// Need to throw an error here
		return ""
	}
	decrypt := make([]byte, blockSize)
	var savedByte byte
	var targetOffset int
	var targetPad int
	for b := 0; b < blockSize; b++ {
		// This bit me for a while. Was simply doing tempBlocks := blocks
		// for each round, but that wasn't reinitializing tempBlocks to
		// the original state of blocks. doing := on a slice is probably
		// actually a pointer to the original slice - need to research
		// golang slice internals...
		tempBlocks := make([]byte, blockSize*2)
		copy(tempBlocks, blocks)

		targetPad = b + 1
		targetOffset = len(blocks) - blockSize - (targetPad)
		savedByte = blocks[targetOffset]

		for j, val := range decrypt {
			if int(val) != 0 {
				tempBlocks[j] = byte(int(tempBlocks[j]) ^ int(val) ^ targetPad)
			}
		}
		hit := false
		for i := 0; i < 256; i++ {
			tempBlocks[targetOffset] = byte(int(savedByte) ^ targetPad ^ i)
			// The check below will fail if the plaintext is legitimately padded with "\x01"
			// Need to think of a better check to accomodate this condition

			if targetOffset == blockSize-1 && tempBlocks[targetOffset] == savedByte {
				continue
			}

			if oracle(tempBlocks) {
				hit = true
				decrypt[targetOffset] = byte(i)
				break
			}
		}
		if targetOffset == blockSize-1 && hit == false {
			decrypt[targetOffset] = byte(1)
		}
	}
	return string(decrypt)
}
