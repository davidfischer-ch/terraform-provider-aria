// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrchestratorWorkflowCompleteExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
resource "aria_orchestrator_category" "root" {
  name      = "TEST_ARIA_PROVIDER"
  type      = "WorkflowCategory"
  parent_id = ""
}

locals {
  workflow_yml = <<EOT
inputForms:
  - schema:
      deploymentID:
        id: deploymentID
        type:
          dataType: string
        label: deploymentID
        constraints:
          required: false
      vRAHost:
        id: vRAHost
        type:
          dataType: reference
          referenceType: 'VRA:Host'
        label: vRAHost
    layout:
      pages:
        - id: page_dnpofm0t
          sections:
            - id: section_owly0g7u
              fields:
                - id: deploymentID
                  display: textField
                  signpostPosition: right-middle
                  state:
                    visible: true
                    read-only: false
            - id: section_qjxtzpr8
              fields:
                - id: vRAHost
                  display: valuePickerTree
                  signpostPosition: right-middle
                  state:
                    visible: true
                    read-only: false
          title: General
    options:
      externalValidations: []
    itemId: ''
workflowSchema:
  display-name: Delete Deployment
  position:
    'y': 50
    x: 20
  input:
    param:
      - name: deploymentID
        description: TODO
        type: string
      - name: vRAHost
        description: TODO
        type: 'VRA:Host'
  output:
    param:
      - name: errorCode
        description: TODO
        type: string
  attrib:
    - value:
        number:
          value: 20
      type: number
      name: sleepTime
    - value:
        string:
          value: ''
      type: string
      name: pathUriDelete
    - type: number
      name: statusCode
    - value:
        string:
          value: ''
      type: string
      name: statusMessage
    - type: Array/string
      name: headers
    - value:
        properties:
          property:
            - key: Accept
              value:
                string:
                  value: application/json
            - key: Content-Type
              value:
                string:
                  value: application/json
      type: Properties
      name: inputHeaders
    - value:
        string:
          value: ''
      type: string
      name: pathUriGet
    - value:
        string:
          value: ''
      type: string
      name: contentAsStringGet
    - value:
        string:
          value: ''
      type: string
      name: contentAsStringDelete
    - value:
        string:
          value: ''
      type: string
      name: requestStatus
    - type: number
      name: count
  workflow-item:
    - in-binding: {}
      position:
        'y': 50
        x: 1020
      name: item0
      type: end
      end-mode: '0'
      comparator: 0
    - display-name: Sleep
      script:
        value: "//Auto-generated script\nif ( sleepTime !== null )  {\n\tSystem.sleep(sleepTime * 1000);\n}else  {\n\tthrow \"'sleepTime' is NULL\"; \n}"
        encoded: false
      in-binding:
        bind:
          - description: Time to sleep in seconds
            name: sleepTime
            type: number
            export-name: sleepTime
      out-binding: {}
      description: Sleep a given number of seconds.
      position:
        'y': 60
        x: 220
      name: item2
      out-name: item6
      type: task
      prototype-id: sleep
      content-mode: x
      comparator: 0
    - display-name: Delete operation
      in-binding:
        bind:
          - description: vRA/C host
            name: host
            type: 'VRA:Host'
            export-name: vRAHost
          - description: Resource path Uri
            name: pathUri
            type: string
            export-name: pathUriDelete
          - description: Request headers
            name: inputHeaders
            type: Properties
            export-name: inputHeaders
      out-binding:
        bind:
          - description: Response status code
            name: statusCode
            type: number
            export-name: statusCode
          - description: Response content
            name: contentAsString
            type: string
            export-name: contentAsStringDelete
          - description: Response status message
            name: statusMessage
            type: string
            export-name: statusMessage
          - description: Response headers
            name: headers
            type: Array/string
            export-name: headers
      description: ' '
      position:
        'y': 60
        x: 420
      name: item5
      out-name: item8
      throw-bind-name: errorCode
      type: link
      linked-workflow-id: 06998d5f-27ea-4c02-b62f-e45140a8072c
      comparator: 0
    - display-name: retrieve PathURIDelete
      script:
        value: var pathUriDelete = "/deployment/api/deployments/"+deploymentID
        encoded: false
      in-binding:
        bind:
          - name: deploymentID
            type: string
            export-name: deploymentID
      out-binding:
        bind:
          - name: pathUriDelete
            type: string
            export-name: pathUriDelete
      description: Simple task with custom script capability.
      position:
        'y': 60
        x: 320
      name: item6
      out-name: item5
      catch-name: item20
      throw-bind-name: errorCode
      type: task
      comparator: 0
    - display-name: retrieve RequestID
      script:
        value: |-
          var contentObj = JSON.parse(contentAsStringDelete);
          var requestID = contentObj.id;
          System.log("Request ID : "+requestID)

          var pathUriGet = "/deployment/api/requests/"+requestID
        encoded: false
      in-binding:
        bind:
          - name: contentAsStringDelete
            type: string
            export-name: contentAsStringDelete
      out-binding:
        bind:
          - name: pathUriGet
            type: string
            export-name: pathUriGet
      description: Simple task with custom script capability.
      position:
        'y': 60
        x: 540
      name: item8
      out-name: item10
      type: task
      comparator: 0
    - display-name: Get operation
      in-binding:
        bind:
          - description: vRA/C host
            name: host
            type: 'VRA:Host'
            export-name: vRAHost
          - description: Resource path uri
            name: pathUri
            type: string
            export-name: pathUriGet
          - description: Request headers
            name: inputHeaders
            type: Properties
            export-name: inputHeaders
      out-binding:
        bind:
          - description: Response status code (HTTP standard - 200/400 etc.)
            name: statusCode
            type: number
            export-name: statusCode
          - description: Response content
            name: contentAsString
            type: string
            export-name: contentAsStringGet
          - description: Response status message
            name: statusMessage
            type: string
            export-name: statusMessage
          - description: Response headers
            name: headers
            type: Array/string
            export-name: headers
      description: ' '
      position:
        'y': 60
        x: 660
      name: item10
      out-name: item16
      type: link
      linked-workflow-id: 36cb38fa-1901-4b4d-840a-33f6368757ea
      comparator: 0
    - display-name: Deleted ?
      script:
        value: |-
          if (requestStatus == "SUCCESSFUL"){
              System.log("Deployment Deleted")
              return true;
          }else {
              return false;
              }
        encoded: false
      in-binding:
        bind:
          - name: requestStatus
            type: string
            export-name: requestStatus
      out-binding: {}
      description: Custom decision based on a custom script.
      position:
        'y': 50
        x: 880
      name: item12
      out-name: item0
      alt-out-name: item15
      type: custom-condition
      comparator: 0
    - display-name: Failed ?
      script:
        value: |-
          if (requestStatus == "FAILED"){
              var contentObj = JSON.parse(contentAsStringGet);
              var failDetails = contentObj.details;
              var errorcode = failDetails
              throw " Deployment not Deleted : "+failDetails
              return false;
          } else {
              return true;
          }
        encoded: false
      in-binding:
        bind:
          - name: requestStatus
            type: string
            export-name: requestStatus
          - name: contentAsStringGet
            type: string
            export-name: contentAsStringGet
      out-binding: {}
      description: Custom decision based on a custom script.
      position:
        'y': 140
        x: 780
      name: item15
      out-name: item10
      alt-out-name: item17
      type: custom-condition
      comparator: 0
    - display-name: Retrieve request Status
      script:
        value: |-
          var contentObj = JSON.parse(contentAsStringGet);
          var requestStatus = contentObj.status;
          System.log("Request Status : "+requestStatus)
        encoded: false
      in-binding:
        bind:
          - name: contentAsStringGet
            type: string
            export-name: contentAsStringGet
      out-binding:
        bind:
          - name: requestStatus
            type: string
            export-name: requestStatus
      description: Simple task with custom script capability.
      position:
        'y': 60
        x: 780
      name: item16
      out-name: item12
      type: task
      comparator: 0
    - in-binding: {}
      position:
        'y': 200
        x: 820
      name: item17
      throw-bind-name: errorCode
      type: end
      end-mode: '1'
      comparator: 0
    - display-name: Retrying ?
      script:
        value: |-

          var count = count + 1 ;
          if (count < 5 ){
              System.log("Faled to delete Deployment - Retrying")
              return true ;
          } else {
               System.log("Faled to delete Deployment")
              return false ;
          }
        encoded: false
      in-binding:
        bind:
          - name: count
            type: number
            export-name: count
      out-binding: {}
      description: Custom decision based on a custom script.
      position:
        'y': 130
        x: 280
      name: item20
      out-name: item2
      alt-out-name: item21
      type: custom-condition
      comparator: 0
    - in-binding: {}
      position:
        'y': 200
        x: 320
      name: item21
      throw-bind-name: errorCode
      type: end
      end-mode: '1'
      comparator: 0
    - display-name: Init variables
      script:
        value: var count = 0
        encoded: false
      in-binding: {}
      out-binding:
        bind:
          - name: count
            type: number
            export-name: count
      description: Simple task with custom script capability.
      position:
        'y': 60
        x: 100
      name: item22
      out-name: item2
      type: task
      comparator: 0
  presentation: {}
  root-name: item22
  object-name: 'workflow:name=generic'
  id: c9868f6f-3c0c-41dd-b6b1-040061e8a239
  version: 0.0.7
  api-version: 6.0.0
  allowed-operations: vfe
  restartMode: 1
  resumeFromFailedMode: 0
  editor-version: '2.0'
EOT

  workflow_data = yamldecode(local.workflow_yml)
}

