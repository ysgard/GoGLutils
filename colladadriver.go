package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stdout, "Syntax: colladadriver.exe <colladafile>\n")
		return
	}
	fname := os.Args[1]
	c, err := ReadColladaFile(fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not process Collada file!\n")
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	c.Debug()
}
