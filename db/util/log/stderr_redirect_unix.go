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

// +build !windows

package log

import (
	"os"

	"golang.org/x/sys/unix"
)

func dupFD(fd uintptr) (uintptr, error) {
	// Warning: failing to set FD_CLOEXEC causes the duplicated file descriptor
	// to leak into subprocesses created by exec.Command. If the file descriptor
	// is a pipe, these subprocesses will hold the pipe open (i.e., prevent
	// EOF), potentially beyond the lifetime of this process.
	//
	// This can break go test's timeouts. go test usually spawns a test process
	// with its stdin and stderr streams hooked up to pipes; if the test process
	// times out, it sends a SIGKILL and attempts to read stdin and stderr to
	// completion. If the test process has itself spawned long-lived
	// subprocesses that hold references to the stdin or stderr pipes, go test
	// will hang until the subprocesses exit, rather defeating the purpose of
	// a timeout.
	nfd, _, errno := unix.Syscall(unix.SYS_FCNTL, fd, unix.F_DUPFD_CLOEXEC, 0)
	if errno != 0 {
		return 0, errno
	}
	return nfd, nil
}

func redirectStderr(f *os.File) error {
	return unix.Dup2(int(f.Fd()), unix.Stderr)
}
