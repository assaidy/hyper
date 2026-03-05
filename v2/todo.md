# hyper v2 TODO

- [ ] Use all-caps for elements. Rename all element functions to uppercase:
    - `Div()` → `DIV()`
    - `P()` → `P()`
    - `Table()` → `TABLE()`
    - `Thead()` → `THEAD()`
    - `DoctypeHtml()` → `DOCTYPE()`
    - And all other elements similarly

- [ ] Remove htmx directory

- [ ] Remove events.go
    - Delete `events.go` entirely
    - All `Event*` constants were only useful for htmx integration (hx.AttrHxOn())
    - Library should be framework-agnostic (supports HTMX, alpinejs, Datastar, or any other framework)

- [ ] Updated test suite and benchmarks
    - Update benchmarks to use `DIV()` syntax
    - Update all test files to use new element names
    - Keep only pure HTML builder functionality

## Goal
Pure, framework-agnostic HTML builder with:
- Fast performance
- Easy to use
- Compile-time type safety
- Helper constants for normal HTML attributes only (AttrClass, AttrId, etc.)
- No dependencies on any frontend framework (HTMX, Datastar, etc.)
