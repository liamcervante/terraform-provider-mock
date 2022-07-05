package client

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.DataSource = DataSource{}

type DataSource struct {
	Client         Local
	CreateResource CreateResource
}

func (d DataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	data := d.CreateResource()

	diags := req.Config.Get(ctx, data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	readResponse := d.CreateResource()
	if err := d.Client.ReadResource(data.GetId(), &readResponse); err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("read error", err.Error()))
		return
	}

	diags = resp.State.Set(ctx, &readResponse)
	resp.Diagnostics.Append(diags...)
}
