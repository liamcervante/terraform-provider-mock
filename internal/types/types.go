package types

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

const (
	Boolean string = "boolean"
	Float   string = "float"
	Integer string = "integer"
	Number  string = "number"
	String  string = "string"

	List   string = "list"
	Map    string = "map"
	Object string = "object"
	Set    string = "set"
)

type Type struct {
	Type string `tfsdk:"type" json:"type"`

	ElementType *Type            `tfsdk:"list_type" json:"list_type,omitempty"`
	ObjectType  map[string]*Type `tfsdk:"object_type" json:"object_type,omitempty"`
}

func (t *Type) ToTerraform5Type() (tftypes.Type, error) {
	switch t.Type {
	case Boolean:
		return tftypes.Bool, nil
	case String:
		return tftypes.String, nil
	case Float, Integer, Number:
		return tftypes.Number, nil
	case List:
		tfType, err := t.ElementType.ToTerraform5Type()
		if err != nil {
			return nil, err
		}
		return tftypes.List{
			ElementType: tfType,
		}, nil
	case Map:
		tfType, err := t.ElementType.ToTerraform5Type()
		if err != nil {
			return nil, err
		}
		return tftypes.Map{
			ElementType: tfType,
		}, nil
	case Object:
		tfTypes := make(map[string]tftypes.Type)
		for name, child := range t.ObjectType {
			tfType, err := child.ToTerraform5Type()
			if err != nil {
				return nil, err
			}
			tfTypes[name] = tfType
		}
		return tftypes.Object{
			AttributeTypes: tfTypes,
		}, nil

	case Set:
		tfType, err := t.ElementType.ToTerraform5Type()
		if err != nil {
			return nil, err
		}
		return tftypes.Set{
			ElementType: tfType,
		}, nil
	default:
		return nil, errors.New("Unrecognized type " + t.Type)
	}
}

func FromTerraform5Type(tft tftypes.Type) (*Type, error) {
	switch {
	case tft.Is(tftypes.Bool):
		return &Type{Type: Boolean}, nil
	case tft.Is(tftypes.String):
		return &Type{Type: String}, nil
	case tft.Is(tftypes.Number):
		return &Type{Type: Number}, nil
	case tft.Is(tftypes.List{}):
		elementType, err := FromTerraform5Type(tft.(tftypes.List).ElementType)
		if err != nil {
			return nil, err
		}
		return &Type{
			Type:        List,
			ElementType: elementType,
		}, nil
	case tft.Is(tftypes.Map{}):
		elementType, err := FromTerraform5Type(tft.(tftypes.Map).ElementType)
		if err != nil {
			return nil, err
		}
		return &Type{
			Type:        Map,
			ElementType: elementType,
		}, nil
	case tft.Is(tftypes.Object{}):

		objectTypes := make(map[string]*Type)
		for name, child := range tft.(tftypes.Object).AttributeTypes {
			typ, err := FromTerraform5Type(child)
			if err != nil {
				return nil, err
			}
			objectTypes[name] = typ
		}
		return &Type{
			Type:       Object,
			ObjectType: objectTypes,
		}, nil
	case tft.Is(tftypes.Set{}):
		elementType, err := FromTerraform5Type(tft.(tftypes.Set).ElementType)
		if err != nil {
			return nil, err
		}
		return &Type{
			Type:        Set,
			ElementType: elementType,
		}, nil
	default:
		return nil, errors.New("Unrecognized type " + tft.String())
	}
}
