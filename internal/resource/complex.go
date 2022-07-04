package resource

import (
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	ComplexAttributes = map[string]tfsdk.Attribute{
		"id": {
			Computed: true,
			PlanModifiers: tfsdk.AttributePlanModifiers{
				tfsdk.UseStateForUnknown(),
			},
			Type: types.StringType,
		},
		"bool": {
			Optional: true,
			Type:     types.BoolType,
		},
		"number": {
			Optional: true,
			Type:     types.NumberType,
		},
		"string": {
			Optional: true,
			Type:     types.StringType,
		},
		"float": {
			Optional: true,
			Type:     types.Float64Type,
		},
		"integer": {
			Optional: true,
			Type:     types.Int64Type,
		},
	}
)

type Complex struct {
	Id types.String `tfsdk:"id" json:"id"`

	Boolean types.Bool    `tfsdk:"bool" json:"bool"`
	Number  types.Number  `tfsdk:"number" json:"number"`
	String  types.String  `tfsdk:"string" json:"string"`
	Float   types.Float64 `tfsdk:"float" json:"float"`
	Integer types.Int64   `tfsdk:"integer" json:"integer"`
}

func (c *Complex) GetId() string {
	return c.Id.Value
}

func (c *Complex) SetId(id string) {
	c.Id = types.String{Value: id}
}

func NewComplex() Data {
	return &Complex{}
}
