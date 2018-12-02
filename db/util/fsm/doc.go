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

/*
Package fsm provides an interface for defining and working with finite-state
machines.

The package is split into two main types: Transitions and Machine. Transitions
is an immutable State graph with Events acting as the directed edges between
different States. The graph is built by calling Compile on a Pattern, which is
meant to be done at init time. This pattern is a mapping from current States to
Events that may be applied on those states to resulting Transitions. The pattern
supports pattern matching on States and Events using wildcards and variable
bindings. To add new transitions to the graph, simply adjust the Pattern
provided to Compile. Transitions are not used directly after creation, instead,
they're used by Machine instances.

Machine is an instantiation of a finite-state machine. It is given a Transitions
graph when it is created to specify its State graph. Since the Transition graph
is itself state-less, multiple Machines can be powered by the same graph
simultaneously. The Machine has an Apply(Event) method, which applies the
provided event to its current state. This does two things:
1. It may move the current State to a new State, according to the Transitions
   graph.
2. It may apply an Action function on the Machine's ExtendedState, which is
   extra state in a Machine that does not contribute to state transition
   decisions, but that can be affected by a state transition.

See example_test.go for a full working example of a state machine with an
associated set of states and events.

This package encourages the Pattern to be declared as a map literal. When
declaring this literal, be careful to not declare two equal keys: they'll result
in the second overwriting the first with no warning because of how Go deals with
map literals. Note that keys that are not technically equal, but where one is a
superset of the other, will work as intended. E.g. the following is permitted:
 Compile(Pattern{
   stateOpen{retryIntent: Any} {
     eventTxnFinish{}: {...}
   }
   stateOpen{retryIntent: True} {
     eventRetriableErr{}: {...}
   }

Members of this package are accessed frequently when implementing a state
machine. For that reason, it is encouraged to dot-import this package in the
file with the transitions Pattern. The respective file should be kept small and
named <name>_fsm.go; our linter doesn't complain about dot-imports in such
files.

*/
package fsm
