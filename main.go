package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Always succeed.
// If the git dir can't be found or
// the HEAD can't be read, treat that
// as "not in a repo" and report nothing.

var (
	flagN = flag.Bool("n", false, "omit trailing NL")
	flagS = flag.Bool("s", false, "print spaces around output")
	flagT = flag.Bool("t", false, "set ANSI terminal title")
)

var (
	branch = []byte("ref: refs/heads/")
	ref    = []byte("ref: ")

	spc        = []byte{' '}
	esc   byte = 0x1b
	bel   byte = 0x07
	title      = []byte{esc, ']', '0', ';'}
)

func main() {
	flag.Parse()
	dir, _ := os.Getwd()
	for !isDir(filepath.Join(dir, ".git")) {
		if dir == "/" || dir == "." {
			return
		}
		dir = filepath.Dir(dir)
	}

	b, err := ioutil.ReadFile(filepath.Join(dir, ".git", "HEAD"))
	if err != nil {
		return
	}
	b = bytes.TrimSpace(b)
	switch {
	case bytes.HasPrefix(b, branch):
		b = b[len(branch):]
	case bytes.HasPrefix(b, ref):
		b = b[len(ref):]
	default: // detached head
		b = b[:8] // abbreviate hash
	}

	if *flagT {
		b = append(title, b...)
		b = append(b, bel)
		os.Stdout.Write(b)
		return
	}

	if *flagS {
		b = append(spc, b...)
		b = append(b, spc...)
	}
	if !*flagN {
		b = append(b, '\n')
	}
	fmt.Print(string(b))
}

// isDir returns whether name is affirmatively
// known to exist and to be a directory.
// The two other possibilities,
// "definitely isn't a dir" or "uncertain
// because of an error" return false.
func isDir(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}
