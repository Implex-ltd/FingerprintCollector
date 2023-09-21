package fingerprint

import (
	"encoding/base64"
)

/*
- Encryption: Fp  -> b64 -> Xor  -> b64
- Decryption: b64 -> Xor -> B64  -> Fp
*/
func Decrypt(data, key string) (string, error) {
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	decryptedData := make([]byte, len(decodedData))
	keyBytes := []byte(key)

	for i := 0; i < len(decodedData); i++ {
		decryptedData[i] = decodedData[i] ^ keyBytes[i%len(keyBytes)]
	}

	return string(decryptedData), nil
}

func Encrypt(data, key string) (string, error) {
    dataBytes := []byte(data)
    keyBytes := []byte(key)

    encryptedData := make([]byte, len(dataBytes))

    for i := 0; i < len(dataBytes); i++ {
        encryptedData[i] = dataBytes[i] ^ keyBytes[i%len(keyBytes)]
    }

    encodedString := base64.StdEncoding.EncodeToString(encryptedData)
    return encodedString, nil
}