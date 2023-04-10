#!/usr/bin/env perl
#
# sgen.pl - switch generator for testing my weird idea.
#
# sometimes perl is still the king yo...

use strict;
use warnings;

use feature 'say';

print <<_START;
// very long switch statement

package main

import (
	"fmt"
	"os"
)

func main() {
	for _, a := range os.Args[1:] {
		fmt.Println(len(doIt(a)))
	}
}

func doIt(n string) string {
	switch n {
_START

my $num = shift @ARGV || 1;

for (1..$num) {
	say qq(\tcase "$_":\n\t\treturn c$_)
};
print <<_END;
	default:
		return "dunno what"
	}
}
_END

# now stick on the constants.
for (1..$num) {
	my $d = "$_ " x ($_ * 50);
	say qq(const c$_ = "$d")
};

