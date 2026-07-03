package command

// uniqDuplicatesOnlyFlag enables -d: emit only groups that repeat.
type uniqDuplicatesOnlyFlag bool

const (
	UniqDuplicatesOnly   uniqDuplicatesOnlyFlag = true
	UniqNoDuplicatesOnly uniqDuplicatesOnlyFlag = false
)

// uniqUniqueOnlyFlag enables -u: emit only groups that do not repeat.
type uniqUniqueOnlyFlag bool

const (
	UniqUniqueOnly   uniqUniqueOnlyFlag = true
	UniqNoUniqueOnly uniqUniqueOnlyFlag = false
)

// uniqCountFlag enables -c: prefix each emitted line with its occurrence count.
type uniqCountFlag bool

const (
	UniqCount   uniqCountFlag = true
	UniqNoCount uniqCountFlag = false
)

// uniqIgnoreCaseFlag enables -i: compare lines case-insensitively.
type uniqIgnoreCaseFlag bool

const (
	UniqIgnoreCase   uniqIgnoreCaseFlag = true
	UniqNoIgnoreCase uniqIgnoreCaseFlag = false
)

type flags struct {
	duplicatesOnlyEnabled uniqDuplicatesOnlyFlag
	uniqueOnlyEnabled     uniqUniqueOnlyFlag
	countEnabled          uniqCountFlag
	ignoreCaseEnabled     uniqIgnoreCaseFlag
}

// fold partitions opts: uniq's own option values are folded into the flag set,
// and every other argument is passed through unchanged for the framework's
// positional classifier.
func fold(opts []any) (flags, []any) {
	var f flags
	rest := make([]any, 0, len(opts))
	for _, o := range opts {
		switch v := o.(type) {
		case uniqDuplicatesOnlyFlag:
			f.duplicatesOnlyEnabled = v
		case uniqUniqueOnlyFlag:
			f.uniqueOnlyEnabled = v
		case uniqCountFlag:
			f.countEnabled = v
		case uniqIgnoreCaseFlag:
			f.ignoreCaseEnabled = v
		default:
			rest = append(rest, o)
		}
	}
	return f, rest
}
