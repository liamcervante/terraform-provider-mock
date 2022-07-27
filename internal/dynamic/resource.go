package dynamic

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type Resource struct {
	Attributes map[string]Attribute
	Blocks     map[string]Block
}

func (resource Resource) ToTerraformSchema(computed bool) (tfsdk.Schema, diag.Diagnostics) {
	var diags diag.Diagnostics

	if _, ok := resource.Attributes["id"]; ok {
		diags.AddError(
			"Found `id` value in top level object",
			"Top level dynamic objects cannot define a value called `id` as the provider will generate an ID for them.",
		)
		return tfsdk.Schema{}, diags
	}

	attributes, err := attributesToTerraformAttributes(resource.Attributes)
	if err != nil {
		diags.AddError("Failed to parse dynamic attributes", err.Error())
		return tfsdk.Schema{}, diags
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

	blocks, err := blocksToTerraformBlocks(resource.Blocks)
	if err != nil {
		diags.AddError("Failed to parse dynamic blocks", err.Error())
		return tfsdk.Schema{}, diags
	}

	return tfsdk.Schema{
		Attributes: attributes,
		Blocks:     blocks,
	}, diags
}
