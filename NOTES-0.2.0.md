# Notes for Binsanity 0.2.0

What is the desired plan for package names?

Need two: the local package name for the declaration; and the full package
name for inclusion in the test files.

Local package name:

1. Have other .go file already in same dir? Use its package name.
2. Have go.mod file in current dir? Use main.
3. Walk up to find closest parent dir with a go.mod file. Use dirname.
4. Default: Use main.

Full package name for importing in test file:

    WAIT: why not write tests for main?  stupid, should have.
    Trick is import full package and then address it as main, use main_test.go

1. Have go.mod in current dir? Use package from that.
2. Walk up until we find a go.mod. Use package from that, plus subdirs.
3. Not found? Um... screwed, complain.

Programmatically then:

1. Find self or ancestor go.mod or fail.
2. Find package name in self dir or use dir if no go.mod or use main.

## WTF happens with versioned import for package name?

So if you have import path is foo.com/boo/v5 and your subpackage is... hm...
gonna need to trim off "v" statements I think? Or... well if you provide it
then it's your own fault and the go.mod shouldn't give it so maybe no problem?

## Hey what about using constants?

Problem is we would need to decode them, right? Because the constant can
only be a string, we need to create a Base64 (might as well compress it too)
of the item and then return that.

Fetching would be slower (how much?) _or_ loading would be slower (decode on
init) _and_ memory would be double (need the const, and the other thing)...
in order to avoid declaring the map. Plus we would need to have a function
with a huge switch to catch every constant because we can't eval, duh!

Seems like these are probably worse than giant maps, but I'm not sure.

TODO: benchmark a giant map of this stuff! binsanity $GOPATH or something.

```go
var pathHashes = map[string]string{} // foo/bar.txt -> hash(content)

const asfdafafsdafsdds = ...

```
