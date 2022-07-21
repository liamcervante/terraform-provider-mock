package dynamic

import "github.com/hashicorp/terraform-plugin-framework/types"

type Resource struct {
	Id types.String `tfsdk:"id" json:"id"`
}

func (r *Resource) GetId() string {
	return r.Id.Value
}

func (r *Resource) SetId(id string) {
	r.Id = types.String{Value: id}
}
