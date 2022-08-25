package client

import (
	"context"
	"github.com/liamcervante/terraform-provider-mock/internal/types"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/liamcervante/terraform-provider-mock/internal/values"
)

var _ tfsdk.Resource = Resource{}
var _ tfsdk.ResourceWithImportState = Resource{}

type Resource struct {
	Client       Local
	UseOnlyState bool
}

func (r Resource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Trace(ctx, "Resource.Create")

	resource := &values.Resource{}
	diags := req.Config.Get(ctx, &resource)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := resource.GetId()
	if err != nil {
		tflog.Info(ctx, "could not retrieve id so generating a new one", map[string]interface{}{
			"resource.GetId": err.Error(),
		})

		if id, err = uuid.GenerateUUID(); err != nil {
			resp.Diagnostics.Append(
				diag.NewAttributeErrorDiagnostic(
					tftypes.NewAttributePath().WithAttributeName("id"),
					"failed to generate id", err.Error()))
			return
		}
		resource.SetId(id)

	}
	ctx = tflog.With(ctx, "id", id)

	// Don't write the resource to disk if we're only using that state.
	if !r.UseOnlyState {
		if err := r.Client.WriteResource(ctx, id, resource); err != nil {
			resp.Diagnostics.Append(diag.NewErrorDiagnostic("write error", err.Error()))
			return
		}
	}

	diags = resp.State.Set(ctx, resource)
	resp.Diagnostics.Append(diags...)
}

func (r Resource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	tflog.Trace(ctx, "Resource.Read")

	resource := &values.Resource{}

	diags := req.State.Get(ctx, &resource)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If we're only using the state then don't try and read this resource from
	// disk.
	if r.UseOnlyState {
		diags = resp.State.Set(ctx, resource)
		resp.Diagnostics.Append(diags...)
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
	tflog.With(ctx, "id", id)

	data, err := r.Client.ReadResource(ctx, id)
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("read error", err.Error()))
		return
	}

	if data == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	typ, err := types.FromTerraform5Type(req.State.Schema.TerraformType(ctx))
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("internal error", err.Error()))
		return
	}
	data.Types = typ.ObjectType

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r Resource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	tflog.Trace(ctx, "Resource.Update")
	resource := &values.Resource{}

	diags := req.Plan.Get(ctx, &resource)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Don't write the updated object to disk if we're only using state.
	if !r.UseOnlyState {
		id, err := resource.GetId()
		if err != nil {
			resp.Diagnostics.Append(
				diag.NewAttributeErrorDiagnostic(
					tftypes.NewAttributePath().WithAttributeName("id"),
					"failed to retrieve id", err.Error()))
			return
		}
		ctx = tflog.With(ctx, "id", id)

		if err := r.Client.UpdateResource(ctx, id, resource); err != nil {
			resp.Diagnostics.Append(diag.NewErrorDiagnostic("update error", err.Error()))
			return
		}
	}

	diags = resp.State.Set(ctx, resource)
	resp.Diagnostics.Append(diags...)
}

func (r Resource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	tflog.Trace(ctx, "Resource.Delete")

	if r.UseOnlyState {
		// Then we have nothing to do here, nothing to delete and the SDK
		// removes it from the state automatically.
		return
	}

	resource := &values.Resource{}
	diags := req.State.Get(ctx, &resource)
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

	if err := r.Client.DeleteResource(ctx, id); err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("delete error", err.Error()))
		return
	}
}

func (r Resource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tflog.Trace(ctx, "Resource.ImportState")
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
