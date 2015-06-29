// 15 april 2015
package clrmamepro

import (
	"io"
	"text/scanner"
	"strconv"
	"unicode"
)

type lexerr struct {
	msg		string
	pos		scanner.Position
}

type lexer struct {
	scanner	*scanner.Scanner
	blocks	[]*Block
	errs		[]lexerr
}

func newLexer(r io.Reader, filename string) *lexer {
	l := new(lexer)
	l.scanner = new(scanner.Scanner)
	l.scanner.Init(r)
	l.scanner.Error = func(s *scanner.Scanner, msg string) {
		l.Error(msg)
	}
	l.scanner.Mode = scanner.ScanIdents | scanner.ScanStrings
	l.scanner.Position.Filename = filename
	l.scanner.IsIdentRune = func(ch rune, i int) bool {
		if ch == '-' {
			return true
		}
		// unfortunately there's no exported function for this in Go (TODO); we had to copy this from the text/scanner source (TODO licensing)
		// the only difference is that we allow a digit to start an identifier
		return ch == '_' || unicode.IsLetter(ch) || unicode.IsDigit(ch)
	}
	return l
}

func (l *lexer) Lex(lval *yySymType) int {
	r := l.scanner.Scan()
	switch r {
	case scanner.EOF:
		return 0
	case scanner.Ident:
		lval.str = l.scanner.TokenText()
		return tokTEXT
	case scanner.String:
		ss := l.scanner.TokenText()
		// the token text is a Go string in a string!
		ss, err := strconv.Unquote(ss)
		if err != nil {
			l.Error(err.Error())
		}
		lval.str = ss
		return tokTEXT
	}
	return int(r)
}

func (l *lexer) Error(s string) {
	l.errs = append(l.errs, lexerr{
		msg:		s,
		pos:		l.scanner.Pos(),
	})
}
