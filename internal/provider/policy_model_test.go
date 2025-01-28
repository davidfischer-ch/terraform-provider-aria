// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

const IDENTIFIER = "7a35997a-cb5c-406b-b1a7-6dc34ac74105"

func TestPolicyModelString(t *testing.T) {
	cases := []struct{ name, typeId, typeName string }{
		{"My Approval policy", "com.vmware.policy.approval", "Approval"},
		{"My Day2 policy", "com.vmware.policy.deployment.action", "Deployment Action"},
		{"My Lease policy", "com.vmware.policy.deployment.lease", "Deployment Lease"},
		{"My Deployment Limit policy", "com.vmware.policy.deployment.limit", "Deployment Limit"},
		{"My Resource Quota policy", "com.vmware.policy.resource.quota", "Resource Quota"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			CheckEqual(
				t,
				PolicyModel{
					Id:     types.StringValue(IDENTIFIER),
					Name:   types.StringValue(tc.name),
					TypeId: types.StringValue(tc.typeId),
				}.String(),
				fmt.Sprintf("%s Policy %s (%s)", tc.typeName, IDENTIFIER, tc.name),
			)
		})
	}
}
