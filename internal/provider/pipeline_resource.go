package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/popsink/terraform-provider-popsink/internal/client"
)

// jsonConnectorTypeValidator validates that connector types in JSON are valid
type jsonConnectorTypeValidator struct{}

// Description returns a description of the validator
func (v jsonConnectorTypeValidator) Description(_ context.Context) string {
	return "validates that source_type and target_type in JSON configuration are valid connector types"
}

// MarkdownDescription returns a markdown description of the validator
func (v jsonConnectorTypeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation
func (v jsonConnectorTypeValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	jsonStr := request.ConfigValue.ValueString()

	var config client.PipelineConfiguration
	if err := json.Unmarshal([]byte(jsonStr), &config); err != nil {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid JSON",
			fmt.Sprintf("Configuration must be valid JSON: %s", err.Error()),
		)
		return
	}

	// Check source_type if provided
	if config.SourceType != nil {
		valid := false
		for _, validType := range validConnectorTypes {
			if *config.SourceType == validType {
				valid = true
				break
			}
		}
		if !valid {
			response.Diagnostics.AddAttributeError(
				request.Path,
				"Invalid Source Type",
				fmt.Sprintf("source_type must be one of: %v, got: %s", validConnectorTypes, *config.SourceType),
			)
		}
	}

	// Check target_type if provided
	if config.TargetType != nil {
		valid := false
		for _, validType := range validConnectorTypes {
			if *config.TargetType == validType {
				valid = true
				break
			}
		}
		if !valid {
			response.Diagnostics.AddAttributeError(
				request.Path,
				"Invalid Target Type",
				fmt.Sprintf("target_type must be one of: %v, got: %s", validConnectorTypes, *config.TargetType),
			)
		}
	}
}

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &pipelineResource{}
	_ resource.ResourceWithConfigure   = &pipelineResource{}
	_ resource.ResourceWithImportState = &pipelineResource{}
)

// NewPipelineResource creates a new pipeline resource
func NewPipelineResource() resource.Resource {
	return &pipelineResource{}
}

// pipelineResource defines the resource implementation
type pipelineResource struct {
	client *client.Client
}

// pipelineResourceModel describes the resource data model
type pipelineResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	TeamID            types.String `tfsdk:"team_id"`
	TeamName          types.String `tfsdk:"team_name"`
	State             types.String `tfsdk:"state"`
	JSONConfiguration types.String `tfsdk:"json_configuration"`
}

// Valid connector types based on the OpenAPI schema
var validConnectorTypes = []string{
	"JOB_SMT",
	"KAFKA_SOURCE",
	"ORACLE_TARGET",
}

// Valid pipeline states
var validPipelineStates = []string{
	"draft",
	"paused",
	"live",
	"error",
	"building",
}

// Metadata returns the resource type name
func (r *pipelineResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pipeline"
}

