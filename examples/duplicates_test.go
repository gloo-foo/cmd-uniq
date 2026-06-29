package uniq_test

import (
	"fmt"

	"github.com/gloo-foo/testable"

	command "github.com/gloo-foo/cmd-uniq"
)

func ExampleUniq_duplicatesOnly() {
	// printf 'apple\napple\nbanana\ncherry\ncherry\n' | uniq -d
	output, _ := testable.Test(
		command.Uniq(command.UniqDuplicatesOnly),
		"apple\napple\nbanana\ncherry\ncherry\n",
	)
	fmt.Print(output)
	// Output:
	// apple
	// cherry
}

func ExampleUniq_uniqueOnly() {
	// printf 'apple\napple\nbanana\ncherry\ncherry\n' | uniq -u
	output, _ := testable.Test(
		command.Uniq(command.UniqUniqueOnly),
		"apple\napple\nbanana\ncherry\ncherry\n",
	)
	fmt.Print(output)
	// Output:
	// banana
}
