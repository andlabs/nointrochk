// <codedate
package clrmamepro

import (
	"testing"
	"strings"
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

game (
	name "[BIOS] LaserActive (USA) (v1.04)"
	description "[BIOS] LaserActive (USA) (v1.04)"
	rom ( name "[BIOS] LaserActive (USA) (v1.04).md" size 131072 crc 50CD3D23 md5 0E7393CD0951D6DDE818FCD4CD819466 sha1 AA811861F8874775075BD3F53008C8AAF59B07DB )
)

game (
	name "[BIOS] Mega-CD (Japan) (1.00S)"
	description "[BIOS] Mega-CD (Japan) (1.00S)"
	rom ( name "[BIOS] Mega-CD (Japan) (1.00S).md" size 131072 crc 79F85384 md5 A3DDCC8483B0368141ADFD99D9A1E466 sha1 230EBFC49DC9E15422089474BCC9FA040F2C57EB )
)

game (
	name "[BIOS] Mega-CD (Japan) (1.00l)"
	description "[BIOS] Mega-CD (Japan) (1.00l)"
	rom ( name "[BIOS] Mega-CD (Japan) (1.00l).md" size 131072 crc F18DDE5B md5 29AD9CE848B49D0F9CEFC294137F653C sha1 0D5485E67C3F033C41D677CC9936AFD6AD618D5F )
)
`

func TestClrmamepro(t *testing.T) {
	r := strings.NewReader(example)
	d, err := NewDatfile(r)
	if err != nil {
		t.Fatalf("error creating new datfile: %v", err)
	}
	for {
		b, err := d.GetBlock()
		if b == nil && err == nil {
			break
		}
		t.Log(b, err)
	}
}
