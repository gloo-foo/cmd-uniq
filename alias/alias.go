// Package alias provides unprefixed names for the uniq command and its flags.
//
//	import uniq "github.com/gloo-foo/cmd-uniq/alias"
//	uniq.Uniq(uniq.DuplicatesOnly)
package alias

import command "github.com/gloo-foo/cmd-uniq"

// Uniq re-exports the constructor.
var Uniq = command.Uniq

// -d flag: emit only repeated groups
const DuplicatesOnly = command.UniqDuplicatesOnly

// default: do not restrict to repeated groups
const NoDuplicatesOnly = command.UniqNoDuplicatesOnly

// -u flag: emit only non-repeated groups
const UniqueOnly = command.UniqUniqueOnly

// default: do not restrict to non-repeated groups
const NoUniqueOnly = command.UniqNoUniqueOnly

// -c flag: prefix each line with its occurrence count
const Count = command.UniqCount

// default: do not prefix counts
const NoCount = command.UniqNoCount

// -i flag: compare lines case-insensitively
const IgnoreCase = command.UniqIgnoreCase

// default: compare lines case-sensitively
const NoIgnoreCase = command.UniqNoIgnoreCase
