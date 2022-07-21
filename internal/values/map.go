package values

import (
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/liamcervante/terraform-provider-fakelocal/internal/types"
)

var _ tftypes.ValueConverter = Map{}
var _ tftypes.ValueCreator = Map{}

type Map struct {
	Type   types.Type       `json:"type"`
	Values map[string]Value `json:"values"`
}

func (m Map) ToTerraform5Value() (interface{}, error) {
	children := make(map[string]interface{})
	for name, value := range m.Values {
		child, err := value.ToTerraform5Value()
		if err != nil {
			return nil, err
		}
		children[name] = child
	}
	return children, nil
}

func (m Map) FromTerraform5Value(value tftypes.Value) error {
	var children map[string]tftypes.Value
	if err := value.As(&children); err != nil {
		return err
	}

	for name, child := range children {
		parsed, err := ValueForType(m.Type)
		if err != nil {
			return err
		}

		if err := parsed.FromTerraform5Value(child); err != nil {
			return err
		}
		m.Values[name] = parsed
	}
	return nil
}
