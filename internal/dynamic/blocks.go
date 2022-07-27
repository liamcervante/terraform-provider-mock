package dynamic

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

const (
	NestingModeList = "list"
	NestingModeSet  = "set"
)

type Block struct {
	Attributes map[string]Attribute `json:"attributes"`
	Blocks     map[string]Block     `json:"blocks"`
	Mode       string               `json:"mode"`
}

func (b Block) ToTerraformBlock() (tfsdk.Block, error) {
	tfAttributes := make(map[string]tfsdk.Attribute)
	tfBlocks := make(map[string]tfsdk.Block)

	for name, attribute := range b.Attributes {
		tfAttribute, err := attribute.ToTerraformAttribute()
		if err != nil {
			return tfsdk.Block{}, err
		}
		tfAttributes[name] = tfAttribute
	}

	for name, block := range b.Blocks {
		tfBlock, err := block.ToTerraformBlock()
		if err != nil {
			return tfsdk.Block{}, err
		}
		tfBlocks[name] = tfBlock
	}

	var nestingMode tfsdk.BlockNestingMode
	switch b.Mode {
	case "", NestingModeList:
		nestingMode = tfsdk.BlockNestingModeList
	case NestingModeSet:
		nestingMode = tfsdk.BlockNestingModeSet
	default:
		return tfsdk.Block{}, errors.New("invalid nesting mode: " + b.Mode)
	}

	return tfsdk.Block{
		Attributes:  tfAttributes,
		Blocks:      tfBlocks,
		NestingMode: nestingMode,
	}, nil
}

func blocksToTerraformBlocks(blocks map[string]Block) (map[string]tfsdk.Block, error) {
	tfBlocks := make(map[string]tfsdk.Block)
	for name, block := range blocks {
		tfBlock, err := block.ToTerraformBlock()
		if err != nil {
			return nil, err
		}
		tfBlocks[name] = tfBlock
	}
	return tfBlocks, nil
}
