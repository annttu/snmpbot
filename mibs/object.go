package mibs

import (
	"fmt"
	"github.com/qmsk/snmpbot/snmp"
)

type Object struct {
	*ID
	Table *Table

	Syntax
}

func (object *Object) Unpack(varBind snmp.VarBind) (interface{}, error) {
	return object.Syntax.Unpack(varBind)
}

func (object *Object) Format(varBind snmp.VarBind) (string, interface{}, error) {
	name := object.FormatOID(varBind.OID())
	value, err := object.Unpack(varBind)

	return name, value, err
}

func (object *Object) FormatIndex(oid snmp.OID) string {
	if object.Table == nil {
		return object.FormatOID(oid)
	}

	if index := object.OID.Index(oid); index == nil {
		return oid.String()
	} else if len(index) == 0 {
		return object.String()
	} else if indexString, err := object.Table.FormatIndex(index); err != nil {
		return fmt.Sprintf("%s::%s%s", object.MIB.Name, object.Name, snmp.OID(index).String())
	} else {
		return fmt.Sprintf("%s::%s%s", object.MIB.Name, object.Name, indexString)
	}
}