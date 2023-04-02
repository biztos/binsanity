// OK, as long as we're here... need to have a way of skipping the content
// checks UNLESS we have some env var set. Because choosing to test content
// hashes at build is stupid.

package main_test

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"

	main "biztos.com/example"
)

const BinsanityAssetMissing = "foo:NOPE" // derived from below
const BinsanityAssetPresent = "foo"      // last one

var BinsanityAssetNames = []string{
	"bar",
	"baz/bat/bloopf",
	"foo",
}
var BinsanityAssetSum = map[string]string{
	"bar":            "947f16f7a547deb1e5d8aa2896cfa00a4903edcae955ec60437bd55de3070c83",
	"baz/bat/bloopf": "2800577c6cda3f97123f9de49b40c0463e8f0c7435f88b66bdca622252dc8c05",
	"foo":            "544df7e68ae9ffcdb7d9d48a844214962812a8bb94eccd1d65fda808e4369ca0",
}

func TestAssetNames(t *testing.T) {
	names := main.AssetNames()
	if len(names) != len(BinsanityAssetNames) {
		t.Fatalf("Wrong number of names:\n  expected: %d\n    actual: %d",
			len(BinsanityAssetNames), len(names))
	}

	// ...moments when you really miss Testify... but NO deps for the
	// generated files!
	for idx, n := range names {
		if n != BinsanityAssetNames[idx] {
			t.Fatalf("Mismatch at %d:\n  expected: %s\n    actual: %s",
				idx, BinsanityAssetNames[idx], n)
		}
	}

}

func TestAssetNotFound(t *testing.T) {

	_, err := main.Asset(BinsanityAssetMissing)
	if err == nil {
		t.Fatal("No error for missing asset.")
	}
	if err.Error() != "Asset not found." {
		t.Fatal("Wrong error for missing asset.")
	}
}

func TestAssetFound(t *testing.T) {

	b, err := main.Asset(BinsanityAssetPresent)
	if err != nil {
		t.Fatal("Error for asset that should not be missing.")
	}
	sum := fmt.Sprintf("%x", sha256.Sum256(b))
	if sum != BinsanityAssetSum[BinsanityAssetPresent] {
		t.Fatal("Wrong sha256 sum for asset data.")
	}
}

func TestMustAssetNotFound(t *testing.T) {

	exp := "Asset not found."
	panicky := func() { main.MustAssetString(BinsanityAssetMissing) }
	AssertPanicsWith(t, panicky, exp, "MustAsset (not found)")

}

func TestMustAssetFound(t *testing.T) {

	b := main.MustAsset(BinsanityAssetPresent)
	sum := fmt.Sprintf("%x", sha256.Sum256(b))
	if sum != BinsanityAssetSum[BinsanityAssetPresent] {
		t.Fatal("Wrong sha256 sum for asset data.")
	}

}

func TestMustAssetStringNotFound(t *testing.T) {

	exp := "Asset not found."
	panicky := func() { main.MustAssetString(BinsanityAssetMissing) }
	AssertPanicsWith(t, panicky, exp, "MustAssetString (not found)")

}

func TestMustAssetStringFound(t *testing.T) {

	s := main.MustAssetString(BinsanityAssetPresent)
	sum := fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
	if sum != BinsanityAssetSum[BinsanityAssetPresent] {
		t.Fatal("Wrong sha256 sum for asset data.")
	}

}

func TestAssetSums(t *testing.T) {
	var want_tests bool
	// This is a little bit overkill but people have habits right?
	boolish := map[string]bool{
		"Y":    true,
		"YES":  true,
		"T":    true,
		"TRUE": true,
		"1":    true,
	}
	flag := strings.ToUpper(os.Getenv("BINSANITY_TEST_CONTENT"))
	want_tests = boolish[flag]
	if want_tests {
		for name, exp := range BinsanityAssetSum {
			b, err := main.Asset(name)
			if err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			sum := fmt.Sprintf("%x", sha256.Sum256(b))
			if sum != exp {
				t.Fatalf("Wrong sha256 sum for data of: %s\n  expected: %s\n    actual: %s",
					name, exp, sum)
			}
		}
	} else {
		t.Skip()
	}
}

// For a more useful version of this see: https://github.com/biztos/testig
func AssertPanicsWith(t *testing.T, f func(), exp string, msg string) {

	panicked := false
	got := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
				got = fmt.Sprintf("%s", r)
			}
		}()
		f()
	}()

	if !panicked {
		t.Fatalf("Function did not panic: %s", msg)
	} else if got != exp {

		t.Fatalf("Panic not as expected: %s\n  expected: %s\n    actual: %s",
			msg, exp, got)
	}

	// (In go testing, success is silent.)

}
