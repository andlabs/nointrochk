// 20 feb 2012
%{
package clrmamepro

type Block struct {
	Name	string
	Texts	map[string]string
	Blocks	map[string]*Block
}

func makeBlock() *Block {
	b := new(Block)
	b.Texts = make(map[string]string)
	b.Blocks = make(map[string]*Block)
	return b
}
%}

%union {
	block	*Block
	str		string
}

%token <str> tokTEXT		// used for both regular words and strings
%type <block> block blockcontents

%start start
%%
start:
		block						{
			l := yylex.(*lexer)
			l.blocks = append(l.blocks, $1)
		}
	|	start block					{
			l := yylex.(*lexer)
			l.blocks = append(l.blocks, $2)
		}
	;

block:
		tokTEXT '(' blockcontents ')'		{
			$$ = $3
			$$.Name = $1
		}
	;

blockcontents:
		tokTEXT tokTEXT				{
			$$ = makeBlock()
			$$.Texts[$1] = $2
		}
	|	block						{
			$$ = makeBlock()
			$$.Blocks[$1.Name] = $1
		}
	|	blockcontents tokTEXT tokTEXT	{
			$$ = $1
			$$.Texts[$2] = $3
		}
	|	blockcontents block				{
			$$ = $1
			$$.Blocks[$2.Name] = $2
		}
	;
