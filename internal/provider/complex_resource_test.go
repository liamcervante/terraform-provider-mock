package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComplexResource(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/complex/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mock_complex_resource.test", "string", "hello"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list.#", "1"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list.0.integer", "0"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "object.bool", "true"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set.#", "2"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set.0.integer", "0"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set.1.integer", "1")),
			},
			{
				Config: LoadFile(t, "testdata/complex/update/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mock_complex_resource.test", "string", "hello"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list.#", "2"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list.0.integer", "0"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list.1.integer", "1"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "object.string", "world"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set.#", "1"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set.0.integer", "0")),
			},
			{
				Config: LoadFile(t, "testdata/complex/delete/main.tf"),
			},
		},
	})
}
func TestAccComplexResourceWithBlocks(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/complex_block/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list_block.#", "2"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list_block.0.integer", "0"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list_block.1.integer", "1"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set_block.#", "1"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set_block.0.integer", "0")),
			},
			{
				Config: LoadFile(t, "testdata/complex_block/update/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list_block.#", "1"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "list_block.0.integer", "0"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set_block.#", "2"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set_block.0.integer", "0"),
					resource.TestCheckResourceAttr("mock_complex_resource.test", "set_block.1.integer", "1")),
			},
			{
				Config: LoadFile(t, "testdata/complex/delete/main.tf"),
			},
		},
	})
}
