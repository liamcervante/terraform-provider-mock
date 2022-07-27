package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDynamicResource(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(LoadFile(t, "testdata/dynamic/dynamic_resources.json")),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/dynamic/create/main.tf"),
				Check:  resource.TestCheckResourceAttr("fakelocal_dynamic_resource.test", "integer", "0"),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/update/main.tf"),
				Check:  resource.TestCheckResourceAttr("fakelocal_dynamic_resource.test", "integer", "1"),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/delete/main.tf"),
			},
		},
	})
}

func TestAccDynamicResourceWithBlocks(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(LoadFile(t, "testdata/dynamic_block/dynamic_resources.json")),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/dynamic_block/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("fakelocal_dynamic_resource.test", "integer", "0"),
					resource.TestCheckResourceAttr("fakelocal_dynamic_resource.test", "nested_list.#", "1"),
					resource.TestCheckResourceAttr("fakelocal_dynamic_resource.test", "nested_list.0.integer", "0")),
			},
			{
				Config: LoadFile(t, "testdata/dynamic_block/update/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("fakelocal_dynamic_resource.test", "integer", "0"),
					resource.TestCheckResourceAttr("fakelocal_dynamic_resource.test", "nested_list.#", "2"),
					resource.TestCheckResourceAttr("fakelocal_dynamic_resource.test", "nested_list.0.integer", "0"),
					resource.TestCheckResourceAttr("fakelocal_dynamic_resource.test", "nested_list.1.integer", "1")),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/delete/main.tf"),
			},
		},
	})
}

func TestAccDynamicResourceWithId(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(LoadFile(t, "testdata/dynamic/dynamic_resources.json")),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/dynamic_with_id/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("fakelocal_dynamic_resource.test", "integer", "0"),
					resource.TestCheckResourceAttr("fakelocal_dynamic_resource.test", "id", "my_id")),
			},
			{
				Config: LoadFile(t, "testdata/dynamic_with_id/update/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("fakelocal_dynamic_resource.test", "integer", "1"),
					resource.TestCheckResourceAttr("fakelocal_dynamic_resource.test", "id", "my_id")),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/delete/main.tf"),
			},
		},
	})
}
