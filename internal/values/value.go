package values

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"reflect"
)

var _ tftypes.ValueConverter = Value{}
var _ tftypes.ValueCreator = Value{}

type Value struct {
	Type attr.Type `tfsdk:"type" json:"type"`

	Boolean bool    `tfsdk:"boolean" json:"boolean"`
	Float   float64 `tfsdk:"boolean" json:"float"`
	Integer int64   `tfsdk:"integer" json:"integer"`
	Number  float64 `tfsdk:"number" json:"number"`
	String  string  `tfsdk:"string" json:"string"`

	List   List   `tfsdk:"list" json:"list"`
	Map    Map    `tfsdk:"map" json:"map"`
	Object Object `tfsdk:"object" json:"object"`
	Set    Set    `tfsdk:"set" json:"set"`
}

func ValueForType(t attr.Type) Value {
	switch target := t.(type) {
	case types.ListType:
		return Value{
			Type: target,
			List: List{
				Type:   target.ElemType,
				Values: []Value{},
			},
		}
	case types.MapType:
		return Value{
			Type: target,
			Map: Map{
				Type:   target.ElemType,
				Values: map[string]Value{},
			},
		}
	case types.SetType:
		return Value{
			Type: target,
			Set: Set{
				Type:   target.ElemType,
				Values: []Value{},
			},
		}
	case types.ObjectType:
		return Value{
			Type: target,
			Object: Object{
				Types:  target.AttrTypes,
				Values: map[string]Value{},
			},
		}
	}

	// Otherwise, just apply the type without extra special consideration.
	return Value{Type: t}
}

func (v Value) SetId(id string) error {
	switch target := v.Type.(type) {
	case types.ObjectType:
		v.Object.Values["id"] = Value{String: id}
		v.Object.Types["id"] = types.StringType
		return nil
	default:
		return errors.New("can only set the ID on an Object type as found: " + reflect.TypeOf(target).String())
	}
}

func (v Value) GetId() (string, error) {
	switch target := v.Type.(type) {
	case types.ObjectType:
		return v.Object.Values["id"].String, nil
	default:
		return "", errors.New("can only retrieve the ID on an Object type as found: " + reflect.TypeOf(target).String())
	}
}

func (v Value) ToTerraform5Value() (interface{}, error) {
	switch reflect.TypeOf(v.Type) {
	case reflect.TypeOf(types.BoolType):
		return tftypes.NewValue(tftypes.Bool, v.Boolean), nil
	case reflect.TypeOf(types.Float64Type):
		return tftypes.NewValue(tftypes.Number, v.Float), nil
	case reflect.TypeOf(types.Int64Type):
		return tftypes.NewValue(tftypes.Number, v.Integer), nil
	case reflect.TypeOf(types.NumberType):
		return tftypes.NewValue(tftypes.Number, v.Number), nil
	case reflect.TypeOf(types.StringType):
		return tftypes.NewValue(tftypes.String, v.String), nil
	case reflect.TypeOf(types.ListType{}):
		return v.List.ToTerraform5Value()
	case reflect.TypeOf(types.MapType{}):
		return v.Map.ToTerraform5Value()
	case reflect.TypeOf(types.SetType{}):
		return v.Set.ToTerraform5Value()
	case reflect.TypeOf(types.ObjectType{}):
		return v.Object.ToTerraform5Value()
	default:
		return tfsdk.Attribute{}, errors.New("unrecognized type: " + v.Type.String())
	}
}

func (v Value) FromTerraform5Value(value tftypes.Value) error {
	value.Type()
	switch reflect.TypeOf(v.Type) {
	case reflect.TypeOf(types.BoolType):
		return value.As(&v.Boolean)
	case reflect.TypeOf(types.Float64Type):
		return value.As(&v.Float)
	case reflect.TypeOf(types.Int64Type):
		return value.As(&v.Integer)
	case reflect.TypeOf(types.NumberType):
		return value.As(&v.Number)
	case reflect.TypeOf(types.StringType):
		return value.As(&v.String)
	case reflect.TypeOf(types.ListType{}):
		return v.List.FromTerraform5Value(value)
	case reflect.TypeOf(types.MapType{}):
		return v.Map.FromTerraform5Value(value)
	case reflect.TypeOf(types.SetType{}):
		return v.Set.FromTerraform5Value(value)
	case reflect.TypeOf(types.ObjectType{}):
		return v.Object.FromTerraform5Value(value)
	default:
		return errors.New("unrecognized type: " + v.Type.String())
	}
}
