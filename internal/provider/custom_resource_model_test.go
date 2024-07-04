// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCustomResourceModelToAPI(t *testing.T) {
	ctx := context.Background()
	inputParameters, diags := types.ListValueFrom(ctx, types.StringType, []string{})
	CheckDiagnostics(t, diags, "", "")
	resource := CustomResourceModel{
		Id:           types.StringValue("d4352b6a-84cd-4729-abbf-f3d83e53c46f"),
		DisplayName:  types.StringValue("My Custom Resource"),
		Description:  types.StringValue("Some description\r\n"),
		ResourceType: types.StringValue("Custom.MyCustom"),
		SchemaType:   types.StringValue("ABX_USER_DEFINED"),
		Status:       types.StringValue("RELEASED"),
		Properties: []PropertyModel{
			{
				Name:        types.StringValue("identifier"),
				Title:       types.StringValue("Identifier"),
				Description: types.StringValue("Identify the resource"),
				Type:        types.StringValue("string"),
			},
			{
				Name:    types.StringValue("replicas"),
				Title:   types.StringValue("Replicas"),
				Type:    types.StringValue("integer"),
				Default: types.StringValue("2"),
			},
			{
				Name:    types.StringValue("enabled"),
				Title:   types.StringValue("Enabled"),
				Type:    types.StringValue("boolean"),
				Default: types.StringValue("true"),
			},
		},
		Create: CustomResourceActionModel{
			Id:              types.StringValue("c974e486-9039-4b84-9152-0e5aa2074d26"),
			Type:            types.StringValue("abx.action"),
			ProjectId:       types.StringValue("175bed78-dd9e-4999-8669-cc62388e9abb"),
			InputParameters: inputParameters,
		},
		Read: CustomResourceActionModel{
			Id:              types.StringValue("7d59017f-cf0d-4f74-8aac-ffa351ba54d8"),
			Type:            types.StringValue("abx.action"),
			ProjectId:       types.StringValue("175bed78-dd9e-4999-8669-cc62388e9abb"),
			InputParameters: inputParameters,
		},
		Update: CustomResourceActionModel{
			Id:              types.StringValue("edb1824c-ca47-4df4-8804-4de3c20a28a4"),
			Type:            types.StringValue("abx.action"),
			ProjectId:       types.StringValue("175bed78-dd9e-4999-8669-cc62388e9abb"),
			InputParameters: inputParameters,
		},
		Delete: CustomResourceActionModel{
			Id:              types.StringValue("d40c0ca0-9d65-463e-a4ee-b1d99c0e23a8"),
			Type:            types.StringValue("abx.action"),
			ProjectId:       types.StringValue("175bed78-dd9e-4999-8669-cc62388e9abb"),
			InputParameters: inputParameters,
		},
		ProjectId: types.StringValue("175bed78-dd9e-4999-8669-cc62388e9abb"),
	}
	raw, diags := resource.ToAPI(ctx)
	CheckDiagnostics(t, diags, "", "")
	CheckEqual(t, raw.Id, "d4352b6a-84cd-4729-abbf-f3d83e53c46f")
	CheckEqual(t, raw.DisplayName, "My Custom Resource")
	CheckEqual(t, raw.Description, "Some description\n")
	CheckEqual(t, raw.ResourceType, "Custom.MyCustom")
	CheckEqual(t, raw.SchemaType, "ABX_USER_DEFINED")
	CheckEqual(t, raw.Status, "RELEASED")
	CheckEqual(t, raw.MainActions["create"].Id, "c974e486-9039-4b84-9152-0e5aa2074d26")
	CheckEqual(t, raw.MainActions["create"].Name, "")
	CheckEqual(t, raw.MainActions["create"].Type, "abx.action")
	CheckEqual(t, raw.MainActions["create"].ProjectId, "175bed78-dd9e-4999-8669-cc62388e9abb")
	CheckDeepEqual(t, raw.MainActions["create"].InputParameters, []string{})
	CheckEqual(t, raw.MainActions["read"].Id, "7d59017f-cf0d-4f74-8aac-ffa351ba54d8")
	CheckEqual(t, raw.MainActions["read"].Name, "")
	CheckEqual(t, raw.MainActions["read"].Type, "abx.action")
	CheckEqual(t, raw.MainActions["read"].ProjectId, "175bed78-dd9e-4999-8669-cc62388e9abb")
	CheckDeepEqual(t, raw.MainActions["read"].InputParameters, []string{})
	CheckEqual(t, raw.MainActions["update"].Id, "edb1824c-ca47-4df4-8804-4de3c20a28a4")
	CheckEqual(t, raw.MainActions["update"].Name, "")
	CheckEqual(t, raw.MainActions["update"].Type, "abx.action")
	CheckEqual(t, raw.MainActions["update"].ProjectId, "175bed78-dd9e-4999-8669-cc62388e9abb")
	CheckDeepEqual(t, raw.MainActions["update"].InputParameters, []string{})
	CheckEqual(t, raw.MainActions["delete"].Id, "d40c0ca0-9d65-463e-a4ee-b1d99c0e23a8")
	CheckEqual(t, raw.MainActions["delete"].Name, "")
	CheckEqual(t, raw.MainActions["delete"].Type, "abx.action")
	CheckEqual(t, raw.MainActions["delete"].ProjectId, "175bed78-dd9e-4999-8669-cc62388e9abb")
	CheckDeepEqual(t, raw.MainActions["delete"].InputParameters, []string{})
	CheckEqual(t, raw.ProjectId, "175bed78-dd9e-4999-8669-cc62388e9abb")
	CheckEqual(t, raw.OrgId, "")
}

