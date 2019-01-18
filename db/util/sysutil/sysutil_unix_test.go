// +build !windows

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
 */

package sysutil

import (
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func TestRefreshSignaledChan(t *testing.T) {
	ch := RefreshSignaledChan()

	if err := unix.Kill(unix.Getpid(), refreshSignal); err != nil {
		t.Error(err)
		return
	}

	select {
	case sig := <-ch:
		if refreshSignal != sig {
			t.Fatalf("expected signal %s, but got %s", refreshSignal, sig)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out while waiting for refresh signal")
	}
}
