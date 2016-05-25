package main

import (
	"crypto/rand"
	"crypto/rsa"
	//"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/alecthomas/kingpin"
)

var encryptCommand = kingpin.Command("encrypt", "encrypts values for the .travis.yml").Action(func(ctx *kingpin.ParseContext) error {
	if len(ctx.Elements) < 2 {
		kingpin.Usage()
		return nil
	}

	err := auth()
	if err != nil {
		return err
	}

	s := slug()
	repo, _, err := client.Repositories.GetFromSlug(s)
	if err != nil {
		return err
	}

	req, err := client.NewRequest("GET", "/repos/"+fmt.Sprint(repo.Id)+"/key", nil, nil)
	if err != nil {
		return err
	}

	var pubkey struct {
		Key         string `json:"key"`
		FingerPrint string `json:"fingerprint"`
	}
	_, err = client.Do(req, &pubkey)
	if err != nil {
		return err
	}

	block, _ := pem.Decode([]byte(pubkey.Key))
	if err != nil {
		return err
	}

	cert, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	rsaPublicKey := cert.(*rsa.PublicKey)

	s = ""
	for _, v := range ctx.Elements[1:] {
		if s != "" {
			s += " "
		}
		s += *v.Value
	}
	content := []byte(s)

	b, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, content)
	if err != nil {
		return err
	}

	s = "  secure: " + base64.StdEncoding.EncodeToString(b)
	fmt.Println(s)
	return nil
})
var encryptArg = encryptCommand.Arg("data", "data to encrypt").Strings()
