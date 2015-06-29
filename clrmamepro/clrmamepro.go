// clrmamepro datfile processor
// 21 feb 2012
package clrmamepro

import (
	"fmt"
	"io"
)

//go:generate go tool yacc parse.y

func Read(r io.Reader, filename string) (b []*Block, errs []string) {
	yyErrorVerbose = true
	l := newLexer(r, filename)
	yyParse(l)
	for _, e := range l.errs {
		errs = append(errs, fmt.Sprintf("%s %s", e.pos, e.msg))
	}
	if len(errs) == 0 {
		return l.blocks, nil
	}
	return nil, errs
}