resource "aria_orchestrator_workflow" "test" {
  name        = local.workflow_data["workflowSchema"]["display-name"]
  description = "Workflow generated by the acceptance tests of Aria provider."
  category_id = aria_orchestrator_category.root.id
  version     = "0.0.0"

  allowed_operations      = local.workflow_data["workflowSchema"]["allowed-operations"]
  attrib                  = jsonencode(local.workflow_data["workflowSchema"]["attrib"])
  object_name             = local.workflow_data["workflowSchema"]["object-name"]
  position                = local.workflow_data["workflowSchema"]["position"]
  presentation            = jsonencode(local.workflow_data["workflowSchema"]["presentation"])
  restart_mode            = local.workflow_data["workflowSchema"]["restartMode"]
  resume_from_failed_mode = local.workflow_data["workflowSchema"]["resumeFromFailedMode"]
  root_name               = local.workflow_data["workflowSchema"]["root-name"]
  workflow_item           = jsonencode(local.workflow_data["workflowSchema"]["workflow-item"])
  input_parameters        = local.workflow_data["workflowSchema"]["input"]["param"]
  output_parameters       = local.workflow_data["workflowSchema"]["output"]["param"]
  input_forms             = jsonencode(local.workflow_data["inputForms"])
  api_version             = local.workflow_data["workflowSchema"]["api-version"]
  editor_version          = local.workflow_data["workflowSchema"]["editor-version"]

  force_delete    = true
  wait_on_catalog = false # Make tests faster

  lifecycle {
    postcondition {
      condition     = self.attrib == jsonencode(local.workflow_data["workflowSchema"]["attrib"])
      error_message = "Attrib is not what we expect!"
    }
    postcondition {
      condition     = self.workflow_item == jsonencode(local.workflow_data["workflowSchema"]["workflow-item"])
      error_message = "Workflow Item is not what we expect!"
    }
    postcondition {
      condition     = self.input_forms == jsonencode(local.workflow_data["inputForms"])
      error_message = "Input Forms is not what we expect!"
    }
  }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_orchestrator_workflow.test", "id"),
					resource.TestCheckResourceAttr("aria_orchestrator_workflow.test", "name", "Delete Deployment"),
					resource.TestCheckResourceAttr("aria_orchestrator_workflow.test", "description", "Workflow generated by the acceptance tests of Aria provider."),
					resource.TestCheckResourceAttr("aria_orchestrator_workflow.test", "version", "0.0.0"),
					resource.TestMatchResourceAttr("aria_orchestrator_workflow.test", "version_id", regexp.MustCompile("[0-9a-f]{40}")),
					resource.TestCheckResourceAttr("aria_orchestrator_workflow.test", "position.x", "20"),
					resource.TestCheckResourceAttr("aria_orchestrator_workflow.test", "position.y", "50"),
					resource.TestCheckResourceAttr("aria_orchestrator_workflow.test", "allowed_operations", "vfe"),
					resource.TestCheckResourceAttr("aria_orchestrator_workflow.test", "object_name", "workflow:name=generic"),
					resource.TestCheckResourceAttr("aria_orchestrator_workflow.test", "presentation", "{}"),
					resource.TestCheckResourceAttr("aria_orchestrator_workflow.test", "restart_mode", "1"),
					resource.TestCheckResourceAttr("aria_orchestrator_workflow.test", "resume_from_failed_mode", "0"),
					resource.TestCheckResourceAttr("aria_orchestrator_workflow.test", "root_name", "item22"),
					resource.TestCheckResourceAttr("aria_orchestrator_workflow.test", "api_version", "6.0.0"),
					resource.TestCheckResourceAttr("aria_orchestrator_workflow.test", "editor_version", "2.0"),
					resource.TestCheckResourceAttr("aria_orchestrator_workflow.test", "force_delete", "true"),
					resource.TestCheckResourceAttr("aria_orchestrator_workflow.test", "wait_on_catalog", "false"),
					resource.TestCheckResourceAttrPair(
						"aria_orchestrator_workflow.test", "category_id",
						"aria_orchestrator_category.root", "id",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
