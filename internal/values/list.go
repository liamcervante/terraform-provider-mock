package values

import (
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/liamcervante/terraform-provider-fakelocal/internal/types"
)

var _ tftypes.ValueConverter = List{}
var _ tftypes.ValueCreator = List{}

type List struct {
	Type   types.Type `json:"type"`
	Values []Value    `json:"values"`
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
		parsed, err := ValueForType(l.Type)
		if err != nil {
			return err
		}

		if err := parsed.FromTerraform5Value(child); err != nil {
			return err
		}
		l.Values = append(l.Values, parsed)
	}
	return nil
}
