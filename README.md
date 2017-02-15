# binsanity

Assets into golang source, with testing!

[![GoDoc][b1]][doc] [![Build Status][b2]][ci] [![Coverage Status][b3]][cov]


[b1]: https://godoc.org/github.com/biztos/binsanity?status.svg
[doc]: https://godoc.org/github.com/biztos/binsanity
[b2]: https://travis-ci.org/biztos/binsanity.svg?branch=master
[ci]: https://travis-ci.org/biztos/binsanity
[b3]: https://coveralls.io/repos/github/biztos/binsanity/badge.svg
[cov]: https://coveralls.io/github/biztos/binsanity

    $ go get -u github.com/biztos/binsanity/...
    $ cd $GOPATH/src/my-project
    $ cat my-assets/foo.txt
    I am foo!
    $ cat mypkg.go
    package mypkg

    func Foo() string {
        return string(MustAsset("foo.txt"))
    }
    $ cat mypkg_test.go
    package mypkg_test

    import (
        "testing"
        "my-project"
    )

    func TestFoo(t *testing.T) {
        if mypkg.Foo() == "" {
            t.Fatal("no foo")
        }
    }
    $ binsanity my-assets
    $ go test -cover
    PASS
    coverage: 100.0% of statements
    ok  	my-project	0.002s
    
This is a work in progress. Inspired by [go-bindata][gbd], which is awesome
but not very testing-friendly.

## Warning - Experimental Software!

I wrote this because I needed a subset of [go-bindata][gbd] functionality and
also needed full test coverage.

It is not thouroughly tested itself, and probably has bugs.  Use at your own
risk, obviously.  Bug reports are welcome via the Github project page.

Also note that **asset data is not compressed.**

## Using binsanity

Normally you will only interact with the `binsanity` command, which will be
built in your `$GOPATH/bin` directory after you get it using:

    go get -u github.com/biztos/binsanity/...

For most cases, you simply run it in your project source directory, with a
single argument: the directory holding your asset files.

These will be encoded into a file `binsanity.go` which will have full
coverage through the `binsanity_test.go` file. The package name will be
whatever is used in your other project files, or `main` if none is found. (If
the package is `main` then no test file will be generated.)

The behavior can be made explicit via command-line options:

* `--package=PKG` -- use PKG instead of guessing the package name.
* `--import=PATH` -- use PATH instead of guessing the packag import path.
* `--output=FILE` -- write package to FILE instead of `binsanity.go`.

The package file will define the following functions:

* `AssetNames() []string` -- return a list of all asset names.
* `Asset(name string) ([]byte,error)` -- return data for an asset.
* `MustAsset(name string) []byte` -- as above, but panic on errors.

These behave as in [go-bindata][gbd].  For most common use-cases you will
want to use `MustAsset`, since you already know what assets you have:

    if opts.Help {
        fmt.Printf("%s", MustAsset("help.md"))
    }

[gbd]: https://github.com/jteeuwen/go-bindata

Good luck, and *design for testing!*
