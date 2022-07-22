package types

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
