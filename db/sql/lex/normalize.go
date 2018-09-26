/*
 * Copyright (c) 2018 Ready Stock
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 *
 * This application uses Open Source components. You can find the
 * source code of their open source projects along with license
 * information below. We acknowledge and are grateful to these
 * developers for their contributions to open source.
 *
 * Project: CockroachDB https://github.com/cockroachdb/cockroach
 * Copyright 2018 The Cockroach Authors.
 * License (Apache License 2.0) https://github.com/cockroachdb/cockroach/blob/master/LICENSE
 *
 * Project: Vitess https://github.com/vitessio/vitess
 * Copyright 2018 Google Inc.
 * License (Apache License 2.0) https://github.com/vitessio/vitess/blob/master/LICENSE
 *
 * Project: Citus https://github.com/citusdata/citus
 * Copyright 2018 Citus Data, Inc.
 * License (GNU Affero General Public License v3.0) https://github.com/citusdata/citus/blob/master/LICENSE
 *
 * Project: pg_query_go https://github.com/lfittl/pg_query_go
 * Copyright 2018 Lukas Fittl
 * License (3-Clause BSD) https://github.com/lfittl/pg_query_go/blob/master/LICENSE
 *
 * Project: pgx https://github.com/jackc/pgx
 * Copyright 2018 Jack Christensen
 * License (MIT) https://github.com/jackc/pgx/blob/master/LICENSE
 *
 * Project: BadgerDB https://github.com/dgraph-io/badger
 * Copyright 2018 Dgraph Labs, Inc. and Contributors
 * License (MIT) https://github.com/dgraph-io/badger/blob/master/LICENSE
 */

package lex

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// Special case normalization rules for Turkish/Azeri lowercase dotless-i and
// uppercase dotted-i. Fold both dotted and dotless 'i' into the ascii i/I, so
// our case-insensitive comparison functions can be locale-invariant. This
// mapping implements case-insensitivity for Turkish and other latin-derived
// languages simultaneously, with the additional quirk that it is also
// insensitive to the dottedness of the i's
var normalize = unicode.SpecialCase{
	unicode.CaseRange{
		Lo: 0x0130,
		Hi: 0x0130,
		Delta: [unicode.MaxCase]rune{
			0x49 - 0x130, // Upper
			0x69 - 0x130, // Lower
			0x49 - 0x130, // Title
		},
		// Delta: [unicode.MaxCase]rune{
		// 	0x49 - 0x130, // Upper
		// 	0x69 - 0x130, // Lower
		// 	0x49 - 0x130, // Title
		// },
	},
	unicode.CaseRange{
		Lo: 0x0131,
		Hi: 0x0131,
		Delta: [unicode.MaxCase]rune{
			0x49 - 0x131, // Upper
			0x69 - 0x131, // Lower
			0x49 - 0x131, // Title
		},
	},
}

// NormalizeName normalizes to lowercase and Unicode Normalization
// Form C (NFC).
func NormalizeName(n string) string {
	lower := strings.Map(normalize.ToLower, n)
	if isASCII(lower) {
		return lower
	}
	return norm.NFC.String(lower)
}
