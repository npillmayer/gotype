/*
Package tree implements an all-purpose tree type.

There are many tree implementations around. This one supports trees
of a fairly simple structure. However, this package makes heavy use
of concurrency for all kinds of tree operations. Tree traversal and
modification are often performed asynchronously by creating pipelines
of concurrent filters. This is done transparently for the client,
only reflected by getting a
promise (https://en.wikipedia.org/wiki/Futures_and_promises)
as a return type.

For small trees the overhead of concurrency may hurt, from a performance
point of view. This package is meant for fairly large DOMs with potentially
complex styling information. However, it is generalized enough to be useful
in other scenarios as well. And to be honest: I wrote it because
concurrency in Go is kind of fun!

We support a set of search & filter functions on tree nodes. Clients will chain
these to perform tasks on nodes (see examples).
You may think of the set of operations to form a small
Domain Specific Language (DSL). This is similar in concept to JQuery, but
of course with a much smaller set of functions.

Navigation functions:

   Parent()                     // find parent for all selected nodes
   AncestorWith(predicate)      // find ancestor with a given predicate
   DescendentsWith(predicate)   // find descendets with a given predicate

Filter functions:

   AttributeIs(key, value)      // filter for nodes with a given attribute value

More operations will follow as I get experience from using the tree in
more real life contexts.

BSD License

Copyright (c) 2017–18, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of Norbert Pillmayer nor the names of its contributors
may be used to endorse or promote products derived from this software
without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
package tree
