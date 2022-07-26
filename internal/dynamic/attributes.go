package dynamic

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/liamcervante/terraform-provider-fakelocal/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
		attributes, err := a.ObjectToTerraformAttribute()
		if err != nil {
			return tfsdk.Attribute{}, err
		}

		attrTypes := make(map[string]attr.Type)
		for name, attribute := range attributes {
			attrTypes[name] = attribute.Type
		}

		return tfsdk.Attribute{
			Optional: a.Optional,
			Required: !a.Optional,
			Type: tftypes.ObjectType{
				AttrTypes: attrTypes,
			},
		}, nil
	default:
		return tfsdk.Attribute{}, errors.New("unrecognized attribute type: " + a.Type)
	}
}

func (a Attribute) ObjectToTerraformAttribute() (map[string]tfsdk.Attribute, error) {
	attributes := make(map[string]tfsdk.Attribute)
	for name, attribute := range a.Object {
		tfAttribute, err := attribute.ToTerraformAttribute()
		if err != nil {
			return nil, err
		}
		attributes[name] = tfAttribute
	}
	return attributes, nil
}

func (a Attribute) ToTerraformSchema(computed bool) (tfsdk.Schema, diag.Diagnostics) {
	var diags diag.Diagnostics

	if a.Type != types.Object {
		diags.AddError(
			"Invalid attribute type",
			"You can only turn objects into terraform schemas")

		return tfsdk.Schema{}, diags
	}

	if _, ok := a.Object["id"]; ok {
		diags.AddError(
			"Found `id` value in top level object",
			"Top level dynamic objects cannot define a value called `id` as the provider will generate an ID for them.",
		)

		return tfsdk.Schema{}, diags
	}

	attributes, err := a.ObjectToTerraformAttribute()
	if err != nil {
		diags.AddError("Failed to parse dynamic attributes", err.Error())
	}

	attributes["id"] = tfsdk.Attribute{
		Required: !computed,
		Optional: computed,
		Computed: computed,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			tfsdk.UseStateForUnknown(),
			tfsdk.RequiresReplace(),
		},
		Type: tftypes.StringType,
	}

	return tfsdk.Schema{
		Attributes: attributes,
	}, diags
}
