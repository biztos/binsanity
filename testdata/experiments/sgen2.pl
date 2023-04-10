#!/usr/bin/env perl
#
# sgen2.pl - second try, using one big variable and lookup as currently doing.

use strict;
use warnings;

use feature 'say';

print <<_START;
// very long switch statement

package main

import (
	"fmt"
	"os"
	"sort"
)

func main() {
	for _, a := range os.Args[1:] {
		fmt.Println(len(doIt(a)))
	}
}

func doIt(n string) string {

	i := sort.SearchStrings(names, n)
	if i == len(names) || names[i] != n {
			return "dunno what"
	}
	return data[i]

}

var names = []string{
_START

my $num = shift @ARGV || 1;

for (1..$num) {
	say qq(\t"$_",)
};
print <<_MID;
}
var data = []string{
_MID
for (1..$num) {
	my $d = "$_ " x ($_ * 50);
	say qq(\t"$d",)
};
print <<_END;
}
_END


