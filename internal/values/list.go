package values

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ tftypes.ValueConverter = List{}
var _ tftypes.ValueCreator = List{}

type List struct {
	Type   attr.Type `tfsdk:"type" json:"type"`
	Values []Value   `tfsdk:"values" json:"values"`
}

func (l List) ToTerraform5Value() (interface{}, error) {
	var children []interface{}
	for _, value := range l.Values {
		child, err := value.ToTerraform5Value()
		if err != nil {
			return nil, err
		}
		children = append(children, child)
	}
	return children, nil
}

func (l List) FromTerraform5Value(value tftypes.Value) error {
	var children []tftypes.Value
	if err := value.As(&children); err != nil {
		return err
	}

	for _, child := range children {
		parsed := ValueForType(l.Type)
		if err := parsed.FromTerraform5Value(child); err != nil {
			return err
		}
		l.Values = append(l.Values, parsed)
	}
	return nil
}
