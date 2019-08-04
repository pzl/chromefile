package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pzl/chromefile/snss"
	log "github.com/sirupsen/logrus"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	//log.SetFormatter()
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "arugment required: file to parse")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	check(err)
	defer f.Close()

	parse(f)
}

func parse(f io.Reader) {
	ver, err := snss.FileInfo(f)
	check(err)
	fmt.Printf("File version: %d\n", ver)

	for err = snss.ReadCommand(f); err == nil; err = snss.ReadCommand(f) {
	}
	if err != nil && err != io.EOF {
		panic(err)
	}

}
