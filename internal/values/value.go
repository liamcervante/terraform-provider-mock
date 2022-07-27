package values

import (
	"errors"
	"math/big"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/liamcervante/terraform-provider-mock/internal/types"
)

type Value struct {
	Boolean *bool      `tfsdk:"boolean" json:"boolean,omitempty"`
	Number  *big.Float `tfsdk:"number" json:"number,omitempty"`
	String  *string    `tfsdk:"string" json:"string,omitempty"`

	List   *[]Value          `tfsdk:"list" json:"list,omitempty"`
	Map    *map[string]Value `tfsdk:"map" json:"map,omitempty"`
	Object *map[string]Value `tfsdk:"object" json:"object,omitempty"`
	Set    *[]Value          `tfsdk:"set" json:"set,omitempty"`
}

func toTerraform5Value(v Value, t *types.Type) (tftypes.Value, error) {
	switch t.Type {
	case types.Boolean:
		return tftypes.NewValue(tftypes.Bool, v.Boolean), nil
	case types.Number, types.Integer, types.Float:
		return tftypes.NewValue(tftypes.Number, v.Number), nil
	case types.String:
		return tftypes.NewValue(tftypes.String, v.String), nil
	case types.List:
		value, typ, err := listToTerraform5Value(v.List, t)
		if err != nil {
			return tftypes.Value{}, err
		}
		return tftypes.NewValue(typ, value), nil
	case types.Map:
		value, typ, err := mapToTerraform5Value(v.Map, t)
		if err != nil {
			return tftypes.Value{}, err
		}
		return tftypes.NewValue(typ, value), nil
	case types.Set:
		value, typ, err := setToTerraform5Value(v.Set, t)
		if err != nil {
			return tftypes.Value{}, err
		}
		return tftypes.NewValue(typ, value), nil
	case types.Object:
		value, typ, err := objectToTerraform5Value(v.Object, t)
		if err != nil {
			return tftypes.Value{}, err
		}
		return tftypes.NewValue(typ, value), nil
	default:
		return tftypes.Value{}, errors.New("unrecognized type: " + t.Type)
	}
}

func fromTerraform5Value(value tftypes.Value) (*Value, error) {
	if value.IsNull() {
		return nil, nil
	}

	typ := value.Type()
	switch {
	case typ.Is(tftypes.Bool):
		values := Value{}
		err := value.As(&values.Boolean)
		return &values, err
	case typ.Is(tftypes.Number):
		values := Value{}
		err := value.As(&values.Number)
		return &values, err
	case typ.Is(tftypes.String):
		values := Value{}
		err := value.As(&values.String)
		return &values, err
	case typ.Is(tftypes.List{}):
		return listFromTerraform5Value(value)
	case typ.Is(tftypes.Map{}):
		return mapFromTerraform5Value(value)
	case typ.Is(tftypes.Set{}):
		return setFromTerraform5Value(value)
	case typ.Is(tftypes.Object{}):
		return objectFromTerraform5Value(value)
	default:
		return nil, errors.New("unrecognized type: " + typ.String())
	}
}

func listToTerraform5Value(v *[]Value, t *types.Type) (interface{}, tftypes.Type, error) {
	tfType, err := t.ToTerraform5Type()
	if err != nil {
		return nil, tftypes.List{}, err
	}

	if v == nil {
		return nil, tfType, nil
	}

	var children []tftypes.Value
	for _, values := range *v {
		child, err := toTerraform5Value(values, t.ElementType)
		if err != nil {
			return children, tfType, err
		}
		children = append(children, child)
	}

	return children, tfType, nil
}

func listFromTerraform5Value(value tftypes.Value) (*Value, error) {
	values := Value{}

	var children []tftypes.Value
	if err := value.As(&children); err != nil {
		return nil, err
	}

	list := []Value{}
	for _, child := range children {
		parsed, err := fromTerraform5Value(child)
		if err != nil {
			return nil, err
		}
		list = append(list, *parsed)
	}

	values.List = &list
	return &values, nil
}

func mapToTerraform5Value(v *map[string]Value, t *types.Type) (interface{}, tftypes.Type, error) {
	tfType, err := t.ToTerraform5Type()
	if err != nil {
		return nil, tftypes.Map{}, err
	}

	if v == nil {
		return nil, tfType, nil
	}

	children := make(map[string]tftypes.Value)
	for name, values := range *v {
		child, err := toTerraform5Value(values, t.ElementType)
		if err != nil {
			return children, tfType, err
		}
		children[name] = child
	}

	return children, tfType, nil
}

func mapFromTerraform5Value(value tftypes.Value) (*Value, error) {
	values := Value{}

	var children map[string]tftypes.Value
	if err := value.As(&children); err != nil {
		return nil, err
	}

	mapValues := make(map[string]Value)
	for name, child := range children {
		parsed, err := fromTerraform5Value(child)
		if err != nil {
			return nil, err
		}
		mapValues[name] = *parsed
	}
	values.Map = &mapValues
	return &values, nil
}

func objectToTerraform5Value(v *map[string]Value, t *types.Type) (interface{}, tftypes.Type, error) {
	tfType, err := t.ToTerraform5Type()
	if err != nil {
		return nil, tftypes.Object{}, err
	}

	if v == nil {
		return nil, tfType, nil
	}

	tfValues := make(map[string]tftypes.Value)
	for name, typ := range t.ObjectType {
		var err error
		if value, ok := (*v)[name]; ok {
			if tfValues[name], err = toTerraform5Value(value, typ); err != nil {
				return tfValues, tftypes.Object{}, err
			}
			continue
		}

		if tfValues[name], err = toTerraform5Value(Value{}, typ); err != nil {
			return tfValues, tftypes.Object{}, err
		}
	}

	return tfValues, tfType, nil
}

func objectFromTerraform5Value(value tftypes.Value) (*Value, error) {
	values := Value{}
	var children map[string]tftypes.Value
	if err := value.As(&children); err != nil {
		return nil, err
	}

	objectValues := make(map[string]Value)
	for name, child := range children {
		parsed, err := fromTerraform5Value(child)
		if err != nil {
			return nil, err
		}
		if parsed != nil {
			objectValues[name] = *parsed
		}
	}
	values.Object = &objectValues
	return &values, nil
}

func setToTerraform5Value(v *[]Value, t *types.Type) (interface{}, tftypes.Type, error) {
	tfType, err := t.ToTerraform5Type()
	if err != nil {
		return nil, tftypes.Set{}, err
	}

	if v == nil {
		return nil, tfType, nil
	}

	var children []tftypes.Value
	for _, value := range *v {
		child, err := toTerraform5Value(value, t.ElementType)
		if err != nil {
			return children, tfType, err
		}
		children = append(children, child)
	}

	return children, tfType, nil
}

func setFromTerraform5Value(value tftypes.Value) (*Value, error) {
	values := Value{}
	var children []tftypes.Value
	if err := value.As(&children); err != nil {
		return nil, err
	}

	set := []Value{}
	for _, child := range children {
		parsed, err := fromTerraform5Value(child)
		if err != nil {
			return nil, err
		}
		set = append(set, *parsed)
	}
	values.Set = &set
	return &values, nil
}
