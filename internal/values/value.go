package values

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/liamcervante/terraform-provider-fakelocal/internal/types"
	"reflect"
)

var _ tftypes.ValueConverter = Value{}
var _ tftypes.ValueCreator = Value{}

type Value struct {
	Type types.Type `json:"type"`

	Boolean bool    `json:"boolean"`
	Float   float64 `json:"float"`
	Integer int64   `json:"integer"`
	Number  float64 `json:"number"`
	String  string  `json:"string"`

	List   List   `json:"list"`
	Map    Map    `json:"map"`
	Object Object `json:"object"`
	Set    Set    `json:"set"`
}

func ValueForType(t types.Type) (Value, error) {
	switch target := t.(type) {
	case types.PrimitiveType:
		return Value{Type: t}, nil
	case types.ListType:
		return Value{
			Type: target,
			List: List{
				Type:   target.Type,
				Values: []Value{},
			},
		}, nil
	case types.MapType:
		return Value{
			Type: target,
			Map: Map{
				Type:   target.Type,
				Values: map[string]Value{},
			},
		}, nil
	case types.SetType:
		return Value{
			Type: target,
			Set: Set{
				Type:   target.Type,
				Values: []Value{},
			},
		}, nil
	case types.ObjectType:
		return Value{
			Type: target,
			Object: Object{
				Types:  target.Types,
				Values: map[string]Value{},
			},
		}, nil
	default:
		return Value{}, errors.New("unrecognized type: " + reflect.TypeOf(target).String())
	}
}

func (v Value) ToTerraform5Value() (interface{}, error) {
	switch v.Type.String() {
	case types.Boolean:
		return tftypes.NewValue(tftypes.Bool, v.Boolean), nil
	case types.Float:
		return tftypes.NewValue(tftypes.Number, v.Float), nil
	case types.Integer:
		return tftypes.NewValue(tftypes.Number, v.Integer), nil
	case types.Number:
		return tftypes.NewValue(tftypes.Number, v.Number), nil
	case types.String:
		return tftypes.NewValue(tftypes.String, v.String), nil
	case types.List:
		return v.List.ToTerraform5Value()
	case types.Map:
		return v.Map.ToTerraform5Value()
	case types.Set:
		return v.Set.ToTerraform5Value()
	case types.Object:
		return v.Object.ToTerraform5Value()
	default:
		return tfsdk.Attribute{}, errors.New("unrecognized type: " + v.Type.String())
	}
}

func (v Value) FromTerraform5Value(value tftypes.Value) error {
	switch v.Type.String() {
	case types.Boolean:
		return value.As(&v.Boolean)
	case types.Float:
		return value.As(&v.Float)
	case types.Integer:
		return value.As(&v.Integer)
	case types.Number:
		return value.As(&v.Boolean)
	case types.String:
		return value.As(&v.String)
	case types.List:
		return v.List.FromTerraform5Value(value)
	case types.Map:
		return v.Map.FromTerraform5Value(value)
	case types.Set:
		return v.Set.FromTerraform5Value(value)
	case types.Object:
		return v.Object.FromTerraform5Value(value)
	default:
		return errors.New("unrecognized type: " + v.Type.String())
	}
}
