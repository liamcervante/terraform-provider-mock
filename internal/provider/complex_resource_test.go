package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccComplexResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: LoadFile(t, "testdata/complex/create.tf"),
				Check:  resource.TestCheckResourceAttr("complex_resource.test", "bool", "true"),
			},
			{
				Config: LoadFile(t, "testdata/complex/update.tf"),
				Check:  resource.TestCheckResourceAttr("complex_resource.test", "bool", "false"),
			},
			{
				Config: LoadFile(t, "testdata/complex/delete.tf"),
			},
		},
	})

	files, err := os.ReadDir("terraform.resource")
	if err != nil {
		t.Fatalf("could not verify resource directory: " + err.Error())
	}

	if len(files) != 0 {
		t.Fatalf("failed to tidy update after test")
	}
}

func LoadFile(t *testing.T, file string) string {
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err.Error())
	}

	return string(data)
}
