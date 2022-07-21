package dynamic

import (
	"errors"

	"github.com/liamcervante/terraform-provider-fakelocal/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type Attribute struct {
	Type     string `json:"type"`
	Optional bool   `json:"optional"`

	List   *List   `json:"list"`
	Map    *Map    `json:"map"`
	Object *Object `json:"object"`
	Set    *Set    `json:"set"`
}

type List struct {
	Attribute Attribute `json:"attribute"`
}

type Map struct {
	Attribute Attribute `json:"attribute"`
}

type Object struct {
	Attributes map[string]Attribute `json:"attributes"`
}

type Set struct {
	Attribute Attribute `json:"attribute"`
}

func (a Attribute) ToTerraformAttribute() (tfsdk.Attribute, error) {
	switch a.Type {
	case types.Boolean:
		return tfsdk.Attribute{
			Optional: a.Optional,
			Type:     tftypes.BoolType,
		}, nil
	case types.Float:
		return tfsdk.Attribute{
			Optional: a.Optional,
			Type:     tftypes.Float64Type,
		}, nil
	case types.Integer:
		return tfsdk.Attribute{
			Optional: a.Optional,
			Type:     tftypes.Int64Type,
		}, nil
	case types.Number:
		return tfsdk.Attribute{
			Optional: a.Optional,
			Type:     tftypes.NumberType,
		}, nil
	case types.String:
		return tfsdk.Attribute{
			Optional: a.Optional,
			Type:     tftypes.StringType,
		}, nil
	case types.List:
		attribute, err := a.List.ToTerraformAttribute()
		if err != nil {
			return tfsdk.Attribute{}, err
		}

		return tfsdk.Attribute{
			Optional: a.Optional,
			Type: tftypes.ListType{
				ElemType: attribute.Type,
			},
		}, nil
	case types.Map:
		attribute, err := a.Map.ToTerraformAttribute()
		if err != nil {
			return tfsdk.Attribute{}, err
		}

		return tfsdk.Attribute{
			Optional: a.Optional,
			Type: tftypes.MapType{
				ElemType: attribute.Type,
			},
		}, nil
	case types.Set:
		attribute, err := a.Set.ToTerraformAttribute()
		if err != nil {
			return tfsdk.Attribute{}, err
		}

		return tfsdk.Attribute{
			Optional: a.Optional,
			Type: tftypes.SetType{
				ElemType: attribute.Type,
			},
		}, nil
	case types.Object:
		attributes, err := a.Object.ToTerraformAttribute()
		if err != nil {
			return tfsdk.Attribute{}, err
		}

		attrTypes := make(map[string]attr.Type)
		for name, attribute := range attributes {
			attrTypes[name] = attribute.Type
		}

		return tfsdk.Attribute{
			Optional: a.Optional,
			Type: tftypes.ObjectType{
				AttrTypes: attrTypes,
			},
		}, nil
	default:
		return tfsdk.Attribute{}, errors.New("unrecognized attribute type: " + a.Type)
	}
}

func (l List) ToTerraformAttribute() (tfsdk.Attribute, error) {
	return l.Attribute.ToTerraformAttribute()
}

func (m Map) ToTerraformAttribute() (tfsdk.Attribute, error) {
	return m.Attribute.ToTerraformAttribute()
}

func (o Object) ToTerraformAttribute() (map[string]tfsdk.Attribute, error) {
	attributes := make(map[string]tfsdk.Attribute)
	for name, attribute := range o.Attributes {
		tfAttribute, err := attribute.ToTerraformAttribute()
		if err != nil {
			return nil, err
		}
		attributes[name] = tfAttribute
	}
	return attributes, nil
}

func (s Set) ToTerraformAttribute() (tfsdk.Attribute, error) {
	return s.Attribute.ToTerraformAttribute()
}
