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

package pgproto

import (
    "bytes"
    "encoding/json"

    "github.com/readystock/pgx/pgio"
)

type Query struct {
    String string
}

func (*Query) Frontend() {}

func (dst *Query) Decode(src []byte) error {
    i := bytes.IndexByte(src, 0)
    if i != len(src)-1 {
        return &invalidMessageFormatErr{messageType: "Query"}
    }

    dst.String = string(src[:i])

    return nil
}

func (src *Query) Encode(dst []byte) []byte {
    dst = append(dst, 'Q')
    dst = pgio.AppendInt32(dst, int32(4+len(src.String)+1))

    dst = append(dst, src.String...)
    dst = append(dst, 0)

    return dst
}

func (src *Query) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct {
        Type   string
        String string
    }{
        Type:   "Query",
        String: src.String,
    })
}
