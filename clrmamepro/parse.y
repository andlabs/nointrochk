// 20 feb 2012
%{
package clrmamepro

type Block struct {
	Name	string
	Texts	map[string]string
	Blocks	map[string]Block
	Error		error
}

func mk(b *Block) {
	if b.Texts == nil {
		b.Texts = make(map[string]string)
	}
	if b.Blocks == nil {
		b.Blocks = make(map[string]Block)
	}
}
%}

%union {
	block	Block
	str		string
}

%token <str> TEXT		// used for both regular words and strings
%type <block> block blockcontents

%start start
%%
start:
		blocks					{
			close(yylex.(*datparse).blocks)
		}
	;

blocks:
		block					{
			yylex.(*datparse).blocks <- $1
		}
	|	blocks block				{
			yylex.(*datparse).blocks <- $2
		}
	;

block:
		TEXT '(' blockcontents ')'		{
			$$ = $3
			$$.Name = $1
		}
	;

blockcontents:
		TEXT TEXT				{
			mk(&$$)
			$$.Texts[$1] = $2
		}
	|	block					{
			mk(&$$)
			$$.Blocks[$1.Name] = $1
		}
	|	blockcontents TEXT TEXT		{
			$$ = $1
			$$.Texts[$2] = $3
		}
	|	blockcontents block			{
			$$ = $1
			$$.Blocks[$2.Name] = $2
		}
	;
