// clrmamepro datfile processor: lexical analyzer
// 20 feb 2012
%{
package main//clrmamepro

type Block struct {
	Name	string
	Texts	map[string]string
	Blocks	map[string]Block
	Error		error
}

type yySymType Block
%}

%token TEXT		// used for both regular words and strings

%start main
%%
main:
		block
	|	main block
	;

block:
		TEXT '(' blockcontents ')'		{
			$$ = $3
			$$.Name = $1
		}
	;

blockcontents:
		TEXT TEXT				{
			$$.Texts[$1] = $2
		}
	|	block					{
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
