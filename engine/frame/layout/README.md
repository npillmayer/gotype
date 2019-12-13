LibNSLayout Architecture
========================

LibNSLayout is a library for performing layout on a Document Object Model
for HTML.  Its purpose is to allow client applications to provide DOM
information and convert that into a render list, which can be displayed
by the client.

Dependencies
------------

Clients of LibNSLayout must use the following additional libraries, because
their types are used in the LibNSLayout interface:

*   `LibDOM` is used to provide the DOM interface.
*   `LibCSS` is used to provide the CSS handling.
*   `LibWapcaplet` is used for interned strings.
*   `LibNSLog` for logging.

Interface
---------

The devision of responsibilities between LibNSLayout and its clients are
as follows:

### Client

*   Fetching the document to be displayed.
*   Creating a CSS selection context (with default user-agent, and user CSS).
*   Generating DOM.
*   Creating a LibNSLayout layout for the document, passing the DOM document,
    CSS selection context, appropriate CSS media descriptor, and scale.
*   Listening to DOM changes.
    *   Fetching resources needed by DOM.
        *   CSS (STYLE elements, and LINK elements):
            *   Parsing the CSS.
            *   Updating CSS selection context as stylesheets are fetched,
                and notifying LibNSLayout.
        *   JavaScript (SCRIPT elements, and LINK elements)
            *   Executing JavaScript.
        *   Favicons (LINK elements.)
        *   Images, Frames, Iframes.
    *   Notifying LibNSLayout of DOM changes.
*   Performing resource fetches on behalf of LibNSLayout.
    *   (Such as when LibNSLayout requires a background image or web font for
        an element due to CSS.)
*   Asking LibNSLayout to perform layout.
    *   Displaying the returned render list.
*   Asking LibNSLayout for layout info (e.g. due to JavaScript.)
*   Passing mouse actions to LibNSLayout.
*   Passing keyboard input to LibNSLayout.
*   Passing scale changes to LibNSLayout.
*   Performing measurement of text; given a string & style, calculating its
    width in pixels.

### LibNSLayout

*   Creates a layout object that's opaque to the client, and returns its
    handle.
*   Performs CSS selection as appropriate when DOM changes.
*   Asking client to fetch a resource that's needed for a computed style.
*   Asking client to measure text.
*   Performs line breaking.
*   Performs layout (if required) when asked by client and returns render list.
*   Performs layout (if required) when asked by client for layout info.
# Box Layout

Module *Layout* is responsible for creating the render trecreating the render tree.
This includes line-breaking for paragraphs.

When the client requires a layout, we would then walk the tree, ensuring every
layout node has a valid `x`, `y`, `width`, and `height`.

## Rendering

The client is responsible for rendering, and any compositing.  When the
client asks the layouter to layout, we return a render list.
