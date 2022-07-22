package client

import (
	"context"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/liamcervante/terraform-provider-fakelocal/internal/values"
)

var _ tfsdk.Resource = Resource{}
var _ tfsdk.ResourceWithImportState = Resource{}

type Resource struct {
	Client Local
}

func (r Resource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	value := values.ValueForType(req.Config.Schema.AttributeType())
	diags := req.Config.Get(ctx, &value)
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

	if err := value.SetId(id); err != nil {
		resp.Diagnostics.Append(
			diag.NewAttributeErrorDiagnostic(
				tftypes.NewAttributePath().WithAttributeName("id"),
				"failed to store id", err.Error()))
		return
	}
	if err := r.Client.WriteResource(id, value); err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("write error", err.Error()))
		return
	}

	diags = resp.State.Set(ctx, &value)
	resp.Diagnostics.Append(diags...)
}

func (r Resource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	value := values.ValueForType(req.State.Schema.AttributeType())

	diags := req.State.Get(ctx, &value)
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

	readResponse := &values.Value{}
	if err := r.Client.ReadResource(id, &readResponse); err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("read error", err.Error()))
		return
	}

	diags = resp.State.Set(ctx, &readResponse)
	resp.Diagnostics.Append(diags...)
}

func (r Resource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	value := values.ValueForType(req.Plan.Schema.AttributeType())

	diags := req.Plan.Get(ctx, &value)
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

	if err := r.Client.UpdateResource(id, value); err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("update error", err.Error()))
		return
	}

	diags = resp.State.Set(ctx, &value)
	resp.Diagnostics.Append(diags...)
}

func (r Resource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	value := values.ValueForType(req.State.Schema.AttributeType())
	diags := req.State.Get(ctx, &value)
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

	if err := r.Client.DeleteResource(id); err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("delete error", err.Error()))
		return
	}
}

func (r Resource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
