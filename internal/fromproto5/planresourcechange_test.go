package fromproto5_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fromproto5"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestPlanResourceChangeRequest(t *testing.T) {
	t.Parallel()

	testProto5Type := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"test_attribute": tftypes.String,
		},
	}

	testProto5Value := tftypes.NewValue(testProto5Type, map[string]tftypes.Value{
		"test_attribute": tftypes.NewValue(tftypes.String, "test-value"),
	})

	testProto5DynamicValue, err := tfprotov5.NewDynamicValue(testProto5Type, testProto5Value)

	if err != nil {
		t.Fatalf("unexpected error calling tfprotov5.NewDynamicValue(): %s", err)
	}

	testFwSchema := &tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"test_attribute": {
				Required: true,
				Type:     types.StringType,
			},
		},
	}

	testCases := map[string]struct {
		input               *tfprotov5.PlanResourceChangeRequest
		resourceSchema      *tfsdk.Schema
		resourceType        tfsdk.ResourceType
		providerMetaSchema  *tfsdk.Schema
		expected            *fwserver.PlanResourceChangeRequest
		expectedDiagnostics diag.Diagnostics
	}{
		"nil": {
			input:    nil,
			expected: nil,
		},
		"empty": {
			input:    &tfprotov5.PlanResourceChangeRequest{},
			expected: nil,
			expectedDiagnostics: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Missing Resource Schema",
					"An unexpected error was encountered when handling the request. "+
						"This is always an issue in the Terraform Provider SDK used to implement the provider and should be reported to the provider developers.\n\n"+
						"Please report this to the provider developer:\n\n"+
						"Missing schema.",
				),
			},
		},
		"config-missing-schema": {
			input: &tfprotov5.PlanResourceChangeRequest{
				Config: &testProto5DynamicValue,
			},
			expected: nil,
			expectedDiagnostics: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Missing Resource Schema",
					"An unexpected error was encountered when handling the request. "+
						"This is always an issue in the Terraform Provider SDK used to implement the provider and should be reported to the provider developers.\n\n"+
						"Please report this to the provider developer:\n\n"+
						"Missing schema.",
				),
			},
		},
		"config": {
			input: &tfprotov5.PlanResourceChangeRequest{
				Config: &testProto5DynamicValue,
			},
			resourceSchema: testFwSchema,
			expected: &fwserver.PlanResourceChangeRequest{
				Config: &tfsdk.Config{
					Raw:    testProto5Value,
					Schema: *testFwSchema,
				},
				ResourceSchema: *testFwSchema,
			},
		},
		"priorprivate": {
			input: &tfprotov5.PlanResourceChangeRequest{
				PriorPrivate: []byte("{}"),
			},
			resourceSchema: testFwSchema,
			expected: &fwserver.PlanResourceChangeRequest{
				PriorPrivate:   []byte("{}"),
				ResourceSchema: *testFwSchema,
			},
		},
		"priorstate-missing-schema": {
			input: &tfprotov5.PlanResourceChangeRequest{
				PriorState: &testProto5DynamicValue,
			},
			expected: nil,
			expectedDiagnostics: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Missing Resource Schema",
					"An unexpected error was encountered when handling the request. "+
						"This is always an issue in the Terraform Provider SDK used to implement the provider and should be reported to the provider developers.\n\n"+
						"Please report this to the provider developer:\n\n"+
						"Missing schema.",
				),
			},
		},
		"priorstate": {
			input: &tfprotov5.PlanResourceChangeRequest{
				PriorState: &testProto5DynamicValue,
			},
			resourceSchema: testFwSchema,
			expected: &fwserver.PlanResourceChangeRequest{
				PriorState: &tfsdk.State{
					Raw:    testProto5Value,
					Schema: *testFwSchema,
				},
				ResourceSchema: *testFwSchema,
			},
		},
		"proposednewstate-missing-schema": {
			input: &tfprotov5.PlanResourceChangeRequest{
				ProposedNewState: &testProto5DynamicValue,
			},
			expected: nil,
			expectedDiagnostics: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Missing Resource Schema",
					"An unexpected error was encountered when handling the request. "+
						"This is always an issue in the Terraform Provider SDK used to implement the provider and should be reported to the provider developers.\n\n"+
						"Please report this to the provider developer:\n\n"+
						"Missing schema.",
				),
			},
		},
		"proposednewstate": {
			input: &tfprotov5.PlanResourceChangeRequest{
				ProposedNewState: &testProto5DynamicValue,
			},
			resourceSchema: testFwSchema,
			expected: &fwserver.PlanResourceChangeRequest{
				ProposedNewState: &tfsdk.Plan{
					Raw:    testProto5Value,
					Schema: *testFwSchema,
				},
				ResourceSchema: *testFwSchema,
			},
		},
		"providermeta-missing-data": {
			input:              &tfprotov5.PlanResourceChangeRequest{},
			resourceSchema:     testFwSchema,
			providerMetaSchema: testFwSchema,
			expected: &fwserver.PlanResourceChangeRequest{
				ProviderMeta: &tfsdk.Config{
					Raw:    tftypes.NewValue(testProto5Type, nil),
					Schema: *testFwSchema,
				},
				ResourceSchema: *testFwSchema,
			},
		},
		"providermeta-missing-schema": {
			input: &tfprotov5.PlanResourceChangeRequest{
				ProviderMeta: &testProto5DynamicValue,
			},
			resourceSchema: testFwSchema,
			expected: &fwserver.PlanResourceChangeRequest{
				// This intentionally should not include ProviderMeta
				ResourceSchema: *testFwSchema,
			},
		},
		"providermeta": {
			input: &tfprotov5.PlanResourceChangeRequest{
				ProviderMeta: &testProto5DynamicValue,
			},
			resourceSchema:     testFwSchema,
			providerMetaSchema: testFwSchema,
			expected: &fwserver.PlanResourceChangeRequest{
				ProviderMeta: &tfsdk.Config{
					Raw:    testProto5Value,
					Schema: *testFwSchema,
				},
				ResourceSchema: *testFwSchema,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, diags := fromproto5.PlanResourceChangeRequest(context.Background(), testCase.input, testCase.resourceType, testCase.resourceSchema, testCase.providerMetaSchema)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
