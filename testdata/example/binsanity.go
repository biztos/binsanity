/* binsanity.go - auto-generated; edit at thine own peril!

More info: https://github.com/biztos/binsanity

*/

package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"sort"
)

// Asset returns the byte content of the asset for the given name, or an error
// if no such asset is available.
func Asset(name string) ([]byte, error) {

	_, found := binsanity_cache[name]
	if !found {
		i := sort.SearchStrings(binsanity_names, name)
		if i == len(binsanity_names) || binsanity_names[i] != name {
			return nil, errors.New("Asset not found.")
		}

		// Not cached, so decode and cache it.
		decoded, err := base64.StdEncoding.DecodeString(binsanity_data[i])
		if err != nil {
			return nil, err
		}
		buf := bytes.NewReader(decoded)
		gzr, err := gzip.NewReader(buf)
		if err != nil {
			return nil, err
		}
		data, err := ioutil.ReadAll(gzr)
		if err != nil {
			return nil, err
		}
		binsanity_cache[name] = data
	}
	return binsanity_cache[name], nil

}

// MustAsset returns the byte content of the asset for the given name, or
// panics if no such asset is available.
func MustAsset(name string) []byte {
	b, err := Asset(name)
	if err != nil {
		panic(err.Error())
	}
	return b
}

// MustAssetString returns the string content of the asset for the given name,
// or panics if no such asset is available.  This is a convenience function
// for string(MustAsset(name)).
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetNames returns the sorted names of the assets.
func AssetNames() []string {
	return binsanity_names
}

// this must remain sorted or everything breaks!
var binsanity_names = []string{
	"bar",
	"baz/bat/bloopf",
	"foo",
}

// only decode once per asset.
var binsanity_cache = map[string][]byte{}

// assets are gzipped and base64 encoded
var binsanity_data = []string{
	"H4sIAEa1KWQAA0tKLFLILFZISiziAgD2TYYNCwAAAA==",
	"H4sIAGG1KWQAA0tKrFLILFZISiwBUzn5+QVpXADkmDgmFQAAAA==",
	"H4sIABxMKWQAA0vLz1fILFZIy8/nAgCgfYmECwAAAA==",
}
