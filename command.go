// Package command implements the uniq filter: it collapses each run of
// adjacent identical lines into a single line, matching GNU uniq semantics.
//
// Only ADJACENT duplicates are collapsed — non-adjacent repeats are preserved,
// exactly as the Unix tool behaves (uniq is typically preceded by sort).
package command

import (
	"bytes"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// runLength is the number of adjacent identical lines in a single group.
type runLength uint64

// group is one run of adjacent lines that compare equal. rep is the
// representative line emitted for the group (the first line of the run).
type group struct {
	rep   []byte
	count runLength
}

// Uniq returns a Command that collapses runs of adjacent duplicate lines.
//
// Flags:
//   - UniqDuplicatesOnly (-d): emit only groups that repeat (count > 1)
//   - UniqUniqueOnly (-u): emit only groups that do not repeat (count == 1)
//   - UniqCount (-c): prefix each emitted line with its occurrence count
//   - UniqIgnoreCase (-i): compare lines case-insensitively
//
// -d and -u are complementary selectors; supplying both selects nothing, which
// is GNU uniq's behavior.
func Uniq(opts ...any) gloo.Command[[]byte, []byte] {
	f, rest := fold(opts)
	gloo.NewParameters[gloo.File, struct{}](rest...)
	render := renderer(f)
	return patterns.Accumulate(func(lines [][]byte) ([][]byte, error) {
		return collapse(lines, f.ignoreCaseEnabled, render), nil
	})
}

// collapse folds adjacent equal lines into groups, then emits each selected
// group through render.
func collapse(lines [][]byte, isIgnoreCase uniqIgnoreCaseFlag, render func(group) ([]byte, bool)) [][]byte {
	out := make([][]byte, 0, len(lines))
	for _, g := range groupLines(lines, isIgnoreCase) {
		if rendered, ok := render(g); ok {
			out = append(out, rendered)
		}
	}
	return out
}

// groupLines partitions lines into runs of adjacent equal lines.
func groupLines(lines [][]byte, isIgnoreCase uniqIgnoreCaseFlag) []group {
	groups := make([]group, 0, len(lines))
	for _, line := range lines {
		if n := len(groups); n > 0 && sameLine(groups[n-1].rep, line, isIgnoreCase) {
			groups[n-1].count++
			continue
		}
		groups = append(groups, group{rep: line, count: 1})
	}
	return groups
}

// sameLine reports whether two lines compare equal under the chosen casing.
func sameLine(a, b []byte, isIgnoreCase uniqIgnoreCaseFlag) bool {
	if isIgnoreCase {
		return bytes.EqualFold(a, b)
	}
	return bytes.Equal(a, b)
}

// renderer builds the per-group rendering function for the active flags. It
// composes a selection predicate (which groups survive) with a formatter (how a
// surviving group is written).
func renderer(f flags) func(group) ([]byte, bool) {
	keep := selector(f)
	format := formatter(f.countEnabled)
	return func(g group) ([]byte, bool) {
		if !keep(g) {
			return nil, false
		}
		return format(g), true
	}
}

// selector returns the predicate deciding whether a group is emitted. -d and -u
// are independent constraints applied conjunctively, so supplying both keeps
// only groups that are simultaneously repeated and unique — i.e. none — which is
// GNU uniq's behavior.
func selector(f flags) func(group) bool {
	keep := keepAll
	if bool(f.duplicatesOnlyEnabled) {
		keep = and(keep, isRepeated)
	}
	if bool(f.uniqueOnlyEnabled) {
		keep = and(keep, isUnique)
	}
	return keep
}

// and composes two group predicates into their conjunction.
func and(left, right func(group) bool) func(group) bool {
	return func(g group) bool { return left(g) && right(g) }
}

// isRepeated reports whether a group has more than one line (-d).
func isRepeated(g group) bool { return g.count > 1 }

// isUnique reports whether a group has exactly one line (-u).
func isUnique(g group) bool { return g.count == 1 }

// keepAll selects every group (default behavior).
func keepAll(group) bool { return true }

// formatter returns the function that renders a surviving group's line, with or
// without the -c count prefix.
func formatter(countEnabled uniqCountFlag) func(group) []byte {
	if bool(countEnabled) {
		return countLine
	}
	return plainLine
}

// plainLine renders a group as its representative line, unmodified.
func plainLine(g group) []byte { return g.rep }

// countLine renders a group as GNU uniq's "%7d %s" count-prefixed line.
func countLine(g group) []byte {
	prefix := make([]byte, 0, len(g.rep)+8)
	prefix = appendCount(prefix, g.count)
	prefix = append(prefix, ' ')
	return append(prefix, g.rep...)
}

// appendCount writes a right-justified, width-7 decimal count, matching the
// field width GNU uniq uses for its -c prefix.
func appendCount(dst []byte, n runLength) []byte {
	const width = 7
	digits := decimal(n)
	for i := len(digits); i < width; i++ {
		dst = append(dst, ' ')
	}
	return append(dst, digits...)
}

// decimal renders n as its base-10 ASCII digits. A group always contains at
// least one line, so n is always >= 1; the loop emits at least one digit.
func decimal(n runLength) []byte {
	var buf [20]byte
	i := len(buf)
	for {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
		if n == 0 {
			return buf[i:]
		}
	}
}
