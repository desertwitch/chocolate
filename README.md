# Chocolate

Layout container for Bubble Tea the easy way ;) .... *hopefully*
Chocolate is meant to create complex and responsive layouts with
themes via flavour working consistent and without hassle together
in TUI applications.

## State / Idea

The first implementation wasn't really fitting my intentions. Some
of the problems that I identified are:

- Restrictive interface (forced to work in a own specific way for chocolate)
- Complicated calculations and inefficient handling of the bubbletea layer
- Force to use chocolate as management layer for the whole bubbletea TUI

So I took the time to work out a way better implementation that solves the main
issues and has way better layout handling.

The layout handling is done by using the cassowary algorithm and is implemented
like the [GTK constraint layouts](https://docs.gtk.org/gtk4/class.ConstraintLayout.html).

Further the idea is now to provide a library that is focused only on styling and
layouts and not on handling the whole bubbletea model handling.

This makes it now able to work out a chocolate specific components set analog
[bubbles](https://github.com/charmbracelet/bubbles), which will use the chocolate as layout workhorse. It is focusing on
UI components in the way, that it shall provide mainly functionality like
dialog (ae.: yes, no) via overlays, status bars, menus, unify the selection of
models that provide input functionality, handling the visibility and input
focus for the underlying models and so on.

This time I'm very convinced that the API is providing what I wanted to
have in the first place.

~~This is a very early state and in theory the whole API can change
at any time.~~
~~Still now is the best time to get *ANY* feedback to change the way~~
~~in the right direction if needed.~~

To get started, please see the [examples](https://github.com/mfulz/chocolate/tree/master/examples)

## Priority

The following tasks are in order regarding my personal priority scoring:

- Providing some basic examples
- Create the bars repository with some useful bars
- Documentation
