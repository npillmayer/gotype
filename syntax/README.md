
Purposes of Macro Languages
----------------------

TeX features a powerful macro language. For newcomers to TeX-typesetting its
language may look cumbersome -- or even ugly or funny -- but it is well thought
out and surprisingly versatile. It serves a couple of purposes, which we will
discuss.

  - Configurability
  - Adaptability
  - Extensibility

### Configurability

When formatting your text you want to be able to configure a lot of things:
fonts, colors, running headers, etc. Necessary configuration is often done
in templates, but plain TeX does not support a templating mechnism. Instead,
configuration is provided through a macro file.

### Adaptability

TeX is used for all kinds of "documents." Obviously, it has been created for
scientific papers and books, possibly heavy on mathematics, but over the year
TeX has found usages besides that as well: music notes, chemistry, school books,
and others. This testifies the versatility gained by the TeX macro language.

### Extensibility

The best software systems usually provide a means of extending their
functionality. TeX does this by providing hooks which may be extended or
overridden by macros. Examples are token registers like `everypar` or
the page output routine. There is even a generic extension mechanism called
*specials*, which can be used for all kinds of additional functionality.


## Lua

Modern versions of TeX support Lua as an extension language. Embedding Lua in an application is all the rage today, and we will follow this tradition. However, we will employ at least one Domain Specific Language: a derivate of MetaPost, which will be called (working title) "Poor Man's MetaPost" (PMMPos). MetaPost includes a powerful macro-language, which we will not support. Instead, we will make PMMPost extensible via Lua.

## Parser Generators

We use ANTLR V4 as our goto-tool for implementing grammars, e.g., for PMMPost.
Additionally we will create an experimental GLR parser. It is intended for parsing Markdown, but the idea still needs some working-on.
