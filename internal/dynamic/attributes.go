package dynamic

import (
	"errors"

	"github.com/liamcervante/terraform-provider-mock/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type Attribute struct {
	Type     string `json:"type"`
	Optional bool   `json:"optional"`

	List   *Attribute           `json:"list,omitempty"`
	Map    *Attribute           `json:"map,omitempty"`
	Object map[string]Attribute `json:"object,omitempty"`
	Set    *Attribute           `json:"set,omitempty"`
}

func (a Attribute) ToTerraformAttribute() (tfsdk.Attribute, error) {
	switch a.Type {
	case types.Boolean:
		return tfsdk.Attribute{
			Optional: a.Optional,
			Required: !a.Optional,
			Type:     tftypes.BoolType,
		}, nil
	case types.Float:
		return tfsdk.Attribute{
			Optional: a.Optional,
			Required: !a.Optional,
			Type:     tftypes.Float64Type,
		}, nil
	case types.Integer:
		return tfsdk.Attribute{
			Optional: a.Optional,
			Required: !a.Optional,
			Type:     tftypes.Int64Type,
		}, nil
	case types.Number:
		return tfsdk.Attribute{
			Optional: a.Optional,
			Required: !a.Optional,
			Type:     tftypes.NumberType,
		}, nil
	case types.String:
		return tfsdk.Attribute{
			Optional: a.Optional,
			Required: !a.Optional,
			Type:     tftypes.StringType,
		}, nil
	case types.List:
		attribute, err := a.List.ToTerraformAttribute()
		if err != nil {
			return tfsdk.Attribute{}, err
		}

		return tfsdk.Attribute{
			Optional: a.Optional,
			Required: !a.Optional,
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
			Required: !a.Optional,
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
			Required: !a.Optional,
			Type: tftypes.SetType{
				ElemType: attribute.Type,
			},
		}, nil
	case types.Object:
		attributes, err := attributesToTerraformAttributes(a.Object)
		if err != nil {
			return tfsdk.Attribute{}, err
		}

		return tfsdk.Attribute{
			Optional:   a.Optional,
			Required:   !a.Optional,
			Attributes: tfsdk.SingleNestedAttributes(attributes),
		}, nil
	default:
		return tfsdk.Attribute{}, errors.New("unrecognized attribute type: " + a.Type)
	}
}

func attributesToTerraformAttributes(attributes map[string]Attribute) (map[string]tfsdk.Attribute, error) {
	tfAttributes := make(map[string]tfsdk.Attribute)
	for name, attribute := range attributes {
		tfAttribute, err := attribute.ToTerraformAttribute()
		if err != nil {
			return nil, err
		}
		tfAttributes[name] = tfAttribute
	}
	return tfAttributes, nil
}
