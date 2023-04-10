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

**SOLVED** just append after the /v5 and Go figures it out.

So if you have import path is foo.com/boo/v5 and your subpackage is... hm...
gonna need to trim off "v" statements I think? Or... well if you provide it
then it's your own fault and the go.mod shouldn't give it so maybe no problem?

## Hey what about using constants?

**CHECKED BUT NOT PROVEN** that it doesn't make a noticeable difference in
the binary size nor in the perceived runtime with a pretty big set -- more
than you should use this kind of tool for really. Build is slow, execute is
quite fast. Ergo: **no obvious reason to use constants at this time.**

---

Question is, do we have any kind of memory advantage (or other advantage) if
we do the big file strings as constants instead of a giant-ass array. Could
bench it I guess.

To do constants would be some like this:

```go

const d1 = "base64 encoded stuff"
const d2 = "more base64 encoded stuff"

func get(name string) (string,error) {
    switch name: {
        case "d1":
            return d1, nil
        case "d2":
            return d2, nil
        default:
            return "", errors.New("not found")
    }
}
```

The obviously bad thing here is we then have a massive switch statement,
because there's no way to address the constant by a variable name.
