package client

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/liamcervante/terraform-provider-mock/internal/values"
)

var _ tfsdk.DataSource = DataSource{}

type DataSource struct {
	Client Local
}

func (d DataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	tflog.Trace(ctx, "DataSource.Read")

	resource := &values.Resource{}

	diags := req.Config.Get(ctx, &resource)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := resource.GetId()
	if err != nil {
		resp.Diagnostics.Append(
			diag.NewAttributeErrorDiagnostic(
				tftypes.NewAttributePath().WithAttributeName("id"),
				"failed to retrieve id", err.Error()))
		return
	}
	ctx = tflog.With(ctx, "id", id)

	data, err := d.Client.ReadDataSource(ctx, id)
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("read error", err.Error()))
		return
	}

	if data == nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("target data source does not exist", fmt.Sprintf("data source at %s could not be found in data directory (%s)", id, d.Client.DataDirectory)))
		return
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}
