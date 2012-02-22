// clrmamepro datfile processor
// 21 feb 2012
package main//clrmamepro

import (
	"errors"
	"os"

	// testing
	"fmt"
	"log"
)

type datparse struct {
	l		*lexer
	blocks	chan Block
	err		error
}

func (d *datparse) Lex(lval *yySymType) int {
	var tok token

	switch tok = <-d.l.toks; tok.typ {
	case tokEOF:
		return 0
	case tokError:
		d.Error(tok.val)
		return 0
	case tokText:
		lval.str = tok.val
		return TEXT
	}
	return int([]rune(tok.val)[0])	// TODO can this be simpler
}

func (d *datparse) Error(e string) {
log.Fatalf("ERROR! %v\n", e)
	d.err = errors.New(e)
}

type Datfile struct {
	dp		*datparse
}

func OpenDatfile(filename string) (d *Datfile, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d = new(Datfile)
	d.dp = new(datparse)
	d.dp.l, err = lex(f)
	if err != nil {
		return nil, err
	}
	d.dp.blocks = make(chan Block)
	go yyParse(d.dp)
	return d, nil
}

func (d *Datfile) GetBlock() (*Block, error) {
	b, ok := <-d.dp.blocks
	if !ok {
		if d.dp.err != nil {	// error?
			return nil, d.dp.err
		}
		return nil, nil		// else eof
	}
	return &b, nil
}


// testing
func (b Block) asString(ntab int) string {
	line := func(s string) string {
		k := ""
		for i := ntab; i != 0; i-- {
			k += "\t"
		}
		return k + s + "\n"
	}
	s := line(b.Name + " (")
	for i, x := range b.Texts {
		s += line("\t" + i + " " + fmt.Sprintf("%q", x))
	}
	for _, x := range b.Blocks {
		s += x.asString(ntab  + 1)
	}
	s += line(")")
	return s
}

func (b Block) String() string {
	return b.asString(0)
}
