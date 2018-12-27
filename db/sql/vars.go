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

package sql

const (
    // PgServerVersion is the latest version of postgres that we claim to support.
    PgServerVersion = "9.5.0"
    // PgServerVersionNum is the latest version of postgres that we claim to support in the numeric format of "server_version_num".
    PgServerVersionNum = "90500"
)

type NTXIsoLevel string

// Transaction isolation levels
const (
    NTXSerializable    = NTXIsoLevel("serializable")
    NTXRepeatableRead  = NTXIsoLevel("repeatable read")
    NTXReadCommitted   = NTXIsoLevel("read committed")
    NTXReadUncommitted = NTXIsoLevel("read uncommitted")
)

type NTXAccessMode string

// Transaction access modes
const (
    ReadWrite = NTXAccessMode("read write")
    ReadOnly  = NTXAccessMode("read only")
)

type NTXStatus int

const (
    NTXNoTransaction   = NTXStatus(0)
    NTXInProgress      = NTXStatus(2)
    NTXPreparedSuccess = NTXStatus(3)
)
