/* binsanity_test.go - auto-generated; edit at your own peril!

To test the checksums for all content, set the environment variable
BINSANITY_TEST_CONTENT to one of: Y,YES,T,TRUE,1 (the Truthy Shortlist).

More info: https://github.com/biztos/binsanity

Generated: 0001-01-01 00:00:00 +0000 UTC

*/

package binsanity_test

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/biztos/binsanity"
)

const BinsanityAssetMissing = "tests.tmpl--NOPE"
const BinsanityAssetPresent = "tests.tmpl"
const BinsanityAssetPresentSum = "4faaed29eb3c52f9d587bbb0aadad8690dfe95e07c14fc9a3eff8b216edfcd44"

var BinsanityAssetNames = []string{

	"code.tmpl",
	"tests.tmpl",
}

var BinsanityAssetSums = []string{
	"dad22db46328ef37e56b714dde7669f6fc7a686c95d028ae970fff52a8425b08",
	"4faaed29eb3c52f9d587bbb0aadad8690dfe95e07c14fc9a3eff8b216edfcd44",
}

func TestAssetNames(t *testing.T) {

	names := binsanity.AssetNames()
	if len(names) != len(BinsanityAssetNames) {
		t.Fatalf("Wrong number of names:\n  expected: %d\n  actual: %d",
			len(BinsanityAssetNames), len(names))
	}

	// ...moments when you really miss Testify... but NO deps for the
	// generated files!
	for idx, n := range names {
		if n != BinsanityAssetNames[idx] {
			t.Fatalf("Mismatch at %d:\n  expected: %s\n  actual: %s",
				idx, BinsanityAssetNames[idx], n)
		}
	}

}

func TestAssetNotFound(t *testing.T) {

	_, err := binsanity.Asset(BinsanityAssetMissing)
	if err == nil {
		t.Fatal("No error for missing asset.")
	}
	if err.Error() != "Asset not found." {
		t.Fatal("Wrong error for missing asset.")
	}
}

func TestAssetFound(t *testing.T) {

	b, err := binsanity.Asset(BinsanityAssetPresent)
	if err != nil {
		t.Fatal("Error for asset that should not be missing.")
	}
	sum := fmt.Sprintf("%x", sha256.Sum256(b))
	if sum != BinsanityAssetPresentSum {
		t.Fatal("Wrong sha256 sum for asset data.")
	}
}

func TestMustAssetNotFound(t *testing.T) {

	exp := "Asset not found."
	panicky := func() { binsanity.MustAssetString(BinsanityAssetMissing) }
	AssertPanicsWith(t, panicky, exp, "MustAsset (not found)")

}

func TestMustAssetFound(t *testing.T) {

	b := binsanity.MustAsset(BinsanityAssetPresent)
	sum := fmt.Sprintf("%x", sha256.Sum256(b))
	if sum != BinsanityAssetPresentSum {
		t.Fatal("Wrong sha256 sum for asset data.")
	}

}

func TestMustAssetStringNotFound(t *testing.T) {

	exp := "Asset not found."
	panicky := func() { binsanity.MustAssetString(BinsanityAssetMissing) }
	AssertPanicsWith(t, panicky, exp, "MustAssetString (not found)")

}

func TestMustAssetStringFound(t *testing.T) {

	s := binsanity.MustAssetString(BinsanityAssetPresent)
	sum := fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
	if sum != BinsanityAssetPresentSum {
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
	if !want_tests {
		t.Skip()
		return
	}
	for idx, name := range BinsanityAssetNames {
		b, err := binsanity.Asset(name)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		exp := BinsanityAssetSums[idx]
		sum := fmt.Sprintf("%x", sha256.Sum256(b))
		if sum != exp {
			t.Fatalf("Wrong sha256 sum for data of: %s\n  expected: %s\n    actual: %s",
				name, exp, sum)
		}
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
