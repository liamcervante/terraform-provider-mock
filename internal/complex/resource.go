package complex

import (
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