func TestCustomResourceModelFromAPI(t *testing.T) {
	ctx := context.Background()
	inputParameters, diags := types.ListValueFrom(ctx, types.StringType, []string{})

	CheckDiagnostics(t, diags, "", "")
	raw := CustomResourceAPIModel{
		Id:           "d4352b6a-84cd-4729-abbf-f3d83e53c46f",
		DisplayName:  "My Custom Resource",
		Description:  "Some description\n",
		ResourceType: "Custom.MyCustom",
		SchemaType:   "ABX_USER_DEFINED",
		Status:       "RELEASED",
		Properties: CustomResourcePropertiesAPIModel{
			Properties: map[string]PropertyAPIModel{
				"identifier": {
					Title:       "Identifier",
					Description: "Identify the resource",
					Type:        "string",
				},
				"replicas": {
					Title:   "Replicas",
					Type:    "integer",
					Default: int64(2),
				},
				"enabled": {
					Title:   "Enabled",
					Type:    "boolean",
					Default: true,
				},
			},
		},
		MainActions: map[string]CustomResourceActionAPIModel{
			"create": {
				Id:              "c974e486-9039-4b84-9152-0e5aa2074d26",
				Type:            "abx.action",
				Name:            "SomeCreateFunction",
				ProjectId:       "175bed78-dd9e-4999-8669-cc62388e9abb",
				InputParameters: []string{},
			},
			"read": {
				Id:              "7d59017f-cf0d-4f74-8aac-ffa351ba54d8",
				Type:            "abx.action",
				Name:            "SomeReadFunction",
				ProjectId:       "175bed78-dd9e-4999-8669-cc62388e9abb",
				InputParameters: []string{},
			},
			"update": {
				Id:              "edb1824c-ca47-4df4-8804-4de3c20a28a4",
				Type:            "abx.action",
				Name:            "SomeUpdateFunction",
				ProjectId:       "175bed78-dd9e-4999-8669-cc62388e9abb",
				InputParameters: []string{},
			},
			"delete": {
				Id:              "d40c0ca0-9d65-463e-a4ee-b1d99c0e23a8",
				Type:            "abx.action",
				Name:            "SomeDeleteFunction",
				ProjectId:       "175bed78-dd9e-4999-8669-cc62388e9abb",
				InputParameters: []string{},
			},
		},
		ProjectId: "175bed78-dd9e-4999-8669-cc62388e9abb",
		OrgId:     "f57768e3-6710-4864-982b-68456c8ea29a",
	}
	resource := CustomResourceModel{}
	diags = resource.FromAPI(ctx, raw)
	CheckDiagnostics(t, diags, "", "")
	CheckEqual(t, resource.Id.ValueString(), "d4352b6a-84cd-4729-abbf-f3d83e53c46f")
	CheckEqual(t, resource.DisplayName.ValueString(), "My Custom Resource")
	CheckEqual(t, resource.Description.ValueString(), "Some description\n")
	CheckEqual(t, resource.ResourceType.ValueString(), "Custom.MyCustom")
	CheckEqual(t, resource.SchemaType.ValueString(), "ABX_USER_DEFINED")
	CheckEqual(t, resource.Status.ValueString(), "RELEASED")
	CheckEqual(t, resource.Create.Id.ValueString(), "c974e486-9039-4b84-9152-0e5aa2074d26")
	CheckEqual(t, resource.Create.Name.ValueString(), "SomeCreateFunction")
	CheckEqual(t, resource.Create.Type.ValueString(), "abx.action")
	CheckEqual(t, resource.Create.ProjectId.ValueString(), "175bed78-dd9e-4999-8669-cc62388e9abb")
	CheckDeepEqual(t, resource.Create.InputParameters, inputParameters)
	CheckEqual(t, resource.Read.Id.ValueString(), "7d59017f-cf0d-4f74-8aac-ffa351ba54d8")
	CheckEqual(t, resource.Read.Name.ValueString(), "SomeReadFunction")
	CheckEqual(t, resource.Read.Type.ValueString(), "abx.action")
	CheckEqual(t, resource.Read.ProjectId.ValueString(), "175bed78-dd9e-4999-8669-cc62388e9abb")
	CheckDeepEqual(t, resource.Read.InputParameters, inputParameters)
	CheckEqual(t, resource.Update.Id.ValueString(), "edb1824c-ca47-4df4-8804-4de3c20a28a4")
	CheckEqual(t, resource.Update.Name.ValueString(), "SomeUpdateFunction")
	CheckEqual(t, resource.Update.Type.ValueString(), "abx.action")
	CheckEqual(t, resource.Update.ProjectId.ValueString(), "175bed78-dd9e-4999-8669-cc62388e9abb")
	CheckDeepEqual(t, resource.Update.InputParameters, inputParameters)
	CheckEqual(t, resource.Delete.Id.ValueString(), "d40c0ca0-9d65-463e-a4ee-b1d99c0e23a8")
	CheckEqual(t, resource.Delete.Name.ValueString(), "SomeDeleteFunction")
	CheckEqual(t, resource.Delete.Type.ValueString(), "abx.action")
	CheckEqual(t, resource.Delete.ProjectId.ValueString(), "175bed78-dd9e-4999-8669-cc62388e9abb")
	CheckDeepEqual(t, resource.Delete.InputParameters, inputParameters)
	CheckEqual(t, resource.ProjectId.ValueString(), "175bed78-dd9e-4999-8669-cc62388e9abb")
	CheckEqual(t, resource.OrgId.ValueString(), "f57768e3-6710-4864-982b-68456c8ea29a")
}
