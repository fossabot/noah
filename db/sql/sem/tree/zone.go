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

package tree

// ZoneSpecifier represents a reference to a configurable zone of the keyspace.
type ZoneSpecifier struct {
	// Only one of NamedZone, Database or TableOrIndex may be set.
	NamedZone    UnrestrictedName
	Database     Name
	TableOrIndex TableNameWithIndex

	// Partition is only respected when Table is set.
	Partition Name
}

// TargetsTable returns whether the zone specifier targets a table or a subzone
// within a table.
func (node ZoneSpecifier) TargetsTable() bool {
	return node.NamedZone == "" && node.Database == ""
}

// TargetsIndex returns whether the zone specifier targets an index.
func (node ZoneSpecifier) TargetsIndex() bool {
	return node.TargetsTable() && (node.TableOrIndex.Index != "" || node.TableOrIndex.SearchTable)
}

// Format implements the NodeFormatter interface.
func (node *ZoneSpecifier) Format(ctx *FmtCtx) {
	if node.Partition != "" {
		ctx.WriteString("PARTITION ")
		ctx.FormatNode(&node.Partition)
		ctx.WriteString(" OF ")
	}
	if node.NamedZone != "" {
		ctx.WriteString("RANGE ")
		ctx.FormatNode(&node.NamedZone)
	} else if node.Database != "" {
		ctx.WriteString("DATABASE ")
		ctx.FormatNode(&node.Database)
	} else {
		if node.TargetsIndex() {
			ctx.WriteString("INDEX ")
		} else {
			ctx.WriteString("TABLE ")
		}
		ctx.FormatNode(&node.TableOrIndex)
	}
}

func (node *ZoneSpecifier) String() string { return AsString(node) }

// ShowZoneConfig represents an EXPERIMENTAL SHOW ZONE CONFIGURATION...
// statement.
type ShowZoneConfig struct {
	ZoneSpecifier
}

// Format implements the NodeFormatter interface.
func (node *ShowZoneConfig) Format(ctx *FmtCtx) {
	if node.ZoneSpecifier == (ZoneSpecifier{}) {
		ctx.WriteString("EXPERIMENTAL SHOW ZONE CONFIGURATIONS")
	} else {
		ctx.WriteString("EXPERIMENTAL SHOW ZONE CONFIGURATION FOR ")
		ctx.FormatNode(&node.ZoneSpecifier)
	}
}

// SetZoneConfig represents an ALTER DATABASE/TABLE... EXPERIMENTAL CONFIGURE
// ZONE statement.
type SetZoneConfig struct {
	ZoneSpecifier
	YAMLConfig Expr
}

// Format implements the NodeFormatter interface.
func (node *SetZoneConfig) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER ")
	ctx.FormatNode(&node.ZoneSpecifier)
	ctx.WriteString(" EXPERIMENTAL CONFIGURE ZONE ")
	ctx.FormatNode(node.YAMLConfig)
}
