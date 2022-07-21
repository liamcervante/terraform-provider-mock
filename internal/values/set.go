package values

import (
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/liamcervante/terraform-provider-fakelocal/internal/types"
)

var _ tftypes.ValueConverter = Set{}
var _ tftypes.ValueCreator = Set{}

type Set struct {
	Type   types.Type `json:"type"`
	Values []Value    `json:"values"`
}

func (s Set) ToTerraform5Value() (interface{}, error) {
	var children []interface{}
	for _, value := range s.Values {
		child, err := value.ToTerraform5Value()
		if err != nil {
			return nil, err
		}
		children = append(children, child)
	}
	return children, nil
}

func (s Set) FromTerraform5Value(value tftypes.Value) error {
	var children []tftypes.Value
	if err := value.As(&children); err != nil {
		return err
	}

	for _, child := range children {
		parsed, err := ValueForType(s.Type)
		if err != nil {
			return err
		}

		if err := parsed.FromTerraform5Value(child); err != nil {
			return err
		}
		s.Values = append(s.Values, parsed)
	}
	return nil
}
