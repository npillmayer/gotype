# Experimental GLR Parser

This package contains an experimental parser which—hopefully—we will use some day to parse Markdown. GLR parsers are rare (outside of academic research) and there is no easy-to-port version in another programming language. We will just muddle ahead and will see where we can get.

GLR parsers rely on a special stack structure, called a GSS. A GSS can hold information about alternative parser states after a conflict (shift/reduce, reduce/reduce) occured. 

For further information see for example

	https://people.eecs.berkeley.edu/~necula/Papers/elkhound_cc04.pdf
	https://cs.au.dk/~amoeller/papers/ambiguity/ambiguity.pdf

This is experimental software, currently not intended for production use in any way.