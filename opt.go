package command

// uniqDuplicatesOnlyFlag enables -d: emit only groups that repeat.
type uniqDuplicatesOnlyFlag bool

const (
	UniqDuplicatesOnly   uniqDuplicatesOnlyFlag = true
	UniqNoDuplicatesOnly uniqDuplicatesOnlyFlag = false
)

func (f uniqDuplicatesOnlyFlag) Configure(flags *flags) { flags.duplicatesOnly = f }

// uniqUniqueOnlyFlag enables -u: emit only groups that do not repeat.
type uniqUniqueOnlyFlag bool

const (
	UniqUniqueOnly   uniqUniqueOnlyFlag = true
	UniqNoUniqueOnly uniqUniqueOnlyFlag = false
)

func (f uniqUniqueOnlyFlag) Configure(flags *flags) { flags.uniqueOnly = f }

// uniqCountFlag enables -c: prefix each emitted line with its occurrence count.
type uniqCountFlag bool

const (
	UniqCount   uniqCountFlag = true
	UniqNoCount uniqCountFlag = false
)

func (f uniqCountFlag) Configure(flags *flags) { flags.count = f }

// uniqIgnoreCaseFlag enables -i: compare lines case-insensitively.
type uniqIgnoreCaseFlag bool

const (
	UniqIgnoreCase   uniqIgnoreCaseFlag = true
	UniqNoIgnoreCase uniqIgnoreCaseFlag = false
)

func (f uniqIgnoreCaseFlag) Configure(flags *flags) { flags.ignoreCase = f }

type flags struct {
	duplicatesOnly uniqDuplicatesOnlyFlag
	uniqueOnly     uniqUniqueOnlyFlag
	count          uniqCountFlag
	ignoreCase     uniqIgnoreCaseFlag
}
