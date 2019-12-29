package deezer

import "crypto/aes"

// ECB implementation since golang technically doesn't support it
// based on https://stackoverflow.com/questions/24072026/golang-aes-ecb-encryption
// input data is automatically zero-padded
func ECB(key, data []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	size := cipher.BlockSize()
	for len(data)%size != 0 {
		data = append(data, '\x00')
	}

	encrypted := make([]byte, len(data))
	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		cipher.Encrypt(encrypted[bs:be], data[bs:be])
	}

	return encrypted, nil
}
