// 29 june 2015
package main

import (
	"fmt"
	"os"
	"io"
	"path/filepath"
	"hash"
	"hash/crc32"
	"crypto/md5"
	"crypto/sha1"
	"strings"
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
	crc := crc32.NewIEEE()
	md := md5.New()
	sha := sha1.New()
	checksums := io.MultiWriter(crc, md, sha)

	_, err := io.Copy(checksums, f)
	if err != nil {
		die("checksumming ROM %q: %v", romname, err)
	}

	compare := func(h hash.Hash, expected string) bool {
		return strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil))) == strings.ToUpper(expected)
	}
	if !compare(crc, rom.Texts[fCRC32]) {
		alert("BAD CRC32", romname)
		return false
	}
	if !compare(md, rom.Texts[fMD5]) {
		alert("BAD MD5", romname)
		return false
	}
	if !compare(sha, rom.Texts[fSHA1]) {
		alert("BAD SHA1", romname)
		return false
	}

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
