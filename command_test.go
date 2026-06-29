package command_test

import (
	"strings"
	"testing"

	"github.com/gloo-foo/testable"
	"github.com/gloo-foo/testable/assertion"

	command "github.com/gloo-foo/cmd-uniq"
)

func lines(t *testing.T, input string, opts ...any) []string {
	t.Helper()
	got, err := testable.TestLines(command.Uniq(opts...), input)
	assertion.NoError(t, err)
	return got
}

// ==============================================================================
// Default: collapse runs of ADJACENT duplicates only
// ==============================================================================

func TestUniq_CollapsesAdjacentRuns(t *testing.T) {
	got := lines(t, "a\nb\nb\nc\nb\nb\nb\nd\n")
	assertion.Lines(t, got, []string{"a", "b", "c", "b", "d"})
}

func TestUniq_PreservesNonAdjacentDuplicates(t *testing.T) {
	// The two "a" lines are not adjacent, so both must survive — this is the
	// defining property of uniq versus a global dedup.
	got := lines(t, "a\nb\na\n")
	assertion.Lines(t, got, []string{"a", "b", "a"})
}

func TestUniq_AllSameCollapseToOne(t *testing.T) {
	got := lines(t, "x\nx\nx\nx\n")
	assertion.Lines(t, got, []string{"x"})
}

func TestUniq_AllUniquePassThrough(t *testing.T) {
	got := lines(t, "a\nb\nc\nd\n")
	assertion.Lines(t, got, []string{"a", "b", "c", "d"})
}

func TestUniq_SingleLine(t *testing.T) {
	assertion.Lines(t, lines(t, "only\n"), []string{"only"})
}

func TestUniq_EmptyInput(t *testing.T) {
	assertion.Empty(t, lines(t, ""))
}

// ==============================================================================
// -d: emit only repeated groups
// ==============================================================================

func TestUniq_DuplicatesOnly(t *testing.T) {
	got := lines(t, "a\nb\nb\nc\nb\nb\nb\nd\n", command.UniqDuplicatesOnly)
	assertion.Lines(t, got, []string{"b", "b"})
}

func TestUniq_DuplicatesOnly_NoneRepeat(t *testing.T) {
	assertion.Empty(t, lines(t, "a\nb\nc\n", command.UniqDuplicatesOnly))
}

func TestUniq_DuplicatesOnly_EmptyInput(t *testing.T) {
	assertion.Empty(t, lines(t, "", command.UniqDuplicatesOnly))
}

// ==============================================================================
// -u: emit only non-repeated groups
// ==============================================================================

func TestUniq_UniqueOnly(t *testing.T) {
	// a(1) b(2) c(1): only the singletons a and c survive.
	got := lines(t, "a\nb\nb\nc\n", command.UniqUniqueOnly)
	assertion.Lines(t, got, []string{"a", "c"})
}

func TestUniq_UniqueOnly_AllRepeat(t *testing.T) {
	assertion.Empty(t, lines(t, "x\nx\ny\ny\n", command.UniqUniqueOnly))
}

func TestUniq_DuplicatesAndUniqueTogetherSelectNothing(t *testing.T) {
	// GNU uniq: -d and -u are complementary; together they match no group.
	got := lines(t, "a\nb\nb\nc\n", command.UniqDuplicatesOnly, command.UniqUniqueOnly)
	assertion.Empty(t, got)
}

// ==============================================================================
// -c: prefix each emitted line with its occurrence count (width-7)
// ==============================================================================

func TestUniq_Count(t *testing.T) {
	got := lines(t, "a\nb\nb\nc\nc\nc\n", command.UniqCount)
	assertion.Lines(t, got, []string{"      1 a", "      2 b", "      3 c"})
}

func TestUniq_Count_MultiDigitWidth(t *testing.T) {
	// Ten adjacent copies: the count "10" right-justifies into the width-7 field.
	got := lines(t, strings.Repeat("z\n", 10), command.UniqCount)
	assertion.Lines(t, got, []string{"     10 z"})
}

func TestUniq_CountWithDuplicatesOnly(t *testing.T) {
	// -c composes with -d: only repeated groups, each count-prefixed.
	got := lines(t, "a\nb\nb\n", command.UniqCount, command.UniqDuplicatesOnly)
	assertion.Lines(t, got, []string{"      2 b"})
}

// ==============================================================================
// -i: case-insensitive comparison
// ==============================================================================

func TestUniq_IgnoreCase(t *testing.T) {
	// "Hello"/"hello"/"HELLO" are adjacent equal under folding; first survives.
	got := lines(t, "Hello\nhello\nHELLO\nworld\n", command.UniqIgnoreCase)
	assertion.Lines(t, got, []string{"Hello", "world"})
}

func TestUniq_CaseSensitiveByDefault(t *testing.T) {
	// Without -i, differing case is not a duplicate.
	got := lines(t, "Hello\nhello\n")
	assertion.Lines(t, got, []string{"Hello", "hello"})
}

func TestUniq_IgnoreCaseWithCount(t *testing.T) {
	got := lines(t, "Hi\nhi\nHI\n", command.UniqIgnoreCase, command.UniqCount)
	assertion.Lines(t, got, []string{"      3 Hi"})
}
