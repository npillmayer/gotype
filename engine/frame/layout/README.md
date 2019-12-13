# Box Layout

Module *Layout* is responsible for creating the render trecreating the render tree.
This includes line-breaking for paragraphs.

When the client requires a layout, we would then walk the tree, ensuring every
layout node has a valid `x`, `y`, `width`, and `height`.

## Rendering

The client is responsible for rendering, and any compositing.  When the
client asks the layouter to layout, we return a render list.
