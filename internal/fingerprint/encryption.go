package fingerprint

import "encoding/base64"

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

	decodedString, err := base64.StdEncoding.DecodeString(string(decryptedData))
	if err != nil {
		return "", err
	}

	return string(decodedString), nil
}
