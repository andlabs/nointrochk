// nointrochk:  check an (extracted) no-intro ROM set against a clrmamepro dat file
// 22 feb 2012

package main

import (
	"fmt"
	"os"
	"io"
	"strings"
	"hash"
	"hash/crc32"
	"crypto/md5"
	"crypto/sha1"
	"log"

//	"clrmamepro"
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
var folder, filename, basename string

func alert(method string) {
	fmt.Printf("%-10s %s\n", method, basename)
}

func passesChecksum(method hash.Hash, mname string, expected string) bool {
	// according to the crypto hash testing code in the Go source code, this is what we do
	// thanks to #go-nuts for fixing a misunderstanding of that code that I had
	// thanks to rog in #go-nuts for the Copy() trick
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("open %s to check %s failed: %v\n", filename, mname, err)
	}
	defer f.Close()
	io.Copy(method, f)
	if strings.ToUpper(fmt.Sprintf("%x", method.Sum(nil))) != strings.ToUpper(expected) {
		alert("BAD " + strings.ToUpper(mname))
		nbad++
		return false
	}
	return true
}

func fileExistsAndSameSize(size string) bool {
	f, err := os.Stat(filename)
	if err != nil {
		alert("MISSING")
		nmiss++
		return false
	}
	if fmt.Sprintf("%d", f.Size()) != size {
		alert("BAD SIZE")
		nbad++
		return false
	}
	return true
}

func do(b Block) {
	var good bool

	stats := b.Blocks[fROM]
	basename = stats.Texts[fFilename]
	filename = folder + string(os.PathSeparator) + basename
	good = fileExistsAndSameSize(stats.Texts[fSize])
	good = good && passesChecksum(crc32.NewIEEE(), "crc32", stats.Texts[fCRC32])
	good = good && passesChecksum(md5.New(), "md5", stats.Texts[fMD5])
	good = good && passesChecksum(sha1.New(), "sha1", stats.Texts[fSHA1])
	if good {
		alert("GOOD")
		ngood++
	}
}

func main() {
	d, err := OpenDatfile(os.Args[1])
	if err != nil {
		log.Fatalf("can't open datfile: %v\n",err)	
	}
	folder = os.Args[2]
	for {
		b, err := d.GetBlock()
		if b == nil && err == nil {
			break
		}
		if err == nil {
			if b.Name == fGame {
				nroms++
				do(*b)
			}
		} else {
			log.Fatalf("error reading block: %v\n", err)
		}
	}
	fmt.Printf("%d ROMs, %d good, %d bad, %d missing (%f%% good)\n",
		nroms, ngood, nbad, nmiss,
		(float64(ngood)/float64(nroms))*100.)
}