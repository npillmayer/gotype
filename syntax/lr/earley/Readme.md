Earley-Parsing
==============

A great introduction to Earley-parsing may be found at
[Loup  Vaillant's Blog](http://loup-vaillant.fr/tutorials/earley-parsing).
Here is what he has to say:

> The biggest advantage of Earley Parsing is its accessibility. Most other tools such as parser generators,
> parsing expression grammars, or combinator libraries feature restrictions that often make them hard
> to use. Use the wrong kind of grammar, and your PEG will enter an infinite loop. Use another wrong
> kind of grammar, and most parser generators will fail. To a beginner, these restrictions feel most
> arbitrary: it looks like it should work, but it doesn't. There are workarounds of course, but
> they make these tools more complex.
>
> Earley parsing Just Works™.
>
> On the flip side, to get this generality we must sacrifice some speed. Earley parsing cannot
> compete with speed demons such as Flex/Bison in terms of raw speed. 

If speed (or the lack thereof) is critical to your project, you should probably grab ANTLR or
Bison. I used both a lot in my programming life. However, there are many scenarios where I wished
I had a more lightweight alternative at hand. Oftentimes I found myself writing recursive-descent
parsers for small ad-hoc languages by hand, sometimes mixing them with the lexer-part of one of
the big players.

