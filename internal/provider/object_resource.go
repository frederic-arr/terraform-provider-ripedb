// Copyright (c) The RIPE DB Provider for Terraform Authors
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/frederic-arr/ripedb-go/ripedb"
	"github.com/frederic-arr/ripedb-go/ripedb/models"
	"github.com/frederic-arr/rpsl-go"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ObjectResource{}

func NewObjectResource() resource.Resource {
	return &ObjectResource{}
}

type ObjectResource struct {
	client *ripedb.RipeDbClient
}

func (r *ObjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object"
}

func (r *ObjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage a generic object in the RIPE Database.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "the ID of the object",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"class": schema.StringAttribute{
				MarkdownDescription: "the class of the object",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "the value of the class",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"attributes": schema.ListNestedAttribute{
				MarkdownDescription: "the attributes of the object. The first attribute will be used as the object class and key",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "the name of the attribute",
							Required:            true,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "the value of the attribute",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (r *ObjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*RipeDbProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *RipeDbProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = data.Client
}

func (r *ObjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ObjectModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resource := data.Class.ValueString()
	obj := modelToObject(&data)

	// The first attribute should always be the class of the object
	// Since the class is given in a separate field, we have to prepend it
	// We also always add the source as it is already specified in the provider
	obj.Attributes = append([]rpsl.Attribute{{Name: resource, Value: data.Value.ValueString()}}, obj.Attributes...)
	obj.Attributes = append(obj.Attributes, rpsl.Attribute{Name: "source", Value: (*r.client).GetSource()})

	m, err := models.ObjectToModel(resource, *obj)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed to validate object with %s schema", resource), err.Error())
		return
	}

	obj, err = (*r.client).CreateObject(resource, obj)
	if err != nil {
		resp.Diagnostics.AddError("failed to create object in RIPE database", err.Error())
		return
	}

	// Remove the first field and timestamps
	// first field: we already specify its data in the .class and .value fields
	filterObject(obj, &data)
	data.Id = types.StringValue(fmt.Sprintf("%s:%s", resource, m.Key()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ObjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ObjectModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()
	idParts := strings.SplitN(id, ":", 2)
	obj, err := (*r.client).GetObject(idParts[0], idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("failed to query RIPE database", err.Error())
		return
	}

	filterObject(obj, &data)
	data.Class = types.StringValue(idParts[0])
	data.Value = types.StringValue(obj.Attributes[0].Value)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ObjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ObjectModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resource := data.Class.ValueString()
	obj := modelToObject(&data)

	// The first attribute should always be the class of the object
	// Since the class is given in a separate field, we have to prepend it
	// We also always add the source as it is already specified in the provider
	obj.Attributes = append([]rpsl.Attribute{{Name: resource, Value: data.Value.ValueString()}}, obj.Attributes...)
	obj.Attributes = append(obj.Attributes, rpsl.Attribute{Name: "source", Value: (*r.client).GetSource()})

	m, err := models.ObjectToModel(resource, *obj)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed to validate object with %s schema", resource), err.Error())
		return
	}

	obj, err = (*r.client).UpdateObject(resource, m.Key(), obj)
	if err != nil {
		resp.Diagnostics.AddError("failed to update RIPE database object", err.Error())
		return
	}

	// Remove the first field and timestamps
	// first field: we already specify its data in the .class and .value fields
	filterObject(obj, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ObjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ObjectModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()
	idParts := strings.SplitN(id, ":", 2)
	_, err := (*r.client).DeleteObject(idParts[0], idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("failed to delete RIPE database object", err.Error())
		return
	}
}

func (r *ObjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.SplitN(req.ID, ":", 2)

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <class>:<primary_key>. Got: %q", req.ID),
		)
		return
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
