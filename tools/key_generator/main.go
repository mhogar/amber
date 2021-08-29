package main

import (
	"authserver/common"
	"authserver/config"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"log"
	"os"
	"path"
)

func main() {
	err := config.InitConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	//parse the flags
	name := flag.String("name", "key", "The name of key. Will create files <name>.private.pem and <name>.public.pem.")
	overwrite := flag.Bool("overwrite", false, "If the key files should be overwriten if they already exist.")
	flag.Parse()

	//run the generator
	err = Run(*name, *overwrite)
	if err != nil {
		log.Fatal(err)
	}
}

type createPEMBlockFunc func(key interface{}) (*pem.Block, error)

// Run runs the key generator with the given inputs.
func Run(name string, overwrite bool) error {
	reader := rand.Reader

	//generate the key
	key, err := rsa.GenerateKey(reader, 2048)
	if err != nil {
		return common.ChainError("error generating key", err)
	}

	//save the private key
	err = savePrivateKey(name, overwrite, key)
	if err != nil {
		return common.ChainError("error saving private key", err)
	}

	//save the public key
	err = savePublicKey(name, overwrite, &key.PublicKey)
	if err != nil {
		return common.ChainError("error saving public key", err)
	}

	return nil
}

func saveKey(name string, overwrite bool, key interface{}, createPEMBlock createPEMBlockFunc) error {
	filename := path.Join(config.GetAppRoot(), "static", "keys", name+".pem")

	//check if the file already exists if we don't want to overwrite it
	if !overwrite {
		_, err := os.Stat(filename)
		if !os.IsNotExist(err) {
			return errors.New("file already exists")
		}
	}

	//create the file
	file, err := os.Create(filename)
	if err != nil {
		return common.ChainError("error creating file", err)
	}
	defer file.Close()

	//create the pem block
	block, err := createPEMBlock(key)
	if err != nil {
		return common.ChainError("error creating pem block", err)
	}

	//encode the block to the file
	err = pem.Encode(file, block)
	if err != nil {
		return common.ChainError("error encoding pem block to file", err)
	}

	log.Println("Created file:", filename)
	return nil
}

func savePrivateKey(name string, overwrite bool, key *rsa.PrivateKey) error {
	return saveKey(name+".private", overwrite, key, func(key interface{}) (*pem.Block, error) {
		//assert key type
		privateKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("key must be an RSA Private Key")
		}

		//create the pem block
		block := &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		}

		return block, nil
	})
}

func savePublicKey(name string, overwrite bool, key *rsa.PublicKey) error {
	return saveKey(name+".public", overwrite, key, func(key interface{}) (*pem.Block, error) {
		//assert key type
		publicKey, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("key must be an RSA Public Key")
		}

		//create the pem block
		block := &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(publicKey),
		}

		return block, nil
	})
}
