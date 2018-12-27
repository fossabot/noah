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

package types

type Decimal Numeric

func (dst *Decimal) Set(src interface{}) error {
    return (*Numeric)(dst).Set(src)
}

func (dst *Decimal) Get() interface{} {
    return (*Numeric)(dst).Get()
}

func (src *Decimal) AssignTo(dst interface{}) error {
    return (*Numeric)(src).AssignTo(dst)
}

func (dst *Decimal) DecodeText(ci *ConnInfo, src []byte) error {
    return (*Numeric)(dst).DecodeText(ci, src)
}

func (dst *Decimal) DecodeBinary(ci *ConnInfo, src []byte) error {
    return (*Numeric)(dst).DecodeBinary(ci, src)
}

func (src *Decimal) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
    return (*Numeric)(src).EncodeText(ci, buf)
}

func (src *Decimal) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
    return (*Numeric)(src).EncodeBinary(ci, buf)
}
