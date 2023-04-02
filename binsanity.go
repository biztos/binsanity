// binsanity.go -- like bindoc but stupider and sane.

// Package binsanity encodes files into Go source, with testing.
//
// Inspired by the bindata package, binsanity aims to provide a minimally
// useful subset of features while also enabiling thorough testing of the
// generated Go source code.
//
// For a much more featureful, but less testable approach see:
//
// https://pkg.go.dev/github.com/jteeuwen/go-bindata
//
// One generally does not use this package directly, but rather the command
// binsanity.
//
// # Differences From Bindata
//
// * Data is not compressed.
//
// * Only the AssetNames, Asset, and MustAsset functions are implemented.
//
// * Edge cases, probably numerous, have not been much considered.
//
// * Not recommended for large assets or large collections of assets.
//
// * Your test coverage will not be reduced. :-)
package binsanity
