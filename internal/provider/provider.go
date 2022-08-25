package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/liamcervante/terraform-provider-mock/internal/dynamic"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/liamcervante/terraform-provider-mock/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.Provider = &provider{}

// provider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type provider struct {
	client client.Local

	// configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the provider was previously configured.
	configured bool

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string

	// reader will read the dynamic resource definitions in the GetResources and
	// GetDataSources functions.
	reader dynamic.Reader

	// useOnlyState tells the provider to read data from state only, and not
	// write any of the data to disk.
	useOnlyState bool
}

// providerData can be used to store client from the Terraform configuration.
type providerData struct {
	ResourceDir  types.String `tfsdk:"resource_directory"`
	DataDir      types.String `tfsdk:"data_directory"`
	UseOnlyState types.Bool   `tfsdk:"use_only_state"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	tflog.Trace(ctx, "provider.Configure")
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	p.client = client.Local{}

	if data.ResourceDir.Null {
		p.client.ResourceDirectory = "terraform.resource"
	} else {
		p.client.ResourceDirectory = data.ResourceDir.Value
	}
	ctx = tflog.With(ctx, "resource_directory", p.client.ResourceDirectory)

	if data.DataDir.Null {
		p.client.DataDirectory = "terraform.data"
	} else {
		p.client.DataDirectory = data.DataDir.Value
	}
	ctx = tflog.With(ctx, "data_directory", p.client.DataDirectory)

	if data.UseOnlyState.Null {
		p.useOnlyState = false
	} else {
		p.useOnlyState = data.UseOnlyState.Value
	}
	ctx = tflog.With(ctx, "use_only_state", p.useOnlyState)

	p.configured = true
	tflog.Trace(ctx, "provider.Configured = true")
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	wd, _ := os.Getwd()
	ctx = tflog.With(ctx, "working_directory", wd)
	tflog.Trace(ctx, "provider.GetResources")

	dynamicResources, err := p.reader.Read()
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("Failed to read dynamic resources", err.Error())
		return nil, diags
	}

	resources := make(map[string]tfsdk.ResourceType)
	for name, resource := range dynamicResources {
		resources[name] = dynamicResourceType{resource}
	}

	resources["mock_complex_resource"] = complexResourceType{}
	resources["mock_simple_resource"] = simpleResourceType{}

	return resources, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	wd, _ := os.Getwd()
	ctx = tflog.With(ctx, "working_directory", wd)
	tflog.Trace(ctx, "provider.GetDataSources")

	dynamicResources, err := p.reader.Read()
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("Failed to read dynamic resources", err.Error())
		return nil, diags
	}

	sources := make(map[string]tfsdk.DataSourceType)
	for name, resource := range dynamicResources {
		sources[name] = dynamicDataSourceType{resource}
	}

	sources["mock_complex_resource"] = complexDataSourceType{}
	sources["mock_simple_resource"] = simpleDataSourceType{}

	return sources, nil
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	tflog.Trace(ctx, "provider.GetSchema")
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"resource_directory": {
				Optional: true,
				Type:     types.StringType,
			},
			"data_directory": {
				Optional: true,
				Type:     types.StringType,
			},
			"use_only_state": {
				Optional: true,
				Type:     types.BoolType,
			},
		},
	}, nil
}

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
			reader:  dynamic.FileReader{File: "dynamic_resources.json"},
		}
	}
}

func NewForTesting(version string, resources string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
			reader:  dynamic.StringReader{Data: resources},
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
func convertProviderType(in tfsdk.Provider) (provider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*provider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return provider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return provider{}, diags
	}

	return *p, diags
}
