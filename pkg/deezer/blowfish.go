package deezer

import (
	"fmt"
	"io"
	"os"

	"crypto/cipher"

	"golang.org/x/crypto/blowfish"
)

const blowfishKey = "g4el58wc0zvf9na1"
const blowfishIV = "\x00\x01\x02\x03\x04\x05\x06\x07"

const fileChunkSize = 2048

// decryptBlowfish decrypts blowfish data
func decryptBlowfish(key, data []byte) ([]byte, error) {
	block, err := blowfish.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, []byte(blowfishIV))

	decrypted := make([]byte, len(data))
	mode.CryptBlocks(decrypted, data)

	return decrypted, nil
}

// getBlowfishKey calculates the key required to decrypt the
// blowfish-encrypted file
func (track *Track) GetBlowfishKey() []byte {
	hash := MD5Hash([]byte(fmt.Sprintf("%d", track.ID)))
	key := []byte(blowfishKey)

	output := make([]byte, 16)
	for i := 0; i < 16; i++ {
		output[i] = hash[i] ^ hash[i+16] ^ key[i]
	}
	return output
}

// DecryptSongFile decrypts the encrypted chunks of a song downloaded
// from deezer
func DecryptSongFile(key []byte, inputPath, outputPath string) error {
	// open files
	inFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inFile.Close()
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	buf := make([]byte, fileChunkSize)
	n, err := inFile.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}

	for chunk := 0; n > 0; chunk++ {
		// only decrypt every third chunk (including first
		// chunk)
		encrypted := (chunk%3 == 0)

		// only decrypt if encrypted and whole chunk
		if encrypted && n == fileChunkSize {
			buf, err = decryptBlowfish(key, buf)
			if err != nil {
				return err
			}
		}

		// write the chunk back
		n, err = outFile.Write(buf)
		if err != nil {
			return err
		}

		// read next chunk
		n, err = inFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
	}

	return nil
}
