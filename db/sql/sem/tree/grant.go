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

// Grant represents a GRANT statement.
type Grant struct {
    Targets  TargetList
    Grantees NameList
}

// TargetList represents a list of targets.
// Only one field may be non-nil.
type TargetList struct {
    Databases NameList
    Tables    TablePatterns
}

// Format implements the NodeFormatter interface.
func (tl *TargetList) Format(ctx *FmtCtx) {
    if tl.Databases != nil {
        ctx.WriteString("DATABASE ")
        ctx.FormatNode(&tl.Databases)
    } else {
        ctx.WriteString("TABLE ")
        ctx.FormatNode(&tl.Tables)
    }
}

// Format implements the NodeFormatter interface.
// func (node *Grant) Format(ctx *FmtCtx) {
//     ctx.WriteString("GRANT ")
//     node.Privileges.Format(ctx.Buffer)
//     ctx.WriteString(" ON ")
//     ctx.FormatNode(&node.Targets)
//     ctx.WriteString(" TO ")
//     ctx.FormatNode(&node.Grantees)
// }

// // GrantRole represents a GRANT <role> statement.
// type GrantRole struct {
//     Roles       NameList
//     Members     NameList
//     AdminOption bool
// }
//
// // Format implements the NodeFormatter interface.
// func (node *GrantRole) Format(ctx *FmtCtx) {
//     ctx.WriteString("GRANT ")
//     ctx.FormatNode(&node.Roles)
//     ctx.WriteString(" TO ")
//     ctx.FormatNode(&node.Members)
//     if node.AdminOption {
//         ctx.WriteString(" WITH ADMIN OPTION")
//     }
// }
