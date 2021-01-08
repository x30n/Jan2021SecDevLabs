package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func PrintHexDecode(s string) {
	str, err := hex.DecodeString(s)
	if err != nil {
		fmt.Println("encoding.hex.DecodeString Failed.")
	}
	fmt.Println(string(str))
}

func FixedXOR(b1 []byte, b2 []byte) ([]byte, error) {
	var ret []byte
	if len(b1) != len(b2) {
		return ret, errors.New("Arguments b1 and b2 must be same length.")
	}
	for i := 0; i < len(b1); i++ {
		ret = append(ret, (b1[i] ^ b2[i]))
	}
	return ret, nil
}

func XORWithKey(ct []byte, key []byte) []byte {
	var ret []byte
	ki := 0
	for i := 0; i < len(ct); i++ {
		ret = append(ret, (ct[i] ^ key[ki]))
		if len(key) > 1 {
			if ki < (len(key) - 1) {
				ki++
			} else {
				ki = 0
			}
		}
	}
	return ret
}

func LangScore(s string) int {
	chars := "etaoin shrdlu"
	score := 0
	s = strings.ToLower(s)
	for _, c := range chars {
		score += strings.Count(s, string(c))
	}
	return score
}

func BruteXORByte(ctext []byte) (string, int, string) {
	score := 0
	key := ""
	ptext := ""
	for i := 0; i <= 255; i++ {
		tptext := string(XORWithKey(ctext, []byte(string(i))))
		tscore := LangScore(tptext)
		if tscore > score {
			score = tscore
			key = string(i)
			ptext = tptext
		}
	}
	return key, score, ptext
}

func GetHamming(str1 string, str2 string) int {
	c := 0
	ba1 := []byte(str1)
	ba2 := []byte(str2)

	for i := 0; i < len(str1); i++ {
		c += strings.Count(strconv.FormatInt(int64(ba1[i]^ba2[i]), 2), "1")
	}
	return c
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func sumSlices(x []float64, y []float64) float64 {

	totalx := 0.0
	for _, valuex := range x {
		totalx += valuex
	}

	totaly := 0.0
	for _, valuey := range y {
		totaly += valuey
	}

	return totalx + totaly
}

func sumSlice(x []float64) float64 {

	totalx := 0.0
	for _, valuex := range x {
		totalx += valuex
	}
	return totalx
}

func findKeySize(cstr string, minlen float64, maxlen float64) float64 {
	var ksz, low_hamm float64
	NUMSAMPLES := 4
	var h_avg float64

	for keysize := minlen; keysize <= maxlen; keysize++ {
		var temp_hamms []float64
		var slices []string
		slices = append(slices, cstr[0:int(keysize)])
		for c := 1; c <= NUMSAMPLES; c++ {
			slices = append(slices, cstr[int(keysize)*c:int(keysize)*(c+1)])
		}

		done := map[string]bool{}

		for _, s := range slices {
			for _, x := range slices {
				if s == x {
					continue
				}
				if !(done[s]) && !(done[x]) {
					temp_hamms = append(temp_hamms, float64(GetHamming(s, x)))
				}
			}
			done[s] = true
		}
		h_avg = ((sumSlice(temp_hamms) / float64(len(temp_hamms))) / keysize)

		if (h_avg < low_hamm) || (low_hamm == 0) {
			low_hamm = h_avg
			ksz = keysize
		}
	}
	return ksz
}

func ChunkStr(ctxt string, bsize int) []string {
	var ret []string
	for i := 0; i <= len(ctxt)-bsize; i += bsize {
		ret = append(ret, ctxt[0+i:bsize+i])
	}
	return ret
}

func TransposeBlocks(blocks []string) []string {
	bsize := len(blocks[0])
	var transposed_blocks []string
	for i := 0; i < bsize; i++ {
		tstr := ""
		for _, block := range blocks {
			tstr += string(block[i])
		}
		transposed_blocks = append(transposed_blocks, tstr)
	}
	return transposed_blocks
}

func CrackXORBlocks(blocks []string) string {
	k := ""
	for _, b := range blocks {
		tk, _, _ := BruteXORByte([]byte(b))
		k += tk
	}
	return k
}

func DetectECB(ciphertext []string) bool {
	ctmap := make(map[string]bool, 0)
	for _, block := range ciphertext {
		ctmap[block] = true
	}
	/*
		Duplicate idx won't be added to map, so size of map will be different
		from size of slice if duplicate blocks exist.
	*/
	return len(ciphertext) != len(ctmap)
}

func PKCS7(blocks []byte, blocksize int) []byte {

	remainder := len(blocks) % blocksize
	if remainder == 0 {
		// If blocks is evenly divisable by blocksize,
		// return blocks + 1 full block of padding
		return append(blocks, bytes.Repeat([]byte{byte(blocksize)}, blocksize)...)
	}
	pad := blocksize - remainder
	return append(blocks, bytes.Repeat([]byte{byte(pad)}, pad)...)
}

func PKCS7UnPad(blocks []byte) ([]byte, error) {
	lastblock := int(blocks[len(blocks)-1])
	if VerifyPKCS7(blocks) {
		return blocks[:len(blocks)-lastblock], nil
	}
	return nil, errors.New("Bad PKCS#7 padding")
}

func VerifyPKCS7(blocks []byte) bool {
	lastbyte := uint(blocks[len(blocks)-1])
	if lastbyte > uint(len(blocks)-1) || lastbyte <= 0 {
		return false
	}
	for i := uint(len(blocks)) - lastbyte; i < uint(len(blocks)); i++ {
		if uint(blocks[i]) != lastbyte {
			return false
		}
	}
	return true
}

func GenerateRandomBytes(n uint) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func ConcatSliceOfStrings(input []string) string {
	outString := ""
	for _, str := range input {
		outString += str
	}
	return outString
}

type Oracle func(string) []byte

/*
* FindBlockSize takes an oracle function as an argument and attempts to determine
* the blocksize. The oracle must consume a single string argument plaintext and
* return a []byte containing the encrypted ciphertext.
* Return Value: FindBlockSize will return the blocksize or 0 if there's an error.
 */
func FindBlockSize(oracle Oracle) int {
	l := 0
	baseline := len(oracle("A"))
	for i := 1; i <= 100; i++ {
		l = len(oracle(strings.Repeat("A", i)))
		if l > baseline {
			return l - baseline
		}
	}
	return 0
}

func FindAddedPlaintextSize(oracle Oracle, blocksize int) int {
	l := 0
	baseline := len(oracle(strings.Repeat("A", blocksize)))
	for i := blocksize; i > 0; i-- {
		l = len(oracle(strings.Repeat("A", i)))
		if l < baseline {
			return l - i
		}
	}
	return 0
}

func CountDuplicateBlocks(blocks []string) int {
	c := 0
	firstBlock := blocks[0]
	for _, v := range blocks {
		if v == firstBlock {
			c++
		}
	}
	return c
}

func PrettyPrintBlocks(blocks []byte, blocksize int) {
	var savedidx int
	for i := 0; i < len(blocks); i += blocksize {
		fmt.Printf("[%16x]", blocks[i:i+blocksize-1])
		savedidx = i
	}
	fmt.Printf("[%16x]\n", blocks[savedidx:])
}
