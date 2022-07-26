package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var (
	dynamicResources = `{"fakelocal_dynamic_resource":{"type":"object","object":{"integer":{"type":"integer","optional":true,"required":false}}}}`
)

func TestAccDynamicResource(t *testing.T) {
	defer CleanupTestingDirectories(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(dynamicResources),
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/dynamic/create.tf"),
				Check:  resource.TestCheckResourceAttr("fakelocal_dynamic_resource.test", "integer", "0"),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/update.tf"),
				Check:  resource.TestCheckResourceAttr("fakelocal_dynamic_resource.test", "integer", "1"),
			},
			{
				Config: LoadFile(t, "testdata/dynamic/delete.tf"),
			},
		},
	})
}
