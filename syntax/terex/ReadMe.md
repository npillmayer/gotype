## TeREx: Term Rewriting Expressions

Parsing generates a parse tree, which is too verbose for further processing.
Instead of long chains of grammar production symbols we usualy prefer a
much more compact AST (abstract syntax tree). One possible variant of
ASTs is a *homogenous* tree, i.e. one where the structure of all nodes
is identical. This makes tree walking easy.

This module provides the core Go data types to create and modify
homogenous trees. Homogenous trees are usually built around some Node type.
However, there is a programming language specialized in homogenous lists and
trees: Lisp (or Clojure, if you prefer). We implement node types which are
reminiscent of Lisp CONS, and call the resulting mini-language TeREx
(Term Rewriting Expressions).

With homogenous tree nodes there is always one caveat: type information of the
implementing programming language is compromised. Therefore, in absence of generics,
the code in this module heavily uses "interface{}" and relies on type switches and
casts. This is sometimes cumbersome to read, but on the other hands brings convenience
for a certain set of operations, including tree walking and tree restructuring.
