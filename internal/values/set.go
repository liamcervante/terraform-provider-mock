package values

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ tftypes.ValueConverter = Set{}
var _ tftypes.ValueCreator = Set{}

type Set struct {
	Type   attr.Type `tfsdk:"type" json:"type"`
	Values []Value   `tfsdk:"values" json:"values"`
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
		parsed := ValueForType(s.Type)
		if err := parsed.FromTerraform5Value(child); err != nil {
			return err
		}
		s.Values = append(s.Values, parsed)
	}
	return nil
}
