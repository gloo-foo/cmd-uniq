package uniq_test

import (
	"fmt"

	"github.com/gloo-foo/testable"

	command "github.com/gloo-foo/cmd-uniq"
)

func ExampleUniq_count() {
	// printf 'apple\napple\nbanana\ncherry\ncherry\ncherry\n' | uniq -c
	output, _ := testable.Test(
		command.Uniq(command.UniqCount),
		"apple\napple\nbanana\ncherry\ncherry\ncherry\n",
	)
	fmt.Print(output)
	// Output:
	//       2 apple
	//       1 banana
	//       3 cherry
}

func ExampleUniq_ignoreCase() {
	// printf 'Apple\napple\nBanana\n' | uniq -i
	output, _ := testable.Test(
		command.Uniq(command.UniqIgnoreCase),
		"Apple\napple\nBanana\n",
	)
	fmt.Print(output)
	// Output:
	// Apple
	// Banana
}
