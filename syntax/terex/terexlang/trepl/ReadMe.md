## Experiments with Term Rewriters

Whenever we work on creating a parser for an new language, there is a point where we want to experiment with test input.
We may do this with unit tests (and we should!), but often enough it's easier to have an interactive parser shell running,
where we can hack in some input strings and immediately receive feedback from the parser.

**T.REPL** (term rewriting REPL) aims to be just that: a test bed for interactive experiments with parsers and term rewriters,
aiding developers to refine AST creation and rewriting.

### Status

T.REPL is work in progress. I am currently using it to develop a grammar for BiDi text setting, extending T.REPL as I go.
It is not yet production ready. As soon as it is more polished, I will extract it as a separate module.
*Please be patient*!

### Example

There is a small demo grammar included, which is used by default. It's a simple grammar for arithmetic expressions
(with limitations).

*More to come…*
