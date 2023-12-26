package main

import (
	flag "flag"
	fmt "fmt"
	"io/fs"
	os "os"
	regex "regexp"
)

func createFilter(hasPathFilter bool, nameRegexpPattern string, pathRegexpPattern string) Filter {
	if hasPathFilter {
		regexp, err := regex.Compile(pathRegexpPattern)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Invalid regex")
			return nil
		}
		return PathFilter{
			PathRegex: regexp,
		}
	} else { //hasNameFilter
		regexp, err := regex.Compile(nameRegexpPattern)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Invalid regex")
			return nil
		}
		return NameFilter{
			NameRegex: regexp,
		}
	}
}

func main() {
	var basePath string
	flag.StringVar(&basePath, "base-path", ".", "base path to start search")

	var pathRegexPattern string
	flag.StringVar(&pathRegexPattern, "path", "", "filter by path (exclusive to name filter)")

	var nameRegexPattern string
	flag.StringVar(&nameRegexPattern, "name", "", "filter by name (exclusive to path filter)")

	flag.Parse()

	hasPathFilter := pathRegexPattern != ""
	hasNameFilter := nameRegexPattern != ""

	if hasPathFilter == hasNameFilter { // not xor
		fmt.Fprintln(os.Stderr, "gofind: Must specify exactly one filter type")
		return
	}

	filter := createFilter(hasPathFilter, nameRegexPattern, pathRegexPattern)

	dir_info, err := os.Stat(basePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gofind: bad basepath,", err)
		return
	}

	IterateJob{
		currentDir: fs.FileInfoToDirEntry(dir_info),
		filter:     &filter,
		basepath:   basePath,
	}.run()

}
