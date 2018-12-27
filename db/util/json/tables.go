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
 * Copyright 2018  'pq' Contributors Portions Copyright (C) 2018 Blake Mizerany
 * License https://github.com/lib/pq/blob/master/LICENSE.md
 *
 * Project: go-linq github.com/ahmetb/go-linq
 * Copyright 2018 Ahmet Alp Balkan
 * License https://github.com/ahmetb/go-linq/blob/master/LICENSE
 */

// Copied from https://github.com/golang/go/blob/master/src/encoding/json/tables.go

package json

import "unicode/utf8"

// safeSet holds the value true if the ASCII character with the given array
// position can be represented inside a JSON string without any further
// escaping.
//
// All values are true except for the ASCII control characters (0-31), the
// double quote ("), and the backslash character ("\").
var safeSet = [utf8.RuneSelf]bool{
    ' ':      true,
    '!':      true,
    '"':      false,
    '#':      true,
    '$':      true,
    '%':      true,
    '&':      true,
    '\'':     true,
    '(':      true,
    ')':      true,
    '*':      true,
    '+':      true,
    ',':      true,
    '-':      true,
    '.':      true,
    '/':      true,
    '0':      true,
    '1':      true,
    '2':      true,
    '3':      true,
    '4':      true,
    '5':      true,
    '6':      true,
    '7':      true,
    '8':      true,
    '9':      true,
    ':':      true,
    ';':      true,
    '<':      true,
    '=':      true,
    '>':      true,
    '?':      true,
    '@':      true,
    'A':      true,
    'B':      true,
    'C':      true,
    'D':      true,
    'E':      true,
    'F':      true,
    'G':      true,
    'H':      true,
    'I':      true,
    'J':      true,
    'K':      true,
    'L':      true,
    'M':      true,
    'N':      true,
    'O':      true,
    'P':      true,
    'Q':      true,
    'R':      true,
    'S':      true,
    'T':      true,
    'U':      true,
    'V':      true,
    'W':      true,
    'X':      true,
    'Y':      true,
    'Z':      true,
    '[':      true,
    '\\':     false,
    ']':      true,
    '^':      true,
    '_':      true,
    '`':      true,
    'a':      true,
    'b':      true,
    'c':      true,
    'd':      true,
    'e':      true,
    'f':      true,
    'g':      true,
    'h':      true,
    'i':      true,
    'j':      true,
    'k':      true,
    'l':      true,
    'm':      true,
    'n':      true,
    'o':      true,
    'p':      true,
    'q':      true,
    'r':      true,
    's':      true,
    't':      true,
    'u':      true,
    'v':      true,
    'w':      true,
    'x':      true,
    'y':      true,
    'z':      true,
    '{':      true,
    '|':      true,
    '}':      true,
    '~':      true,
    '\u007f': true,
}
