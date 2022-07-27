package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func ProviderFactories(resources string) map[string]func() (tfprotov6.ProviderServer, error) {
	provider := NewForTesting("test", resources)()
	return map[string]func() (tfprotov6.ProviderServer, error){
		"mock": providerserver.NewProtocol6WithError(provider),
	}
}

func LoadFile(t *testing.T, file string) string {
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err.Error())
	}

	return string(data)
}

func CleanupTestingDirectories(t *testing.T) {
	files, err := os.ReadDir("terraform.resource")
	if err != nil {
		if os.IsNotExist(err) {
			return // Then it's fine.
		}

		t.Fatalf("could not verify resource directory: " + err.Error())
	}
	defer os.Remove("terraform.resource")

	if len(files) != 0 {
		t.Fatalf("failed to tidy update after test")
	}
}
