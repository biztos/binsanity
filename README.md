# binsanity

Asset files into golang source, with testing! Inspired by [go-bindata][gbd].

[![GoDoc][docbadge]][doc] [![Coverage Status][covbadge]][cov]

[docbadge]: https://pkg.go.dev/badge/github.com/biztos/binsanity.svg
[doc]: https://pkg.go.dev/github.com/biztos/binsanity
[covbadge]: https://coveralls.io/repos/github/biztos/binsanity/badge.svg
[cov]: https://coveralls.io/github/biztos/binsanity
[gbd]: https://github.com/jteeuwen/go-bindata

You usually interact directly with the command-line application:

```bash
$ go get -u github.com/biztos/binsanity/v1/...
$ binsanity --help
```

For information on the Go package itself, please refer to the official
[documentation][doc].

**Please use version v1 or higher, earlier versions may not work with `go.mod`.**

## Using binsanity

Specify your asset directory, and `binsanity` will create two files for you:
a source file and a test file. The test file provides 100% coverage of the
source file.

```bash
$ cd src/mypkg
$ binsanity my-asset-dir
$ go test -cover .
ok      github.com/you/mypkg 0.012s  coverage: 100.0% of statements
```

You can pass custom values for the output file, package name, and module
(imported for testing). By default the source file will be `binsanity.go` and
the test file will `binsanity_test.go`; the package and module are taken from
your project directory.

The generated source file defines the following functions:

- `AssetNames() []string` -- return a list of all asset names.
- `Asset(name string) ([]byte,error)` -- return data for an asset.
- `MustAsset(name string) []byte` -- as above, but panic on errors.
- `MustAssetString(name string) string` -- as above, but for strings.

Note that the design of `binsanity` is probably not conducive to very large
asset collections. Data is compressed, but also Base64-encoded; and the
lookup and caching system is fast but could potentially more than double your
memory usage. If you are very worried about efficiency, you should not use
`binsanity`. But if you are more interested in convenience and test coverage,
then you probably should. :-)

Issues are welcome if they are reproducible.

Good luck, and _design for testing!_
