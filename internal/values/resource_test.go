package values

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/liamcervante/terraform-provider-mock/internal/types"
)

func TestResource_Symmetry(t *testing.T) {
	testCases := []struct {
		TestCase string
		Resource Resource
	}{
		{
			TestCase: "basic",
			Resource: Resource{
				Types: map[string]*types.Type{
					"number": {Type: types.Number},
				},
				Values: map[string]Value{
					"number": {Number: big.NewFloat(0)},
				},
			},
		},
		{
			TestCase: "missing_object",
			Resource: Resource{
				Types: map[string]*types.Type{
					"object": {
						Type: types.Object,
						ObjectType: map[string]*types.Type{
							"number": {Type: types.Number},
						},
					},
				},
				Values: map[string]Value{},
			},
		},
		{
			TestCase: "missing_list",
			Resource: Resource{
				Types: map[string]*types.Type{
					"list": {
						Type:        types.List,
						ElementType: &types.Type{Type: types.Number},
					},
				},
				Values: map[string]Value{},
			},
		},
		{
			TestCase: "empty_list",
			Resource: Resource{
				Types: map[string]*types.Type{
					"list": {
						Type:        types.List,
						ElementType: &types.Type{Type: types.Number},
					},
				},
				Values: map[string]Value{
					"list": {
						List: &[]Value{},
					},
				},
			},
		},
		{
			TestCase: "missing_map",
			Resource: Resource{
				Types: map[string]*types.Type{
					"map": {
						Type:        types.Map,
						ElementType: &types.Type{Type: types.Number},
					},
				},
				Values: map[string]Value{},
			},
		},
		{
			TestCase: "missing_set",
			Resource: Resource{
				Types: map[string]*types.Type{
					"set": {
						Type:        types.Set,
						ElementType: &types.Type{Type: types.Number},
					},
				},
				Values: map[string]Value{},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.TestCase, func(t *testing.T) {
			CheckSymmetry(t, testCase.Resource)
		})
	}
}

func toJson(t *testing.T, obj Resource) string {
	data, err := json.Marshal(obj)
	if err != nil {
		t.Fatalf(err.Error())
	}
	return string(data)
}

func CheckResourceEqual(t *testing.T, expected, actual Resource) {
	expectedStr := toJson(t, expected)
	actualStr := toJson(t, actual)
	if expectedStr != actualStr {
		t.Fatalf("expected:\n\t%s\nactual:\n\t%s\n", expectedStr, actualStr)
	}
}

func CheckSymmetry(t *testing.T, resource Resource) {
	raw, err := resource.ToTerraform5Value()
	if err != nil {
		t.Fatalf(err.Error())
	}

	tfType, err := resource.Type().ToTerraform5Type()
	if err != nil {
		t.Fatalf(err.Error())
	}

	value := tftypes.NewValue(tfType, raw)
	actual := Resource{}
	if err := actual.FromTerraform5Value(value); err != nil {
		t.Fatalf(err.Error())
	}

	CheckResourceEqual(t, resource, actual)
}
