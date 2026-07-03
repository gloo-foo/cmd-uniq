package uniq_test

import (
	"fmt"

	"github.com/gloo-foo/testable"

	command "github.com/gloo-foo/cmd-uniq"
)

func ExampleUniq_basic() {
	// echo "apple\napple\nbanana\nbanana\napple" | uniq
	output, _ := testable.Test(
		command.Uniq(),
		"apple\napple\nbanana\nbanana\napple\n",
	)
	fmt.Print(output)
	// Output:
	// apple
	// banana
	// apple
}
