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

// This code was derived from https://github.com/youtube/vitess.

package tree

// Truncate represents a TRUNCATE statement.
type Truncate struct {
    Tables       NormalizableTableNames
    DropBehavior DropBehavior
}

// Format implements the NodeFormatter interface.
func (node *Truncate) Format(ctx *FmtCtx) {
    ctx.WriteString("TRUNCATE TABLE ")
    sep := ""
    for i := range node.Tables {
        ctx.WriteString(sep)
        ctx.FormatNode(&node.Tables[i])
        sep = ", "
    }
    if node.DropBehavior != DropDefault {
        ctx.WriteByte(' ')
        ctx.WriteString(node.DropBehavior.String())
    }
}
