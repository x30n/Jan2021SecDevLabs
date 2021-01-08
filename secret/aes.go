package main

import (
	"bytes"
	"crypto/aes"
	"encoding/binary"
)

func ECBDecrypt(key []byte, ciphertext []byte) []byte {

	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err)
	}

	bs := block.BlockSize()
	if len(ciphertext)%bs != 0 {
		panic("Need a multiple of the blocksize")
	}

	plaintext := make([]byte, len(ciphertext))
	var pt []byte
	pt = plaintext
	for len(ciphertext) > 0 {
		block.Decrypt(plaintext, ciphertext)
		ciphertext = ciphertext[bs:]
		plaintext = plaintext[bs:]
	}
	return pt
}

func ECBEncrypt(key []byte, plaintext []byte) []byte {

	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err)
	}

	bs := block.BlockSize()
	if len(plaintext)%bs != 0 {
		panic("Need a multiple of the blocksize")
	}

	ciphertext := make([]byte, len(plaintext))
	var ct []byte
	ct = ciphertext
	for len(plaintext) > 0 {
		block.Encrypt(ciphertext, plaintext)
		plaintext = plaintext[bs:]
		ciphertext = ciphertext[bs:]

	}
	return ct
}

func CBCEncrypt(key []byte, plaintext []byte, iv []byte) []byte {
	var ct []byte
	x := iv
	for i := 0; i <= len(plaintext)-len(key); i += len(key) {
		temp_block, err := FixedXOR(plaintext[i:i+len(key)], x)
		if err != nil {
			panic(err)
		}
		temp_block = ECBEncrypt(key, temp_block)
		ct = append(ct, temp_block...)
		x = temp_block
	}
	return ct
}

func CBCDecrypt(key []byte, ciphertext []byte, iv []byte) []byte {
	var pt []byte
	x := iv
	for i := 0; i <= len(ciphertext)-len(key); i += len(key) {
		next_x := ciphertext[i : i+len(key)]
		temp_block := ECBDecrypt(key, ciphertext[i:i+len(key)])
		t, err := FixedXOR(temp_block, x)
		if err != nil {
			panic(err)
		}
		pt = append(pt, t...)
		x = next_x
	}
	return pt
}

func CTR(message []byte, key []byte, nonce uint64) ([]byte, error) {
	var bcount uint64
	buf := new(bytes.Buffer)
	var out []byte
	err := binary.Write(buf, binary.LittleEndian, nonce)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, bcount)
	if err != nil {
		return nil, err
	}
	keystream := ECBEncrypt(key, buf.Bytes())
	j := 0
	for i := 0; i < len(message); i++ {
		if j == len(buf.Bytes()) {
			bcount++
			buf = new(bytes.Buffer)
			err := binary.Write(buf, binary.LittleEndian, nonce)
			if err != nil {
				return nil, err
			}
			err = binary.Write(buf, binary.LittleEndian, bcount)
			if err != nil {
				return nil, err
			}
			j = 0
			keystream = ECBEncrypt(key, buf.Bytes())
		}
		out = append(out, byte(message[i]^keystream[j]))
		j++
	}
	return []byte(out), nil
}
