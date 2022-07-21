package values

import (
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/liamcervante/terraform-provider-fakelocal/internal/types"
)

var _ tftypes.ValueConverter = Object{}
var _ tftypes.ValueCreator = Object{}

type Object struct {
	Types  map[string]types.Type `json:"types"`
	Values map[string]Value      `json:"values"`
}

func (o Object) ToTerraform5Value() (interface{}, error) {
	children := make(map[string]interface{})
	for name, value := range o.Values {
		child, err := value.ToTerraform5Value()
		if err != nil {
			return nil, err
		}
		children[name] = child
	}
	return children, nil
}

func (o Object) FromTerraform5Value(value tftypes.Value) error {
	var children map[string]tftypes.Value
	if err := value.As(&children); err != nil {
		return err
	}

	for name, child := range children {
		parsed, err := ValueForType(o.Types[name])
		if err != nil {
			return err
		}

		if err := parsed.FromTerraform5Value(child); err != nil {
			return err
		}
		o.Values[name] = parsed
	}
	return nil
}
