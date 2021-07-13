package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/mehix/sshtunnel/pkg/tunnel"
	"golang.org/x/crypto/ssh"
)

const secretFile = ".secret"

func init() {
	flag.Parse()

	if flag.NArg() != 3 {
		log.Fatalln("Not enough arguments")
	}
}

func main() {

	srvr := flag.Arg(0)
	pkFile := flag.Arg(1)
	dest := flag.Arg(2)

	t := tunnel.NewSSHTunnel(srvr, pkAuth(pkFile), dest)

	t.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)

	if err := t.Start(); err != nil {
		log.Println(err)
	}
}

func pkAuth(f string) ssh.AuthMethod {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		log.Println(err)
		return nil
	}

	s, _ := ioutil.ReadFile(secretFile)

	signer, err := ssh.ParsePrivateKeyWithPassphrase(b, s)
	if err != nil {
		log.Println(err)
		return nil
	}

	return ssh.PublicKeys(signer)
}
