package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSimpleResource(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/simple/create/main.tf"),
				Check:  resource.TestCheckResourceAttr("fakelocal_simple_resource.test", "integer", "0"),
			},
			{
				Config: LoadFile(t, "testdata/simple/update/main.tf"),
				Check:  resource.TestCheckResourceAttr("fakelocal_simple_resource.test", "integer", "1"),
			},
			{
				Config: LoadFile(t, "testdata/simple/delete/main.tf"),
			},
		},
	})
}

func TestAccSimpleResourceWithDrift(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/simple/create/main.tf"),
				Check: func(state *terraform.State) error {
					id := state.RootModule().Resources["fakelocal_simple_resource.test"].Primary.Attributes["id"]
					return os.Remove(fmt.Sprintf("terraform.resource/%s.json", id))
				},
				ExpectNonEmptyPlan: true,
			},
			{
				Config: LoadFile(t, "testdata/simple/update/main.tf"),
				Check:  resource.TestCheckResourceAttr("fakelocal_simple_resource.test", "integer", "1"),
			},
			{
				Config: LoadFile(t, "testdata/simple/delete/main.tf"),
			},
		},
	})
}

func TestAccSimpleResourceWithId(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/simple_with_id/create/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("fakelocal_simple_resource.test", "integer", "0"),
					resource.TestCheckResourceAttr("fakelocal_simple_resource.test", "id", "my_id")),
			},
			{
				Config: LoadFile(t, "testdata/simple_with_id/update/main.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("fakelocal_simple_resource.test", "integer", "1"),
					resource.TestCheckResourceAttr("fakelocal_simple_resource.test", "id", "my_id")),
			},
			{
				Config: LoadFile(t, "testdata/simple/delete/main.tf"),
			},
		},
	})
}
