// clrmamepro datfile processor
// 21 feb 2012
package clrmamepro

import (
	"io"
	"errors"
)

//go:generate go tool yacc parse.y

type datparse struct {
	l		*lexer
	blocks	chan *Block
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

func NewDatfile(r io.Reader) (d *Datfile, err error) {
	d = new(Datfile)
	d.dp = new(datparse)
	d.dp.l, err = lex(r)
	if err != nil {
		return nil, err
	}
	d.dp.blocks = make(chan *Block)
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
	return b, nil
}
