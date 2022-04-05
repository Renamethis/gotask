package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// Format directory size to necessary units of measurement
func formatSize(size int64) string {
	if size > 1024*1024*1024 {
		return strconv.FormatInt(size/(1024*1024*1024), 10) + " GB"
	} else if size > 1024*1024 {
		return strconv.FormatInt(size/(1024*1024), 10) + " MB"
	} else if size > 1024 {
		return strconv.FormatInt(size/(1024), 10) + " KB"
	} else {
		return strconv.FormatInt(size, 10) + " B"
	}
}

// Print directory and subdirectories size
func dirSize(limiter chan bool, path string, wg *sync.WaitGroup) {
	// Wait other gorountines if number of processes exceeded
	limiter <- true
	defer func() {
		<-limiter
		wg.Done()
	}()
	var size int64
	// Get absolute path of given path
	absolute, _ := filepath.Abs(path)
	// Enumerate all files and subdirectories in directory
	err := filepath.Walk(absolute, func(subpath string, info os.FileInfo,
		err error) error {
		if err != nil {
			return err
		}
		// Recursive call goroutine for subdirectory
		absoluteSub, err := filepath.Abs(subpath)
		if info.IsDir() && absolute != absoluteSub {
			wg.Add(1)
			go dirSize(limiter, subpath, wg)
		} else {
			size += info.Size()
		}
		return err
	})
	// Print size or error
	if err == nil {
		fmt.Println(absolute, " - ", formatSize(size))
	} else {
		fmt.Println("Error " + err.Error())
	}
}
func main() {
	start := time.Now()
	var wg sync.WaitGroup
	maxProcesses, _ := strconv.Atoi(os.Args[1])
	limiter := make(chan bool, maxProcesses)
	// Iterate all command line arguments
	for _, path := range os.Args[2:] {
		wg.Add(1)
		// Run goroutine for calculating size of directory
		go dirSize(limiter, path, &wg)
	}
	// Wait goroutines
	wg.Wait()
	fmt.Printf("Time took: %v\n", time.Since(start))
}
