package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type dirs []string

type patternWalker struct {
	pats        []string
	matchSuffix bool
	skipDirs    dirs
}

func (d *dirs) String() string {
	var list []string
	for _, c := range *d {
		list = append(list, c)
	}

	return fmt.Sprintf("[<%s>]", strings.Join(list, " "))
}

func (d *dirs) Set(value string) error {
	comps := strings.Split(value, ",")
	for _, c := range comps {
		*d = append(*d, c)
	}
	return nil
}

const (
	usage = `The walk utility recursively descends the directory tree for the specified directory
and prints files that match the specified substring patterns.

Usage:
  walk <directory> [pat1 pat2 ...]

Flags:`
)

func help() {
	fmt.Fprintln(os.Stderr, usage)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	var skipDirs dirs
	// only flag, -x, treats patterns as file extensions (uses HasPrefix vs Contains)
	matchSuffix := flag.Bool("x", false, "treat patterns as file extensions")
	flag.Var(&skipDirs, "s", "skip this directory")

	flag.Usage = help
	flag.Parse()
	args := flag.Args()

	// at the very least, user must specify the directory
	if len(args) < 1 {
		help()
	}

	p := patternWalker{pats: args[1:], matchSuffix: *matchSuffix, skipDirs: skipDirs}
	if err := filepath.Walk(path.Clean(args[0]), p.Walker); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Walker defines the walker function that's called for each subdirectory and file traversed.
func (w *patternWalker) Walker(full string, info os.FileInfo, err error) error {
	if info == nil {
		return err
	}

	// just the filename
	fn := info.Name()

	if info.IsDir() {
		// SKIP all directories that start with "." (execept ".")
		if len(full) > 1 && strings.HasPrefix(fn, ".") {
			return filepath.SkipDir
		}

		for _, d := range w.skipDirs {
			if d == fn {
				return filepath.SkipDir
			}
		}
	} else {
		match := false
		for _, pat := range w.pats {
			if w.matchSuffix {
				if strings.HasSuffix(fn, pat) {
					match = true
					break
				}
			} else {
				if strings.Contains(fn, pat) {
					match = true
					break
				}
			}
		}

		// print the file if it either matches the pattern or no patterns were specified (MATCH ALL)
		if match || len(w.pats) == 0 {
			fmt.Println(full)
		}
	}

	return nil
}
