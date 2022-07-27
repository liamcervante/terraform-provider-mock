package values

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/liamcervante/terraform-provider-fakelocal/internal/types"
)

var _ tftypes.ValueConverter = &Resource{}
var _ tftypes.ValueCreator = &Resource{}

type Resource struct {
	Types  map[string]*types.Type `tfsdk:"types" json:"types"`
	Values map[string]Value       `tfsdk:"values" json:"values"`
}

func (r *Resource) SetId(id string) {
	r.Types["id"] = &types.Type{Type: types.String}
	r.Values["id"] = Value{String: &id}
}

func (r Resource) GetId() (string, error) {
	if id, ok := r.Values["id"]; ok && id.String != nil {
		return *id.String, nil
	}
	return "", errors.New("ID not set or ID not set as String type")
}

func (r Resource) Type() *types.Type {
	return &types.Type{
		Type:       types.Object,
		ObjectType: r.Types,
	}
}

func (r *Resource) ToTerraform5Value() (interface{}, error) {
	value, _, err := objectToTerraform5Value(&r.Values, r.Type())

	//fmt.Printf("ToTerraform5Value\n")
	//fmt.Printf("FROM\n\t%v\n", toJSON(r.Values))
	//fmt.Printf("TO\n\t%v\n\n", tftypes.NewValue(typ, value))
	return value, err
}

func (r *Resource) FromTerraform5Value(value tftypes.Value) error {
	var err error

	typ, err := types.FromTerraform5Type(value.Type())
	if err != nil {
		return err
	}

	values, err := fromTerraform5Value(value)
	if err != nil {
		return err
	}

	r.Values = *values.Object
	r.Types = typ.ObjectType

	//fmt.Printf("FromTerraform5Value\n")
	//fmt.Printf("FROM\n\t%v\n", value)
	//fmt.Printf("TO\n\t%v\n\n", toJSON(r.Values))

	return nil
}

//func toJSON(obj map[string]Value) string {
//	data, err := json.Marshal(obj)
//	if err != nil {
//		panic(err)
//	}
//	return string(data)
//}
