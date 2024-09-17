package provider

import (
	"cmp"
	"context"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"slices"
)

var _ resource.Resource = &PoolResource{}
var _ resource.ResourceWithImportState = &PoolResource{}

func NewPoolResource() resource.Resource {
	return &PoolResource{}
}

type PoolResource struct {
}

type PoolResourceModel struct {
	Resources           types.Set    `tfsdk:"resources"`
	Borrowers           types.Set    `tfsdk:"borrowers"`
	ResourcesToBorrower types.Map    `tfsdk:"resource_borrowers"`
	BorrowerToResources types.Map    `tfsdk:"borrower_resources"`
	Id                  types.String `tfsdk:"id"`
}

func (r *PoolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_strings"
}

func (r *PoolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A pool of strings that are exclusively checked out on create and checked in on destroy.",

		Attributes: map[string]schema.Attribute{
			"resources": schema.SetAttribute{
				ElementType: basetypes.StringType{},
				Required:    false,
				Sensitive:   false,
				Computed:    false,
				Optional:    true,
				Description: "The resources which are available to be checked out.",
			},
			"borrowers": schema.SetAttribute{
				ElementType: basetypes.StringType{},
				Required:    false,
				Sensitive:   false,
				Computed:    false,
				Optional:    true,
				Description: "The borrowers interested in checking out a resource.",
			},
			"borrower_resources": schema.MapAttribute{
				ElementType: basetypes.StringType{},
				Computed:    true,
				Sensitive:   false,
				Description: "A map from borrower to checked out resource.",
			},
			"resource_borrowers": schema.MapAttribute{
				ElementType: basetypes.StringType{},
				Computed:    true,
				Sensitive:   false,
				Description: "A map from borrower to checked out resource.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The identifier for the resource pool.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *PoolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
}

func asString(value attr.Value) string {
	return value.(types.String).ValueString()
}

func reverseMap(input map[string]attr.Value) map[string]attr.Value {
	result := make(map[string]attr.Value)
	for k, v := range input {
		result[asString(v)] = types.StringValue(k)
	}
	return result
}

func compareAttrs(a attr.Value, b attr.Value) int {
	return cmp.Compare(asString(a), asString(b))
}

func updateMappings(oldMappings map[string]attr.Value, newKeys []attr.Value, newValues []attr.Value) map[string]attr.Value {
	newMappings := make(map[string]attr.Value)
	newValueSet := make(map[string]bool)
	slices.SortFunc(newKeys, compareAttrs)
	slices.SortFunc(newValues, compareAttrs)
	for _, v := range newValues {
		newValueSet[asString(v)] = true
	}
	for _, k := range newKeys {
		oldValue, hadMapping := oldMappings[asString(k)]
		if hadMapping {
			_, isValidValue := newValueSet[asString(oldValue)]
			if isValidValue {
				newMappings[asString(k)] = oldValue
			} else {
				for _, v := range newValues {
					_, isAssigned := newMappings[asString(v)]
					if !isAssigned {
						newMappings[asString(k)] = v
						break
					}
				}
			}
		} else {
			for _, v := range newValues {
				_, isAssigned := newMappings[asString(v)]
				if !isAssigned {
					newMappings[asString(k)] = v
					break
				}
			}
		}
	}
	return newMappings
}

func (r *PoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PoolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	data.Id = types.StringValue(uuid.NewString())
	borrowerToResource := updateMappings(map[string]attr.Value{}, data.Borrowers.Elements(), data.Resources.Elements())
	resourceToBorrower := reverseMap(borrowerToResource)
	data.BorrowerToResources = types.MapValueMust(types.StringType, borrowerToResource)
	data.ResourcesToBorrower = types.MapValueMust(types.StringType, resourceToBorrower)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PoolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var stateData PoolResourceModel
	var planData PoolResourceModel
	var newStateData PoolResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &newStateData)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	borrowerToResource := updateMappings(stateData.BorrowerToResources.Elements(), planData.Borrowers.Elements(), planData.Resources.Elements())
	resourceToBorrower := reverseMap(borrowerToResource)
	newStateData.BorrowerToResources = types.MapValueMust(types.StringType, borrowerToResource)
	newStateData.ResourcesToBorrower = types.MapValueMust(types.StringType, resourceToBorrower)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newStateData)...)
}

func (r *PoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PoolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *PoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
