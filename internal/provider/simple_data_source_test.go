package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSimpleDataSource(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(""),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/simple_datasource/get/main.tf"),
				Check:  resource.TestCheckResourceAttr("data.fakelocal_simple_resource.test", "integer", "0"),
			},
		},
	})
}
