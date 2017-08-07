# Walk
A simplified `find` command that recursively walks a directory and finds files that match the given patterns.  All files that begin with "." are ignored.  Walk is written in **Go**, is (currently) under 100 lines and code and should easily compile on all platforms that support Go.

## Install

1. checkout the code
2. cd walk
3. go install

If go/bin is in your PATH, simply type `walk`.

## Command line help:

```
The walk utility recursively descends the directory tree for the specified directory
and prints files that match the specified substring patterns.

Usage:
  walk <directory> [pat1 pat2 ...]

Flags:
  -x	treat patterns as file extensions
```

## Examples

Find all web related files in the current directory:

    walk -x . js css html

Find all files that contain 2017 in the filename:

    walk . 2017

Find all files in another directory

    walk ~/Sites/img

