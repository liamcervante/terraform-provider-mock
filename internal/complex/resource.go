package complex

import (
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/liamcervante/terraform-provider-fakelocal/internal/simple"
)

func Attributes(maxDepth int) map[string]tfsdk.Attribute {
	attrs := make(map[string]tfsdk.Attribute)
	attrs["id"] = tfsdk.Attribute{
		Computed: true,
		Optional: true,
		PlanModifiers: tfsdk.AttributePlanModifiers{
			tfsdk.UseStateForUnknown(),
			tfsdk.RequiresReplace(),
		},
		Type: types.StringType,
	}

	for name, attribute := range simple.CoreAttributes {
		attrs[name] = attribute
	}

	for name, attribute := range object(0, maxDepth) {
		attrs[name] = attribute
	}

	return attrs
}

func Blocks(maxDepth int) map[string]tfsdk.Block {
	return map[string]tfsdk.Block{} // blocks(0, maxDepth)
}

func blocks(depth, maxDepth int) map[string]tfsdk.Block {
	if depth == maxDepth {
		return nil
	}

	blks := make(map[string]tfsdk.Block)

	blks["list_block"] = tfsdk.Block{
		Attributes:  object(depth+1, maxDepth),
		Blocks:      blocks(depth+1, maxDepth),
		NestingMode: tfsdk.BlockNestingModeList,
	}

	blks["set_block"] = tfsdk.Block{
		Attributes:  object(depth+1, maxDepth),
		Blocks:      blocks(depth+1, maxDepth),
		NestingMode: tfsdk.BlockNestingModeSet,
	}

	return blks
}

func object(depth, maxDepth int) map[string]tfsdk.Attribute {
	attrs := make(map[string]tfsdk.Attribute)

	for name, attribute := range simple.CoreAttributes {
		attrs[name] = attribute
	}

	if depth < maxDepth {
		attrs["list"] = tfsdk.Attribute{
			Optional:   true,
			Attributes: tfsdk.ListNestedAttributes(object(depth+1, maxDepth)),
		}
		attrs["map"] = tfsdk.Attribute{
			Optional:   true,
			Attributes: tfsdk.MapNestedAttributes(object(depth+1, maxDepth)),
		}
		attrs["object"] = tfsdk.Attribute{
			Optional:   true,
			Attributes: tfsdk.SingleNestedAttributes(object(depth+1, maxDepth)),
		}
		attrs["set"] = tfsdk.Attribute{
			Optional:   true,
			Attributes: tfsdk.SetNestedAttributes(object(depth+1, maxDepth)),
		}
	}

	return attrs
}
