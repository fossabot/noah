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
 *
 * Project: Sonyflake https://github.com/sony/sonyflake
 * Copyright 2018 Sony Corporation
 * License (MIT) https://github.com/sony/sonyflake/blob/master/LICENSE
 *
 * Project: Raft https://github.com/hashicorp/raft
 * Copyright 2018 HashiCorp
 * License (MPL-2.0) https://github.com/hashicorp/raft/blob/master/LICENSE
 *
 * Project: pq github.com/lib/pq
 * Copyright 2018  'pq' Contributors Portions Copyright (C) 2011 Blake Mizerany
 * License https://github.com/lib/pq/blob/master/LICENSE.md
 */

package lex

import (
    "unicode"
    "unicode/utf8"
)

// isASCII returns true if all the characters in s are ASCII.
func isASCII(s string) bool {
    for _, c := range s {
        if c > unicode.MaxASCII {
            return false
        }
    }
    return true
}

// IsDigit returns true if the character is between 0 and 9.
func IsDigit(ch int) bool {
    return ch >= '0' && ch <= '9'
}

// IsHexDigit returns true if the character is a valid hexadecimal digit.
func IsHexDigit(ch int) bool {
    return (ch >= '0' && ch <= '9') ||
        (ch >= 'a' && ch <= 'f') ||
        (ch >= 'A' && ch <= 'F')
}

// lookaheadKeywords are those keywords for which we need one token
// of lookahead extra to determine their token type.
var lookaheadKeywords = map[string]struct{}{
    "between":    {},
    "ilike":      {},
    "in":         {},
    "like":       {},
    "of":         {},
    "ordinality": {},
    "similar":    {},
    "time":       {},
}

var reservedKeywords = map[string]struct{}{}

// isReservedKeyword returns true if the keyword is reserved, or needs
// one extra token of lookahead.
func isReservedKeyword(s string) bool {
    if _, ok := reservedKeywords[s]; ok {
        return true
    }
    if _, ok := lookaheadKeywords[s]; ok {
        return true
    }
    return false
}

// isBareIdentifier returns true if the input string is a permissible bare SQL
// identifier.
func isBareIdentifier(s string) bool {
    if len(s) == 0 || !IsIdentStart(int(s[0])) || (s[0] >= 'A' && s[0] <= 'Z') {
        return false
    }
    // Keep track of whether the input string is all ASCII. If it is, we don't
    // have to bother running the full Normalize() function at the end, which is
    // quite expensive.
    isASCII := s[0] < utf8.RuneSelf
    for i := 1; i < len(s); i++ {
        if !IsIdentMiddle(int(s[i])) {
            return false
        }
        if s[i] >= 'A' && s[i] <= 'Z' {
            // Non-lowercase identifiers aren't permissible.
            return false
        }
        if s[i] >= utf8.RuneSelf {
            isASCII = false
        }
    }
    return isASCII || NormalizeName(s) == s
}

// IsIdentStart returns true if the character is valid at the start of an identifier.
func IsIdentStart(ch int) bool {
    return (ch >= 'A' && ch <= 'Z') ||
        (ch >= 'a' && ch <= 'z') ||
        (ch >= 128 && ch <= 255) ||
        (ch == '_')
}

// IsIdentMiddle returns true if the character is valid inside an identifier.
func IsIdentMiddle(ch int) bool {
    return IsIdentStart(ch) || IsDigit(ch) || ch == '$'
}
