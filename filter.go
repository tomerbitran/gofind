package main

import (
	"io/fs"
	regex "regexp"
)

type Filter interface {
	Filter(file fs.DirEntry, basepath string) bool
}

type PathFilter struct {
	PathRegex regex.Regexp
}

func (pf PathFilter) Filter(file fs.DirEntry, basepath string) bool {
	return pf.PathRegex.MatchString(basepath)
}

type NameFilter struct {
	NameRegex regex.Regexp
}

func (nf NameFilter) Filter(file fs.DirEntry, basepath string) bool {
	return nf.NameRegex.MatchString(file.Name())
}
