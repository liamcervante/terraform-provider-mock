package client

import (
	"context"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/liamcervante/terraform-provider-fakelocal/internal/types"
)

var _ tfsdk.Resource = Resource{}
var _ tfsdk.ResourceWithImportState = Resource{}

type Resource struct {
	Client         Local
	CreateResource CreateResource
}

func (r Resource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	data := r.CreateResource()

	types.CreateRootType(req.Config.Schema)

	diags := req.Config.Get(ctx, data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		resp.Diagnostics.Append(
			diag.NewAttributeErrorDiagnostic(
				tftypes.NewAttributePath().WithAttributeName("id"),
				"failed to generate id", err.Error()))
		return
	}

	data.SetId(id)
	if err := r.Client.WriteResource(id, data); err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("write error", err.Error()))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r Resource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	data := r.CreateResource()

	diags := req.State.Get(ctx, data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	readResponse := r.CreateResource()
	if err := r.Client.ReadResource(data.GetId(), &readResponse); err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("read error", err.Error()))
		return
	}

	diags = resp.State.Set(ctx, &readResponse)
	resp.Diagnostics.Append(diags...)
}

func (r Resource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	data := r.CreateResource()

	diags := req.Plan.Get(ctx, data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.Client.UpdateResource(data.GetId(), data); err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("update error", err.Error()))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r Resource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	data := r.CreateResource()

	diags := req.State.Get(ctx, data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.Client.DeleteResource(data.GetId()); err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("delete error", err.Error()))
		return
	}
}

func (r Resource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
