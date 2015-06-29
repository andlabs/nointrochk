// 29 june 2015
package clrmamepro

import (
	"testing"
	"strings"
	"reflect"
)

var example = `clrmamepro (
	name "Sega - Mega Drive - Genesis"
	description "Sega - Mega Drive - Genesis"
	version 20150628-135736
	comment "no-intro | www.no-intro.org"
)

game (
	name "[BIOS] CDX (USA) (v2.21X)"
	description "[BIOS] CDX (USA) (v2.21X)"
	rom ( name "[BIOS] CDX (USA) (v2.21X).md" size 131072 crc D48C44B5 md5 BACA1DF271D7C11FE50087C0358F4EB5 sha1 2B125C0545AFA089B617F2558E686EA723BDC06E )
)

game (
	name "[BIOS] LaserActive (USA) (v1.02)"
	description "[BIOS] LaserActive (USA) (v1.02)"
	rom ( name "[BIOS] LaserActive (USA) (v1.02).md" size 131072 crc 3B10CF41 md5 691C3FD368211280D268645C0EFD2EFF sha1 8AF162223BB12FC19B414F126022910372790103 )
)

game (
	name "[BIOS] LaserActive (Japan) (v1.02)"
	description "[BIOS] LaserActive (Japan) (v1.02)"
	rom ( name "[BIOS] LaserActive (Japan) (v1.02).md" size 131072 crc 00EEDB3A md5 A5A2F9AAE57D464BC66B80EE79C3DA6E sha1 26237B333DB4A4C6770297FA5E655EA95840D5D9 )
)
`

var expected = []*Block{
	&Block{
		Name:	"clrmamepro",
		Texts:	map[string]string{
			"name":		"Sega - Mega Drive - Genesis",
			"description":	"Sega - Mega Drive - Genesis",
			"version":		"20150628-135736",
			"comment":	"no-intro | www.no-intro.org",
		},
		Blocks:	map[string]*Block{},		// needed because of how reflect.DeepEqual() works
	},
	&Block{
		Name:	"game",
		Texts:	map[string]string{
			"name":		"[BIOS] CDX (USA) (v2.21X)",
			"description":	"[BIOS] CDX (USA) (v2.21X)",
		},
		Blocks:	map[string]*Block{
			"rom":	&Block{
				Name:	"rom",
				Texts:	map[string]string{
					"name":	"[BIOS] CDX (USA) (v2.21X).md",
					"size":	"131072",
					"crc":	"D48C44B5",
					"md5":	"BACA1DF271D7C11FE50087C0358F4EB5",
					"sha1":	"2B125C0545AFA089B617F2558E686EA723BDC06E",
				},
				Blocks:	map[string]*Block{},
			},
		},
	},
	&Block{
		Name:	"game",
		Texts:	map[string]string{
			"name":		"[BIOS] LaserActive (USA) (v1.02)",
			"description":	"[BIOS] LaserActive (USA) (v1.02)",
		},
		Blocks:	map[string]*Block{
			"rom":	&Block{
				Name:	"rom",
				Texts:	map[string]string{
					"name":	"[BIOS] LaserActive (USA) (v1.02).md",
					"size":	"131072",
					"crc":	"3B10CF41",
					"md5":	"691C3FD368211280D268645C0EFD2EFF",
					"sha1":	"8AF162223BB12FC19B414F126022910372790103",
				},
				Blocks:	map[string]*Block{},
			},
		},
	},
	&Block{
		Name:	"game",
		Texts:	map[string]string{
			"name":		"[BIOS] LaserActive (Japan) (v1.02)",
			"description":	"[BIOS] LaserActive (Japan) (v1.02)",
		},
		Blocks:	map[string]*Block{
			"rom":	&Block{
				Name:	"rom",
				Texts:	map[string]string{
					"name":	"[BIOS] LaserActive (Japan) (v1.02).md",
					"size":	"131072",
					"crc":	"00EEDB3A",
					"md5":	"A5A2F9AAE57D464BC66B80EE79C3DA6E",
					"sha1":	"26237B333DB4A4C6770297FA5E655EA95840D5D9",
				},
				Blocks:	map[string]*Block{},
			},
		},
	},
}

func TestClrmamepro(t *testing.T) {
	r := strings.NewReader(example)
	d, err := NewDatfile(r)
	if err != nil {
		t.Fatalf("error creating new datfile: %v", err)
	}
	for i, e := range expected {
		// TODO change to Next()
		b, err := d.GetBlock()
		if err != nil {
			t.Fatalf("error reading block: %v", err)
		}
		if b == nil {
			t.Fatalf("unexpected end of stream reading block %d", i)
		}
		if !reflect.DeepEqual(b, e) {
			t.Fatalf("block %d differs:\nexpected %#v\ngot %#v", i, e, b)
		}
	}
	b, err := d.GetBlock()
	if err != nil {
		t.Fatalf("error finishing reading blocks: %v", err)
	}
	if b != nil {
		t.Fatalf("another block was returned when no more were expected")
	}
}
