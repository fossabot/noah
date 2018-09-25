package npgx

import (
	"github.com/Ready-Stock/Noah/db/sql/pgio"
	"github.com/Ready-Stock/Noah/db/sql/pgwire/pgproto"
	"github.com/Ready-Stock/Noah/db/sql/types"
	"github.com/Ready-Stock/pgx/pgtype"
)



type PgError struct {
	Severity         string
	Code             string
	Message          string
	Detail           string
	Hint             string
	Position         int32
	InternalPosition int32
	InternalQuery    string
	Where            string
	SchemaName       string
	TableName        string
	ColumnName       string
	DataTypeName     string
	ConstraintName   string
	File             string
	Line             int32
	Routine          string
}

func (pe PgError) Error() string {
	return pe.Severity + ": " + pe.Message + " (SQLSTATE " + pe.Code + ")"
}

func (c *Conn) rxErrorResponse(msg *pgproto.ErrorResponse) PgError {
	err := PgError{
		Severity:         msg.Severity,
		Code:             msg.Code,
		Message:          msg.Message,
		Detail:           msg.Detail,
		Hint:             msg.Hint,
		Position:         msg.Position,
		InternalPosition: msg.InternalPosition,
		InternalQuery:    msg.InternalQuery,
		Where:            msg.Where,
		SchemaName:       msg.SchemaName,
		TableName:        msg.TableName,
		ColumnName:       msg.ColumnName,
		DataTypeName:     msg.DataTypeName,
		ConstraintName:   msg.ConstraintName,
		File:             msg.File,
		Line:             msg.Line,
		Routine:          msg.Routine,
	}

	if err.Severity == "FATAL" {
		c.die(err)
	}

	return err
}

func (c *Conn) rxBackendKeyData(msg *pgproto.BackendKeyData) {
	c.pid = msg.ProcessID
	c.secretKey = msg.SecretKey
}

func (c *Conn) rxReadyForQuery(msg *pgproto.ReadyForQuery) {
	c.pendingReadyForQueryCount--
	c.txStatus = msg.TxStatus
}

func (c *Conn) rxRowDescription(msg *pgproto.RowDescription) []FieldDescription {
	fields := make([]FieldDescription, len(msg.Fields))
	for i := 0; i < len(fields); i++ {
		fields[i].Name = msg.Fields[i].Name
		fields[i].Table = types.OID(msg.Fields[i].TableOID)
		fields[i].AttributeNumber = msg.Fields[i].TableAttributeNumber
		fields[i].DataType = types.OID(msg.Fields[i].DataTypeOID)
		fields[i].DataTypeSize = msg.Fields[i].DataTypeSize
		fields[i].Modifier = msg.Fields[i].TypeModifier
		fields[i].FormatCode = msg.Fields[i].Format
	}
	return fields
}


// appendParse appends a PostgreSQL wire protocol parse message to buf and returns it.
func appendParse(buf []byte, name string, query string, parameterOIDs []pgtype.OID) []byte {
	buf = append(buf, 'P')
	sp := len(buf)
	buf = pgio.AppendInt32(buf, -1)
	buf = append(buf, name...)
	buf = append(buf, 0)
	buf = append(buf, query...)
	buf = append(buf, 0)

	buf = pgio.AppendInt16(buf, int16(len(parameterOIDs)))
	for _, oid := range parameterOIDs {
		buf = pgio.AppendUint32(buf, uint32(oid))
	}
	pgio.SetInt32(buf[sp:], int32(len(buf[sp:])))

	return buf
}

// appendDescribe appends a PostgreSQL wire protocol describe message to buf and returns it.
func appendDescribe(buf []byte, objectType byte, name string) []byte {
	buf = append(buf, 'D')
	sp := len(buf)
	buf = pgio.AppendInt32(buf, -1)
	buf = append(buf, objectType)
	buf = append(buf, name...)
	buf = append(buf, 0)
	pgio.SetInt32(buf[sp:], int32(len(buf[sp:])))

	return buf
}

// appendSync appends a PostgreSQL wire protocol sync message to buf and returns it.
func appendSync(buf []byte) []byte {
	buf = append(buf, 'S')
	buf = pgio.AppendInt32(buf, 4)

	return buf
}

// appendBind appends a PostgreSQL wire protocol bind message to buf and returns it.
func appendBind(
	buf []byte,
	destinationPortal,
	preparedStatement string,
	connInfo *pgtype.ConnInfo,
	parameterOIDs []pgtype.OID,
	arguments []interface{},
	resultFormatCodes []int16,
) ([]byte, error) {
	buf = append(buf, 'B')
	sp := len(buf)
	buf = pgio.AppendInt32(buf, -1)
	buf = append(buf, destinationPortal...)
	buf = append(buf, 0)
	buf = append(buf, preparedStatement...)
	buf = append(buf, 0)

	buf = pgio.AppendInt16(buf, int16(len(parameterOIDs)))
	for i, oid := range parameterOIDs {
		buf = pgio.AppendInt16(buf, chooseParameterFormatCode(connInfo, oid, arguments[i]))
	}

	buf = pgio.AppendInt16(buf, int16(len(arguments)))
	for i, oid := range parameterOIDs {
		var err error
		buf, err = encodePreparedStatementArgument(connInfo, buf, oid, arguments[i])
		if err != nil {
			return nil, err
		}
	}

	buf = pgio.AppendInt16(buf, int16(len(resultFormatCodes)))
	for _, fc := range resultFormatCodes {
		buf = pgio.AppendInt16(buf, fc)
	}
	pgio.SetInt32(buf[sp:], int32(len(buf[sp:])))

	return buf, nil
}

// appendExecute appends a PostgreSQL wire protocol execute message to buf and returns it.
func appendExecute(buf []byte, portal string, maxRows uint32) []byte {
	buf = append(buf, 'E')
	sp := len(buf)
	buf = pgio.AppendInt32(buf, -1)

	buf = append(buf, portal...)
	buf = append(buf, 0)
	buf = pgio.AppendUint32(buf, maxRows)

	pgio.SetInt32(buf[sp:], int32(len(buf[sp:])))

	return buf
}

// appendQuery appends a PostgreSQL wire protocol query message to buf and returns it.
func appendQuery(buf []byte, query string) []byte {
	buf = append(buf, 'Q')
	sp := len(buf)
	buf = pgio.AppendInt32(buf, -1)
	buf = append(buf, query...)
	buf = append(buf, 0)
	pgio.SetInt32(buf[sp:], int32(len(buf[sp:])))

	return buf
}