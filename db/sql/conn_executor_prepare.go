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

import (
    "fmt"
    "github.com/kataras/go-errors"
    "github.com/readystock/noah/db/sql/pgwire/pgerror"
    "github.com/readystock/noah/db/sql/pgwire/pgwirebase"
    "github.com/readystock/noah/db/sql/plan"
    "github.com/readystock/noah/db/sql/types"
    nodes "github.com/readystock/pg_query_go/nodes"
    "strconv"
)

func (ex *connExecutor) execPrepare(parseCmd PrepareStmt) error {
    if parseCmd.Name != "" {
        if _, ok := ex.prepStmtsNamespace.prepStmts[parseCmd.Name]; ok {
            err := pgerror.NewErrorf(
                pgerror.CodeDuplicatePreparedStatementError,
                "prepared statement %q already exists", parseCmd.Name,
            )
            return err
        }
    } else {
        // Deallocate the unnamed statement, if it exists.
        ex.deletePreparedStmt("")
    }

    ps, err := ex.addPreparedStmt(parseCmd.Name, parseCmd.PGQuery, parseCmd.TypeHints)
    if err != nil {
        return err
    }

    // Convert the inferred SQL types back to an array of pgwire Oids.
    inTypes := make([]types.OID, 0, len(ps.TypeHints))
    if len(ps.TypeHints) > pgwirebase.MaxPreparedStatementArgs {
        return pgwirebase.NewProtocolViolationErrorf(
            "more than %d arguments to prepared statement: %d",
            pgwirebase.MaxPreparedStatementArgs, len(ps.TypeHints))
    }
    for k, t := range ps.TypeHints {
        i, err := strconv.Atoi(k)
        if err != nil || i < 1 {
            return pgerror.NewErrorf(
                pgerror.CodeUndefinedParameterError, "invalid placeholder name: $%s", k)
        }
        // Placeholder names are 1-indexed; the arrays in the protocol are
        // 0-indexed.
        i--
        // Grow inTypes to be at least as large as i. Prepopulate all
        // slots with the hints provided, if any.
        for j := len(inTypes); j <= i; j++ {
            inTypes = append(inTypes, 0)
            if j < len(parseCmd.RawTypeHints) {
                inTypes[j] = parseCmd.RawTypeHints[j]
            }
        }
        // OID to Datum is not a 1-1 mapping (for example, int4 and int8
        // both map to TypeInt), so we need to maintain the types sent by
        // the client.
        if inTypes[i] != 0 {
            continue
        }
        inTypes[i] = t.GetOID()
    }
    for i, t := range inTypes {
        if t == 0 {
            return pgerror.NewErrorf(
                pgerror.CodeIndeterminateDatatypeError,
                "could not determine data type of placeholder $%d", i+1)
        }
    }
    // Remember the inferred placeholder types so they can be reported on
    // Describe.
    ps.InTypes = inTypes

    return nil
}

