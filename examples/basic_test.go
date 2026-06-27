package uniq_test

import (
	"fmt"

	command "github.com/gloo-foo/cmd-uniq"
	"github.com/gloo-foo/testable"
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
