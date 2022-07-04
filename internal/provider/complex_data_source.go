package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/liamcervante/terraform-provider-fakelocal/internal/resource"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.DataSourceType = complexDataSourceType{}

type complexDataSourceType struct{}

func (t complexDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: resource.ComplexAttributes,
	}, nil
}

func (t complexDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return resource.DataSource{
		Client:         provider.client,
		CreateResource: resource.NewComplex,
	}, diags
}
