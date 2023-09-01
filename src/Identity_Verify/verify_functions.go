package Identity_Verify

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"os"
)

func GenerateKey(userProvide bool) {
	if userProvide == false {
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			panic(err)
		}

		// x509 serialize
		tempText, _ := x509.MarshalECPrivateKey(privateKey)
		if err != nil {
			panic(err)
		}

		block := pem.Block{
			Type:  "ECDSA private key",
			Bytes: tempText,
		}

		// pem coding
		privateFile, _ := os.Create("private.pem")
		pem.Encode(privateFile, &block)

		// get public key
		publicKey := privateKey.PublicKey

		tempText, _ = x509.MarshalPKIXPublicKey(&publicKey)

		block = pem.Block{
			Type:  "ECDSA public key",
			Bytes: tempText,
		}
		publicFile, _ := os.Create("public.pem")
		pem.Encode(publicFile, &block)

		privateFile.Close()
		publicFile.Close()
	}
	return
}

func Sign(message []byte, privateName string) (rText, sText []byte) {
	file, err := os.Open(privateName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	buf := make([]byte, fileInfo.Size())
	file.Read(buf)

	// pem decode
	block, _ := pem.Decode(buf)

	// x509 deserialize
	privateKey, _ := x509.ParseECPrivateKey(block.Bytes)

	// generate message's hash
	apiHash := sha256.New()
	apiHash.Write(message)
	msgHash := apiHash.Sum(nil)

	r, s, _ := ecdsa.Sign(rand.Reader, privateKey, msgHash)

	// serialize
	rText, _ = r.MarshalText()
	sText, _ = s.MarshalText()
	return
}

func Verify(message, rText, sText []byte, publicName string) bool {
	file, _ := os.Open(publicName)
	defer file.Close()
	fileInfo, _ := file.Stat()
	buf := make([]byte, fileInfo.Size())
	file.Read(buf)

	block, _ := pem.Decode(buf)
	pubInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)
	publicKey := pubInterface.(*ecdsa.PublicKey)

	apiHash := sha256.New()
	apiHash.Write(message)
	msgHash := apiHash.Sum(nil)

	var r, s big.Int
	r.UnmarshalText(rText)
	s.UnmarshalText(sText)
	return ecdsa.Verify(publicKey, msgHash, &r, &s)
}
