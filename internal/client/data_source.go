package client

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/liamcervante/terraform-provider-fakelocal/internal/values"
)

var _ tfsdk.DataSource = DataSource{}

type DataSource struct {
	Client Local
}

func (d DataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	value := values.ValueForType(req.Config.Schema.AttributeType())

	diags := req.Config.Get(ctx, &value)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := value.GetId()
	if err != nil {
		resp.Diagnostics.Append(
			diag.NewAttributeErrorDiagnostic(
				tftypes.NewAttributePath().WithAttributeName("id"),
				"failed to retrieve id", err.Error()))
		return
	}

	readResponse := values.ValueForType(req.Config.Schema.AttributeType())
	if err := d.Client.ReadResource(id, &readResponse); err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("read error", err.Error()))
		return
	}

	diags = resp.State.Set(ctx, &readResponse)
	resp.Diagnostics.Append(diags...)
}
