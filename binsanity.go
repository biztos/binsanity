/* binsanity.go - auto-generated; edit at your own peril!

More info: https://github.com/biztos/binsanity

Generated: 0001-01-01 00:00:00 +0000 UTC

*/

package binsanity

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io"
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

		// We ignore errors because we controlled the data from the begining.
		// It's not perfect but it seems better than having additional funcs
		// hanging around that might confuse the user: tried that already, not
		// nicer.
		decoded, _ := base64.StdEncoding.DecodeString(binsanity_data[i])
		buf := bytes.NewReader(decoded)
		gzr, _ := gzip.NewReader(buf)
		defer gzr.Close()
		data, _ := io.ReadAll(gzr)

		// Not cached, so decode and cache it.
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
	"code.tmpl",
	"tests.tmpl",
}

// only decode once per asset.
var binsanity_cache = map[string][]byte{}

// assets are gzipped and base64 encoded
var binsanity_data = []string{
	"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",
}
