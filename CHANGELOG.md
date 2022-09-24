# Changelog

## v1.0.3

**Summary**: Modernization of this package to Go standards in 2022, mostly
through internal cleanups.

**New Features**: None

The last time this package was significantly modified was 2016, which predates
`cmp`, subtests, the modern Protocol Buffer implementation, and numerous Go
practices that emerged in the intervening years.  The new release is tested
against Go 1.19, though I expect it would work with Go 1.13 just fine.

Finally, I declared bankruptcy on the vendored test fixtures and opted for
creating my own.  This is due to the underlying implementation of the generated
code in conjunction with working with a moving target that is an external data
model representation.

## v1.0.2

**Summary**: Tagged version with Go module support.

**New Features**: None

End-users wanted a tagged release that includes Go module support.