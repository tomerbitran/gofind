package main

import (
	"fmt"
	"io/fs"
	"os"
	filepath "path/filepath"
	"sync"
)

type Job interface {
	run()
}

type IterateJob struct {
	currentDir fs.DirEntry
	filter     *Filter
	basepath   string
}

type EntryJob struct {
	wg       *sync.WaitGroup
	file     fs.DirEntry
	basepath string
	filter   *Filter
}

func (eg EntryJob) run() {
	defer (*eg.wg).Done()

	if eg.file.IsDir() {
		// add iterate job to jobs
		base_path := filepath.Join(eg.basepath, eg.file.Name())

		IterateJob{

			currentDir: eg.file,
			filter:     eg.filter,
			basepath:   base_path,
		}.run()

	}
	// Filter the message here
	if (*eg.filter).Filter(eg.file, eg.basepath) {
		fmt.Println(filepath.Join(eg.basepath, eg.file.Name()))
	}
	return
}

func (ij IterateJob) run() {
	// Create a new job for each file
	files, err := getFilesInDirectory(ij.basepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gofind:", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(files))

	for _, file := range files {
		entryJob := EntryJob{
			wg:       &wg,
			file:     file,
			basepath: ij.basepath,
			filter:   ij.filter,
		}
		go entryJob.run()
	}

	wg.Wait()

	return
}

func getFilesInDirectory(dir string) ([]fs.DirEntry, error) {
	direntries, error := os.ReadDir(dir)
	if error != nil {
		return nil, error
	}
	return direntries, nil
}
