// 19 feb 2012
// see also Rob Pike's talk "Lexical Scanning in Go"

package clrmamepro

import (
	"fmt"
	"io"
	"io/ioutil"
	"unicode"
	"unicode/utf8"
	"errors"
)

const eof = utf8.RuneError

type toktype int
const (
	tokError toktype = iota
	tokEOF
	tokText
	tokChar
)

type token struct {
	typ	toktype
	val	string
}

func (t token) String() string {
	k := map[toktype]string{
		tokError:	"error",
		tokEOF:	"eof",
		tokText:	"text",
		tokChar:	"char",
	}
	return fmt.Sprintf("%s %q", k[t.typ], t.val)
}

type lexer struct {
	input	string
	start, pos	int
	width	int			// size of current rune
	err		error
	toks		chan token
}

// emit:  send out a token and advance the input pointer
func (l *lexer) emit(typ toktype) {
	switch typ {
	case tokEOF:
		l.toks <- token{typ, ""}
	case tokError:
		l.toks <- token{typ, l.err.Error()}
	default:
		l.toks <- token{typ, l.input[l.start:l.pos]}
	}
	l.start = l.pos
}

// ignore:  drop this character
func (l *lexer) ignore() {
	l.start = l.pos
}

// next:  read next rune
func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		return eof
	}
	// apparently I can't use l.width in :=
	r, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = width
	l.pos += l.width
	return r
}

// back:  undo that
func (l *lexer) back() {
	l.pos -= l.width
}

type stateFn func(l *lexer) stateFn

func isWordChar(r rune) bool {
	return ('0' <= r && r <= '9') ||
		('A' <= r && r <= 'Z') ||
		('a' <= r && r <= 'z') ||
		r == '-'
}

func lexEOF(l *lexer) stateFn {
	l.emit(tokEOF)
	return nil
}

func lexError(l *lexer) stateFn {
	l.emit(tokError)
	return nil
}

// ( or ) or some other single character that results in an error
func lexChar(l *lexer) stateFn {
	l.emit(tokChar)
	return lexMain
}

func lexWord(l *lexer) stateFn {
	for {
		r := l.next()
		if !isWordChar(r) {
			l.back()
			break
		}
	}
	l.emit(tokText)
	return lexMain
}

func lexString(l *lexer) stateFn {
	for {
		r := l.next()
		if r == eof {
			l.err = errors.New("eof in string")
			return lexError
		}
		if r == '"' {
			l.back()		// take only the inside
			break
		}
	}
	l.emit(tokText)
	l.next()				// NOW skip "
	return lexMain
}	

func lexMain(l *lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		return lexEOF
	case unicode.IsSpace(r):
		l.ignore()
		return lexMain		// TODO should we optimize this?
	case r == '"':
		l.ignore()
		return lexString
	case isWordChar(r):
		return lexWord
	}
	return lexChar			// let the parser deal with invalid characters
}

func (l *lexer) run() {
	for state := lexMain; state != nil; {
		state = state(l)
	}
	close(l.toks)	// signal end of input
}

func lex(src io.Reader) (*lexer, error) {
	b, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}
	l := &lexer{
		input:	string(b),
		toks: 	make(chan token),
	}
	go l.run()
	return l, nil
}
