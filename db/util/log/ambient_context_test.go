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

package log

import (
	"context"
	"testing"

    "github.com/opentracing/opentracing-go"

	"github.com/cockroachdb/cockroach/pkg/util/tracing"
)

func TestAnnotateCtxTags(t *testing.T) {
	ac := AmbientContext{}
	ac.AddLogTag("a", 1)
	ac.AddLogTag("b", 2)

	ctx := ac.AnnotateCtx(context.Background())
	if exp, val := "[a1,b2] test", MakeMessage(ctx, "test", nil); val != exp {
		t.Errorf("expected '%s', got '%s'", exp, val)
	}

	ctx = context.Background()
	ctx = WithLogTag(ctx, "a", 10)
	ctx = WithLogTag(ctx, "aa", nil)
	ctx = ac.AnnotateCtx(ctx)

	if exp, val := "[aa,a1,b2] test", MakeMessage(ctx, "test", nil); val != exp {
		t.Errorf("expected '%s', got '%s'", exp, val)
	}
}

func TestAnnotateCtxSpan(t *testing.T) {
	tracer := tracing.NewTracer()
	tracer.SetForceRealSpans(true)

	ac := AmbientContext{Tracer: tracer}
	ac.AddLogTag("ambient", nil)

	// Annotate a context that has an open span.

	sp1 := tracer.StartSpan("root")
	tracing.StartRecording(sp1, tracing.SingleNodeRecording)
	ctx1 := opentracing.ContextWithSpan(context.Background(), sp1)
	Event(ctx1, "a")

	ctx2, sp2 := ac.AnnotateCtxWithSpan(ctx1, "child")
	Event(ctx2, "b")

	Event(ctx1, "c")
	sp2.Finish()
	sp1.Finish()

	if err := tracing.TestingCheckRecordedSpans(tracing.GetRecording(sp1), `
		span root:
			event: a
			event: c
		span child:
			event: [ambient] b
	`); err != nil {
		t.Fatal(err)
	}

	// Annotate a context that has no span.

	ac.Tracer = tracer
	ctx, sp := ac.AnnotateCtxWithSpan(context.Background(), "s")
	tracing.StartRecording(sp, tracing.SingleNodeRecording)
	Event(ctx, "a")
	sp.Finish()
	if err := tracing.TestingCheckRecordedSpans(tracing.GetRecording(sp), `
	  span s:
			event: [ambient] a
	`); err != nil {
		t.Fatal(err)
	}
}

func TestAnnotateCtxNodeStoreReplica(t *testing.T) {
	// Test the scenario of a context being continually re-annotated as it is
	// passed down a call stack.
	n := AmbientContext{}
	n.AddLogTagInt("n", 1)
	s := n
	s.AddLogTagInt("s", 2)
	r := s
	r.AddLogTagInt("r", 3)

	ctx := n.AnnotateCtx(context.Background())
	ctx = s.AnnotateCtx(ctx)
	ctx = r.AnnotateCtx(ctx)
	if exp, val := "[n1,s2,r3] test", MakeMessage(ctx, "test", nil); val != exp {
		t.Errorf("expected '%s', got '%s'", exp, val)
	}
	if bottom := contextBottomTag(ctx); bottom != r.tags {
		t.Errorf("expected %p, got %p", r.tags, bottom)
	}
}

func TestResetAndAnnotateCtx(t *testing.T) {
	ac := AmbientContext{}
	ac.AddLogTagInt("a", 1)

	ctx := context.Background()
	ctx = WithLogTag(ctx, "b", 2)
	ctx = ac.ResetAndAnnotateCtx(ctx)
	if exp, val := "[a1] test", MakeMessage(ctx, "test", nil); val != exp {
		t.Errorf("expected '%s', got '%s'", exp, val)
	}
}
