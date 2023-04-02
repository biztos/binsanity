// main_test.go -- use Example to test main(), just for instance.
//
// it's generally shitty we have to resort to this, and might not be safe,
// and anyway testing main() is not our problem except we will have to do it
// for the cmd, so might as well get back in shape.

package main

func Example() {

	main()

	// Output:
	// foo is foo
	//
	// bar is bar
	//
	// baz is bat is bloopf
	//
	// For doobie:  Asset not found.
}
