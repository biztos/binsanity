// binspkg.go -- just the binsanity package description.

// Package binsanity encodes files into Go source, with testing.
//
// Inspired by the bindata package, binsanity aims to provide a useful subset
// of features while also enabiling thorough testing of the generated Go
// source code.
//
// For a much more featureful, but less testable approach see:
//
// https://pkg.go.dev/github.com/jteeuwen/go-bindata
//
// One generally does not use this package directly, but rather the command
// binsanity.
//
// This generates two Go source files: one with following functions, and one
// testing those functions:
//
// # Asset - return an asset's data as a []byte
//
// # MustAsset - retrieve an asset's bytes or panic if not found
//
// # MustAssetString - call MustAsset and return its result as a string
//
// # AssetNames - return a list of asset names as a []string
//
// Assets are gzipped and base64-encoded; they are decoded and inflated only
// once, with the result cached.
//
// The resulting source files introduce no dependencies outside the Go
// standard library.
package binsanity
