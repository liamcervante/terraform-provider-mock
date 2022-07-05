package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/liamcervante/terraform-provider-fakelocal/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.DataSourceType = complexDataSourceType{}

type complexDataSourceType struct{}

func (t complexDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: client.ComplexAttributes,
	}, nil
}

func (t complexDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return client.DataSource{
		Client:         provider.client,
		CreateResource: client.NewComplex,
	}, diags
}
