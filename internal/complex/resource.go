package complex

import (
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/liamcervante/terraform-provider-fakelocal/internal/client"
)

var (
	Attributes = map[string]tfsdk.Attribute{
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

type Resource struct {
	Id types.String `tfsdk:"id" json:"id"`

	Boolean types.Bool    `tfsdk:"bool" json:"bool"`
	Number  types.Number  `tfsdk:"number" json:"number"`
	String  types.String  `tfsdk:"string" json:"string"`
	Float   types.Float64 `tfsdk:"float" json:"float"`
	Integer types.Int64   `tfsdk:"integer" json:"integer"`
}

func (r *Resource) GetId() string {
	return r.Id.Value
}

func (r *Resource) SetId(id string) {
	r.Id = types.String{Value: id}
}

func New() client.Data {
	return &Resource{}
}
