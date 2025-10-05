package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/popsink/terraform-provider-popsink/internal/client"
)

// normalizeRetentionConfig filters the API response to only include non-empty/non-null fields
// that would be meaningful in the Terraform configuration
func normalizeRetentionConfig(config client.BrokerConfiguration) map[string]any {
	normalized := make(map[string]any)

	// Copy all non-empty fields from the original config
	for key, value := range config {
		// Skip null values and empty strings
		if value == nil || value == "" {
			continue
		}
		normalized[key] = value
	}

	return normalized
}

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &envResource{}
	_ resource.ResourceWithConfigure   = &envResource{}
	_ resource.ResourceWithImportState = &envResource{}
)

// NewEnvResource creates a new environment resource
func NewEnvResource() resource.Resource {
	return &envResource{}
}

// envResource defines the resource implementation
type envResource struct {
	client *client.Client
}

// envResourceModel describes the resource data model
type envResourceModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	UseRetention           types.Bool   `tfsdk:"use_retention"`
	RetentionConfiguration types.String `tfsdk:"retention_configuration"`
}

// Metadata returns the resource type name
func (r *envResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_env"
}

// Schema defines the resource schema
func (r *envResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Popsink environment resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the environment.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the environment.",
				Required:    true,
			},
			"use_retention": schema.BoolAttribute{
				Description: "Whether message retention is enabled for this environment.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"retention_configuration": schema.StringAttribute{
				Description: "Retention policy configuration as a JSON string. Only used when use_retention is true.",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *envResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates the resource
func (r *envResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan envResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create environment
	createReq := &client.EnvCreate{
		Name:         plan.Name.ValueString(),
		UseRetention: plan.UseRetention.ValueBool(),
	}

	// Handle retention configuration if provided
	if !plan.RetentionConfiguration.IsNull() && !plan.RetentionConfiguration.IsUnknown() && plan.RetentionConfiguration.ValueString() != "" {
		var retentionConfig client.BrokerConfiguration
		if err := json.Unmarshal([]byte(plan.RetentionConfiguration.ValueString()), &retentionConfig); err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("retention_configuration"),
				"Invalid JSON",
				fmt.Sprintf("Could not parse retention_configuration as JSON: %s", err.Error()),
			)
			return
		}
		createReq.RetentionConfiguration = &retentionConfig
	}

	env, err := r.client.CreateEnv(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Environment",
			fmt.Sprintf("Could not create environment: %s", err.Error()),
		)
		return
	}

	// Update state with created environment
	plan.ID = types.StringValue(env.ID)
	plan.Name = types.StringValue(env.Name)
	plan.UseRetention = types.BoolValue(env.UseRetention)

	if env.RetentionConfiguration != nil {
		// Normalize the retention configuration to match the planned configuration
		normalizedConfig := normalizeRetentionConfig(*env.RetentionConfiguration)
		retentionJSON, err := json.Marshal(normalizedConfig)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Marshaling Retention Configuration",
				fmt.Sprintf("Could not marshal retention configuration: %s", err.Error()),
			)
			return
		}
		plan.RetentionConfiguration = types.StringValue(string(retentionJSON))
	} else {
		plan.RetentionConfiguration = types.StringNull()
	}

	tflog.Info(ctx, "Created environment", map[string]any{"id": env.ID})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the resource state
func (r *envResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state envResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	env, err := r.client.GetEnv(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Environment",
			fmt.Sprintf("Could not read environment %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// If environment not found, remove from state
	if env == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Update state
	state.Name = types.StringValue(env.Name)
	state.UseRetention = types.BoolValue(env.UseRetention)

	if env.RetentionConfiguration != nil {
		// Normalize the retention configuration to match the planned configuration
		normalizedConfig := normalizeRetentionConfig(*env.RetentionConfiguration)
		retentionJSON, err := json.Marshal(normalizedConfig)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Marshaling Retention Configuration",
				fmt.Sprintf("Could not marshal retention configuration: %s", err.Error()),
			)
			return
		}
		state.RetentionConfiguration = types.StringValue(string(retentionJSON))
	} else {
		state.RetentionConfiguration = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource
func (r *envResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan envResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state envResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update request
	updateReq := &client.EnvUpdate{}

	if !plan.Name.Equal(state.Name) {
		name := plan.Name.ValueString()
		updateReq.Name = &name
	}

	if !plan.UseRetention.Equal(state.UseRetention) {
		useRetention := plan.UseRetention.ValueBool()
		updateReq.UseRetention = &useRetention
	}

	if !plan.RetentionConfiguration.Equal(state.RetentionConfiguration) {
		if !plan.RetentionConfiguration.IsNull() && !plan.RetentionConfiguration.IsUnknown() && plan.RetentionConfiguration.ValueString() != "" {
			var retentionConfig client.BrokerConfiguration
			if err := json.Unmarshal([]byte(plan.RetentionConfiguration.ValueString()), &retentionConfig); err != nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("retention_configuration"),
					"Invalid JSON",
					fmt.Sprintf("Could not parse retention_configuration as JSON: %s", err.Error()),
				)
				return
			}
			updateReq.RetentionConfiguration = &retentionConfig
		} else {
			// Set to nil to remove the retention configuration
			updateReq.RetentionConfiguration = nil
		}
	}

	// Update environment
	env, err := r.client.UpdateEnv(ctx, state.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Environment",
			fmt.Sprintf("Could not update environment %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Update state
	plan.ID = types.StringValue(env.ID)
	plan.Name = types.StringValue(env.Name)
	plan.UseRetention = types.BoolValue(env.UseRetention)

	if env.RetentionConfiguration != nil {
		// Normalize the retention configuration to match the planned configuration
		normalizedConfig := normalizeRetentionConfig(*env.RetentionConfiguration)
		retentionJSON, err := json.Marshal(normalizedConfig)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Marshaling Retention Configuration",
				fmt.Sprintf("Could not marshal retention configuration: %s", err.Error()),
			)
			return
		}
		plan.RetentionConfiguration = types.StringValue(string(retentionJSON))
	} else {
		plan.RetentionConfiguration = types.StringNull()
	}

	tflog.Info(ctx, "Updated environment", map[string]any{"id": env.ID})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource
func (r *envResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state envResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Note: Based on the API specification, there doesn't appear to be a DELETE endpoint for environments
	// For now, we'll just remove it from the Terraform state and log a warning
	resp.Diagnostics.AddWarning(
		"Environment Deletion Not Supported",
		fmt.Sprintf("The Popsink API does not support deleting environments. Environment %s has been removed from Terraform state but still exists in the API. Please manage environment deletion through the Popsink web interface or API directly.", state.ID.ValueString()),
	)

	tflog.Warn(ctx, "Environment removed from state (deletion not supported by API)", map[string]any{"id": state.ID.ValueString()})
}

// ImportState imports the resource state
func (r *envResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
