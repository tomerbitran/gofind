package main

import (
	fmt "fmt"
	"io/fs"
	os "os"
	regex "regexp"
)

func main() {

	//wd, err := os.Getwd()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	wd_info, err := os.Stat(".")
	if err != nil {
		fmt.Println(err)
		return
	}

	IterateJob{
		currentDir: fs.FileInfoToDirEntry(wd_info),
		filter: NameFilter{
			NameRegex: *regex.MustCompile(".*"),
		},
		basepath: ".",
	}.run()

}