func (ex *connExecutor) execBind(bindCmd BindStmt) error {
    portalName := bindCmd.PortalName
    // The unnamed portal can be freely overwritten.
    if portalName != "" {
        if _, ok := ex.prepStmtsNamespace.portals[portalName]; ok {
            return pgerror.NewErrorf(pgerror.CodeDuplicateCursorError, "portal %q already exists", portalName)
        }
    } else {
        // Deallocate the unnamed portal, if it exists.
        ex.deletePortal("")
    }

    ps, ok := ex.prepStmtsNamespace.prepStmts[bindCmd.PreparedStatementName]
    if !ok {
        return pgerror.NewErrorf(pgerror.CodeInvalidSQLStatementNameError, "unknown prepared statement %q", bindCmd.PreparedStatementName)
    }

    numQArgs := uint16(len(ps.InTypes))
    qArgFormatCodes := bindCmd.ArgFormatCodes

    // If a single code is specified, it is applied to all arguments.
    if len(qArgFormatCodes) != 1 && len(qArgFormatCodes) != int(numQArgs) {
        return pgwirebase.NewProtocolViolationErrorf("wrong number of format codes specified: %d for %d arguments", len(qArgFormatCodes), numQArgs)
    }
    // If a single format code was specified, it applies to all the arguments.
    if len(qArgFormatCodes) == 1 {
        fmtCode := qArgFormatCodes[0]
        qArgFormatCodes = make([]pgwirebase.FormatCode, numQArgs)
        for i := range qArgFormatCodes {
            qArgFormatCodes[i] = fmtCode
        }
    }

    if len(bindCmd.Args) != int(numQArgs) {
        return pgwirebase.NewProtocolViolationErrorf("expected %d arguments, got %d", numQArgs, len(bindCmd.Args))
    }

    qargs := plan.QueryArguments{}
    for i, arg := range bindCmd.Args {
        k := strconv.Itoa(i + 1)
        t := ps.InTypes[i]
        if arg == nil {
            // nil indicates a NULL argument value.
            qargs[k] = &types.Unknown{
                Status: types.Null,
            }
        } else {
            dt, ok := ex.typeInfo.DataTypeForOID(t)
            if !ok {
                return errors.New("invalid thing")
            }
            switch qArgFormatCodes[i] {
            case pgwirebase.FormatBinary:
                binaryDecode, ok := dt.Value.(types.BinaryDecoder)
                if !ok {
                    return errors.New("type is not a text decoder")
                }
                if err := binaryDecode.DecodeBinary(ex.typeInfo, arg); err != nil {
                    return err
                }
                qargs[k] = binaryDecode.(types.Value)
            case pgwirebase.FormatText:
                textDecode, ok := dt.Value.(types.TextDecoder)
                if !ok {
                    return errors.New("type is not a text decoder")
                }
                if err := textDecode.DecodeText(ex.typeInfo, arg); err != nil {
                    return err
                }
                qargs[k] = textDecode.(types.Value)
            default:
                return errors.New("invalid format code")
            }
            // d, err := pgwirebase.DecodeOidDatum(t, qArgFormatCodes[i], arg)
            // if err != nil {
            //     if _, ok := err.(*pgerror.Error); ok {
            //         return err
            //     }
            //     return pgwirebase.NewProtocolViolationErrorf(
            //         "error in argument for $%d: %s", i+1, err.Error())
            //
            // }
            // qargs[k] = d
        }
    }

    // numCols := len(ps.Columns)
    // if (len(bindCmd.OutFormats) > 1) && (len(bindCmd.OutFormats) != numCols) {
    //     return pgwirebase.NewProtocolViolationErrorf(
    //         "expected 1 or %d for number of format codes, got %d",
    //         numCols, len(bindCmd.OutFormats))
    // }
    //
    // columnFormatCodes := bindCmd.OutFormats
    // if len(bindCmd.OutFormats) == 1 {
    //     // Apply the format code to every column.
    //     columnFormatCodes = make([]pgwirebase.FormatCode, numCols)
    //     for i := 0; i < numCols; i++ {
    //         columnFormatCodes[i] = bindCmd.OutFormats[0]
    //     }
    // }

    // Create the new PreparedPortal.
    if err := ex.addPortal(portalName, bindCmd.PreparedStatementName, ps.PreparedStatement, qargs, qArgFormatCodes); err != nil {
        return err
    }
    return nil
}

func (ex *connExecutor) addPortal(
    portalName string,
    psName string,
    stmt *PreparedStatement,
    qargs plan.QueryArguments,
    outFormats []pgwirebase.FormatCode,
) error {
    if _, ok := ex.prepStmtsNamespace.portals[portalName]; ok {
        panic(fmt.Sprintf("portal already exists: %q", portalName))
    }
    portal := ex.newPreparedPortal(stmt, qargs, outFormats)
    ex.prepStmtsNamespace.portals[portalName] = portalEntry{
        PreparedPortal: &portal,
        psName:         psName,
    }
    ex.prepStmtsNamespace.prepStmts[psName].portals[portalName] = struct{}{}
    return nil
}

func (ex *connExecutor) addPreparedStmt(name string, stmt nodes.Stmt, parseTypeHints plan.PlaceholderTypes) (*PreparedStatement, error) {
    if _, ok := ex.prepStmtsNamespace.prepStmts[name]; ok {
        panic(fmt.Sprintf("prepared statement already exists: %q", name))
    }
    // Prepare the query. This completes the typing of placeholders.
    prepared, err := ex.prepare(stmt, parseTypeHints)
    if err != nil {
        return nil, err
    }
    ex.prepStmtsNamespace.prepStmts[name] = prepStmtEntry{
        PreparedStatement: prepared,
        portals:           make(map[string]struct{}),
    }
    return prepared, nil
}

func (ex *connExecutor) prepare(stmt nodes.Stmt, parseTypeHints plan.PlaceholderTypes) (*PreparedStatement, error) {
    prepared := &PreparedStatement{
        TypeHints: parseTypeHints,
        Statement: &stmt,
    }

    if stmt == nil {
        return prepared, nil
    }

    // prepared.Columns = getPlanColumns(stmt)

    return prepared, nil
}

func (ex *connExecutor) deletePreparedStmt(name string) {
    psEntry, ok := ex.prepStmtsNamespace.prepStmts[name]
    if !ok {
        return
    }
    for portalName := range psEntry.portals {
        ex.deletePortal(portalName)
    }
    delete(ex.prepStmtsNamespace.prepStmts, name)
}

func (ex *connExecutor) deletePortal(name string) {
    portalEntry, ok := ex.prepStmtsNamespace.portals[name]
    if !ok {
        return
    }
    delete(ex.prepStmtsNamespace.portals, name)
    delete(ex.prepStmtsNamespace.prepStmts[portalEntry.psName].portals, name)
}
