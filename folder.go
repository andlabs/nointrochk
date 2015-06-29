// 29 june 2015
package main

import (
	"os"
	"path/filepath"
	"sort"
	"sync"
)

var files map[string]struct{}
var filesLock sync.Mutex

func collectFilenames(path string) (err error) {
	filesLock.Lock()
	defer filesLock.Unlock()

	files = make(map[string]struct{})
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) {
		if err != nil {
			return err
		}
		files[path] = struct{}{}
		return nil
	})
}

// we don't even need to use package os to see if the file exists; we already know it does from the map above
func exists(name string) bool {
	filesLock.Lock()
	defer filesLock.Unlock()

	_, ok := files[name]
	return ok
}

func markProcessed(name string) {
	filesLock.Lock()
	defer filesLock.Unlock()

	delete(files, name)
}

func printLeftovers() {
	filesLock.Lock()
	defer filesLock.Unlock()

	// print the list sorted for neatness
	m := make([]string, 0, len(files))
	for f, _ := range files {
		m = append(m, f)
	}
	sort.Strings(m)
	for _, f := range m {
		alert("EXTRA", m)
	}
}