// Schema defines the resource schema
func (r *pipelineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Popsink pipeline resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the pipeline.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the pipeline.",
				Required:    true,
			},
			"team_id": schema.StringAttribute{
				Description: "The UUID of the team that owns the pipeline.",
				Required:    true,
			},
			"team_name": schema.StringAttribute{
				Description: "The name of the team that owns the pipeline.",
				Computed:    true,
			},
			"state": schema.StringAttribute{
				Description: "The state of the pipeline. Valid values: draft, paused, live, error, building.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(validPipelineStates...),
				},
			},
			"json_configuration": schema.StringAttribute{
				Description: "The complete configuration of the pipeline as a JSON string. " +
					"The source_type and target_type fields must be one of: JOB_SMT, KAFKA_SOURCE, ORACLE_TARGET.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					jsonConnectorTypeValidator{},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *pipelineResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *pipelineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan pipelineResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse JSON configuration
	var config client.PipelineConfiguration
	if err := json.Unmarshal([]byte(plan.JSONConfiguration.ValueString()), &config); err != nil {
		resp.Diagnostics.AddError(
			"Invalid JSON Configuration",
			fmt.Sprintf("Could not parse json_configuration: %s", err.Error()),
		)
		return
	}

	// Create pipeline
	createReq := &client.PipelineCreate{
		Name:              plan.Name.ValueString(),
		TeamID:            plan.TeamID.ValueString(),
		State:             client.PipelineState(plan.State.ValueString()),
		JSONConfiguration: &config,
	}

	pipeline, err := r.client.CreatePipeline(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Pipeline",
			fmt.Sprintf("Could not create pipeline: %s", err.Error()),
		)
		return
	}

	// Update state with created pipeline
	plan.ID = types.StringValue(pipeline.ID)
	plan.Name = types.StringValue(pipeline.Name)
	plan.TeamID = types.StringValue(pipeline.TeamID)
	plan.TeamName = types.StringValue(pipeline.TeamName)
	plan.State = types.StringValue(string(pipeline.State))

	// Keep the original configuration from the plan instead of using API response
	// The API may return a different format or incomplete configuration
	// plan.JSONConfiguration already contains the original configuration

	tflog.Info(ctx, "Created pipeline", map[string]any{"id": pipeline.ID})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the resource state
func (r *pipelineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state pipelineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pipeline, err := r.client.GetPipeline(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Pipeline",
			fmt.Sprintf("Could not read pipeline %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// If pipeline not found, remove from state
	if pipeline == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Update state
	state.Name = types.StringValue(pipeline.Name)
	state.TeamID = types.StringValue(pipeline.TeamID)
	state.TeamName = types.StringValue(pipeline.TeamName)
	state.State = types.StringValue(string(pipeline.State))

	// Keep the original configuration from the state instead of using API response
	// The API may return a different format or incomplete configuration
	// state.JSONConfiguration already contains the current configuration

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource
func (r *pipelineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan pipelineResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state pipelineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse JSON configuration
	var config client.PipelineConfiguration
	if err := json.Unmarshal([]byte(plan.JSONConfiguration.ValueString()), &config); err != nil {
		resp.Diagnostics.AddError(
			"Invalid JSON Configuration",
			fmt.Sprintf("Could not parse json_configuration: %s", err.Error()),
		)
		return
	}

	// Build update request
	updateReq := &client.PipelineUpdate{}

	if !plan.Name.Equal(state.Name) {
		name := plan.Name.ValueString()
		updateReq.Name = &name
	}

	if !plan.TeamID.Equal(state.TeamID) {
		teamID := plan.TeamID.ValueString()
		updateReq.TeamID = &teamID
	}

	if !plan.State.Equal(state.State) {
		pipelineState := client.PipelineState(plan.State.ValueString())
		updateReq.State = &pipelineState
	}

	if !plan.JSONConfiguration.Equal(state.JSONConfiguration) {
		updateReq.JSONConfiguration = &config
	}

	// Update pipeline
	pipeline, err := r.client.UpdatePipeline(ctx, state.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Pipeline",
			fmt.Sprintf("Could not update pipeline %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Update state
	plan.ID = types.StringValue(pipeline.ID)
	plan.Name = types.StringValue(pipeline.Name)
	plan.TeamID = types.StringValue(pipeline.TeamID)
	plan.TeamName = types.StringValue(pipeline.TeamName)
	plan.State = types.StringValue(string(pipeline.State))

	// Keep the original configuration from the plan instead of using API response
	// The API may return a different format or incomplete configuration
	// plan.JSONConfiguration already contains the updated configuration

	tflog.Info(ctx, "Updated pipeline", map[string]any{"id": pipeline.ID})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource
func (r *pipelineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state pipelineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePipeline(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Pipeline",
			fmt.Sprintf("Could not delete pipeline %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	tflog.Info(ctx, "Deleted pipeline", map[string]any{"id": state.ID.ValueString()})
}

// ImportState imports the resource state
func (r *pipelineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
