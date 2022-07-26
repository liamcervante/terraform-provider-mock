package simple

import (
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	attributes = map[string]tfsdk.Attribute{
		"id": {
			Computed: true,
			Optional: true,
			PlanModifiers: tfsdk.AttributePlanModifiers{
				tfsdk.UseStateForUnknown(),
				tfsdk.RequiresReplace(),
			},
			Type: types.StringType,
		},
	}

	CoreAttributes = map[string]tfsdk.Attribute{
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

func Attributes() map[string]tfsdk.Attribute {
	attrs := make(map[string]tfsdk.Attribute)
	for name, attribute := range attributes {
		attrs[name] = attribute
	}
	for name, attribute := range CoreAttributes {
		attrs[name] = attribute
	}
	return attrs
}
