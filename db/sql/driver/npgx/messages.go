package npgx

import (
	"github.com/Ready-Stock/Noah/db/sql/pgwire/pgproto"
	"github.com/Ready-Stock/Noah/db/sql/types"
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