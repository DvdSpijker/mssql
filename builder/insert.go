package builder

import (
	"github.com/go-rel/rel"
	"github.com/go-rel/sql/builder"
)

// Insert builder.
type Insert struct {
	BufferFactory builder.BufferFactory
}

// Build sql query and its arguments.
func (i Insert) Build(table string, primaryField string, mutates map[string]rel.Mutate) (string, []interface{}) {
	var (
		buffer            = i.BufferFactory.Create()
		_, identityInsert = mutates[primaryField]
		arguments         = make([]interface{}, 0, len(mutates))
	)

	if identityInsert {
		buffer.WriteString("IF OBJECTPROPERTY(OBJECT_ID('")
		buffer.WriteEscape(table)
		buffer.WriteString("'), 'TableHasIdentity') = 1 ")
		buffer.WriteString("SET IDENTITY_INSERT ")
		buffer.WriteEscape(table)
		buffer.WriteString(" ON; ")
	}

	buffer.WriteString("INSERT INTO ")
	buffer.WriteEscape(table)
	buffer.WriteString(" (")

	index := 0
	for field, mut := range mutates {
		if mut.Type == rel.ChangeSetOp {
			if index > 0 {
				buffer.WriteByte(',')
			}

			buffer.WriteEscape(field)
			arguments = append(arguments, mut.Value)
			index++
		}
	}

	buffer.WriteString(")")

	if primaryField != "" {
		buffer.WriteString(" OUTPUT [INSERTED].")
		buffer.WriteEscape(primaryField)
	}

	buffer.WriteString(" VALUES (")

	for index := range arguments {
		if index > 0 {
			buffer.WriteByte(',')
		}

		buffer.WriteValue(arguments[index])
	}

	buffer.WriteString(");")

	if identityInsert {
		buffer.WriteString("IF OBJECTPROPERTY(OBJECT_ID('")
		buffer.WriteEscape(table)
		buffer.WriteString("'), 'TableHasIdentity') = 1 ")
		buffer.WriteString(" SET IDENTITY_INSERT ")
		buffer.WriteEscape(table)
		buffer.WriteString(" OFF; ")
	}

	return buffer.String(), buffer.Arguments()
}
