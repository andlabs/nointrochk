// 29 june 2015
package main

import (
	"fmt"
	"os"
	"io"
	"path/filepath"
	"github.com/andlabs/nointrochk/clrmamepro"
)

// fields in dat file
const (
	fGame = "game"
	fROM = "rom"
	fFilename = "name"
	fSize = "size"
	fCRC32 = "crc"
	fMD5 = "md5"
	fSHA1 = "sha1"
)

var nroms, ngood, nbad, nmiss uint

func checksumsPass(romname string, rom *clrmamepro.Block, f io.Reader) bool {
	// TODO
	return true
}

func check(b *clrmamepro.Block, folder string) {
	if b.Name != fGame {
		return
	}
	nroms++
	rom := b.Blocks[fROM]
	if rom == nil {
		die("checking ROM: game block without rom block: %#v", b)
	}
	romname := rom.Texts[fFilename]
	if !exists(romname) {
		alert("MISSING", romname)
		nmiss++
		return
	}
	defer markProcessed(romname)

	filename := filepath.Join(folder, romname)
	f, err := os.Open(filename)
	if err != nil {
		die("opening ROM %q to check it: %v", filename, err)
	}
	defer f.Close()

	// first check size
	s, err := f.Stat()
	if err != nil {
		die("getting ROM %q stats to check size: %v", romname, err)
	}
	if fmt.Sprint(s.Size()) != rom.Texts[fSize] {
		alert("BAD SIZE", romname)
		nbad++
		return
	}

	// now check checksums
	if !checksumsPass(romname, rom, f) {
		nbad++
		return
	}

	// otherwise it all checked out
	alert("GOOD", romname)
	ngood++
}
