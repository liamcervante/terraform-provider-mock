package types

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	Boolean string = "boolean"
	Float   string = "float"
	Integer string = "integer"
	Number  string = "number"
	String  string = "string"

	List   string = "list"
	Map    string = "map"
	Object string = "type"
	Set    string = "set"
)

func Test(interface{}) Type {

}

func GetType(t attr.Type) Type {

	attribute := tfsdk.Attribute{
		Type: types.BoolType,
	}

	b := types.BoolType

	t.(types.BoolType)

	switch target := t.(type) {
	case attr.Type(types.BoolType):
		return PrimitiveType{
			Type: Boolean,
		}
	case types.MapType:
		return MapType{
			Type: GetType(target.ElemType),
		}
	case types.ObjectType
	}
}

type Type interface {
	String() string
}

type PrimitiveType struct {
	Type string
}

func (p PrimitiveType) String() string {
	return List
}

type ListType struct {
	Type Type
}

func (l ListType) String() string {
	return List
}

type MapType struct {
	Type Type
}

func (m MapType) String() string {
	return Map
}

type ObjectType struct {
	Types map[string]Type
}

func (o ObjectType) String() string {
	return Object
}

type SetType struct {
	Type Type
}

func (s SetType) String() string {
	return Set
}
