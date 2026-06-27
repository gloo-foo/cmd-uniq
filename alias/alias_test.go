package alias_test

import (
	"slices"
	"testing"

	uniq "github.com/gloo-foo/cmd-uniq/alias"
	"github.com/gloo-foo/testable"
)

// The alias package re-exports the constructor and flag constants under
// unprefixed names. A mis-wired re-export (say, UniqueOnly bound to the
// duplicates constant, or Uniq bound to the wrong function) compiles cleanly, so
// only behavior can prove the wiring. Each test exercises one re-export and
// asserts the GNU uniq output it must produce.
//
// input groups: A (x3), b (x1), C (x2). Choosing distinct group sizes makes the
// -d / -u / -c selectors observably different from one another and from default.
const input = "A\nA\nA\nb\nC\nC\n"

func assertLines(t *testing.T, got, want []string) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func run(t *testing.T, want []string, opts ...any) {
	t.Helper()
	lines, err := testable.TestLines(uniq.Uniq(opts...), input)
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, lines, want)
}

func TestAlias_DefaultCollapsesEachGroup(t *testing.T) {
	run(t, []string{"A", "b", "C"})
}

func TestAlias_DuplicatesOnlyKeepsRepeatedGroups(t *testing.T) {
	// -d: only A (x3) and C (x2) repeat; the singleton b is dropped.
	run(t, []string{"A", "C"}, uniq.DuplicatesOnly)
}

func TestAlias_UniqueOnlyKeepsNonRepeatedGroups(t *testing.T) {
	// -u: only the singleton b survives.
	run(t, []string{"b"}, uniq.UniqueOnly)
}

func TestAlias_CountPrefixesEachGroup(t *testing.T) {
	// -c: GNU's width-7 right-justified count, one space, then the line.
	run(t, []string{"      3 A", "      1 b", "      2 C"}, uniq.Count)
}

func TestAlias_IgnoreCaseFoldsAdjacentCasing(t *testing.T) {
	// -i: "a" and "A" are adjacent equal, collapsing to the first ("A").
	lines, err := testable.TestLines(uniq.Uniq(uniq.IgnoreCase), "A\na\nb\n")
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, lines, []string{"A", "b"})
}

func TestAlias_DisabledFlagsMatchDefault(t *testing.T) {
	// The No* constants are the disabled forms: they must behave exactly like
	// passing no flag at all.
	run(t, []string{"A", "b", "C"},
		uniq.NoDuplicatesOnly, uniq.NoUniqueOnly, uniq.NoCount, uniq.NoIgnoreCase)
}
