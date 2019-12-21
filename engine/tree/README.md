Concurrent Tree Operations
==========================

There are many tree implementations around. This one supports trees
of a fairly simple structure. However, this package makes heavy use
of concurrency for all kinds of tree operations. Tree traversal and
modification are often performed asynchronously by creating pipelines
of concurrent filters. This is done transparently for the client,
only reflected by getting a
[promise](https://en.wikipedia.org/wiki/Futures_and_promises)
as a return type.

For small trees the overhead of concurrency may hurt, from a performance
point of view. This package is meant for fairly large DOMs with potentially
complex styling information. However, it is generalized enough to be useful
in other scenarios as well. And to be honest: I wrote it because
concurrency in Go is kind of fun!

DSL
---

We support a set of search & filter functions on tree nodes. Clients will chain
these to perform tasks on nodes (see examples below).
You may think of the set of operations to form a small
Domain Specific Language (DSL). This is similar in concept to JQuery, but
of course with a much smaller set of functions.

Navigation functions:

    TopDown(action)              // operate on all nodes, starting at the top
    Parent()                     // find parent for all selected nodes
    AncestorWith(predicate)      // find ancestor with a given predicate
    DescendentsWith(predicate)   // find descendets with a given predicate
    AllDescendents()             // find all descendents of selected nodes

Filter functions:

    AttributeIs(key, value)      // filter for nodes with a given attribute value
    SetAttribute(key, value)     // set an attribute value for nodes
    Filter(userfunc)             // apply a user-provided filter function

More operations will follow as I get experience from using the tree in
more real life contexts.

Tree Walker
---

Operations on trees start with creating a *Walker* for the tree.
A typical usage of a Walker looks like this ("*FindNodesAndDoSomething*()" is
a placeholder for a sequence of DSL function calls):

    w := NewWalker(rootnode)
    futureResult := w.FindNodesAndDoSomething(...).Promise()
    nodes, err := futureResult()

Clients always have to call *Promise*() as the final link of the
DSL expression chain, even if they do not expect the expression to
return a non-empty set of nodes.