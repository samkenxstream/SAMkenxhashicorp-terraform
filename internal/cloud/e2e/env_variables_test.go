package main

import (
	"context"
	"fmt"
	"testing"
)

func Test_cloud_organization_env_var(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	org, cleanup := createOrganization(t)
	t.Cleanup(cleanup)

	cases := testCases{
		"with TF_ORGANIZATION set": {
			operations: []operationSets{
				{
					prep: func(t *testing.T, orgName, dir string) {
						remoteWorkspace := "cloud-workspace"
						tfBlock := terraformConfigCloudBackendOmitOrg(remoteWorkspace)
						writeMainTF(t, tfBlock, dir)
					},
					commands: []tfCommand{
						{
							command:           []string{"init"},
							expectedCmdOutput: `Terraform Cloud has been successfully initialized!`,
						},
						{
							command:         []string{"apply", "-auto-approve"},
							postInputOutput: []string{`Apply complete!`},
						},
					},
				},
			},
			validations: func(t *testing.T, orgName string) {
				expectedName := "cloud-workspace"
				ws, err := tfeClient.Workspaces.Read(ctx, org.Name, expectedName)
				if err != nil {
					t.Fatal(err)
				}
				if ws == nil {
					t.Fatalf("Expected workspace %s to be present, but is not.", expectedName)
				}
			},
		},
	}

	testRunner(t, cases, 0, fmt.Sprintf("TF_ORGANIZATION=%s", org.Name))

}

func Test_cloud_workspace_env_var(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cases := testCases{
		"with TF_WORKSPACE set": {
			operations: []operationSets{
				{
					prep: func(t *testing.T, orgName, dir string) {
						tfBlock := terraformConfigCloudBackendOmitWorkspaces(orgName)
						writeMainTF(t, tfBlock, dir)
					},
					commands: []tfCommand{
						{
							command:           []string{"init"},
							expectedCmdOutput: `Terraform Cloud has been successfully initialized!`,
						},
						{
							command:         []string{"apply", "-auto-approve"},
							postInputOutput: []string{`Apply complete!`},
						},
					},
				},
				{
					prep: func(t *testing.T, orgName, dir string) {
						tfBlock := terraformConfigCloudBackendOmitWorkspaces(orgName)
						writeMainTF(t, tfBlock, dir)
					},
					commands: []tfCommand{
						{
							command:           []string{"init"},
							expectedCmdOutput: `Terraform Cloud has been successfully initialized!`,
						},
						{
							command:           []string{"workspace", "show"},
							expectedCmdOutput: "wkspace",
						},
					},
				},
			},
			validations: func(t *testing.T, orgName string) {
				expectedName := "wkspace"
				ws, err := tfeClient.Workspaces.Read(ctx, orgName, expectedName)
				if err != nil {
					t.Fatal(err)
				}
				if ws == nil {
					t.Fatalf("Expected workspace %s to be present, but is not.", expectedName)
				}
			},
		},
	}

	testRunner(t, cases, 1, `TF_WORKSPACE=wkspace`)
}

func Test_cloud_null_config(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	org, cleanup := createOrganization(t)
	t.Cleanup(cleanup)

	cases := testCases{
		"with all env vars set": {
			operations: []operationSets{
				{
					prep: func(t *testing.T, orgName, dir string) {
						tfBlock := terraformConfigCloudBackendOmitConfig()
						writeMainTF(t, tfBlock, dir)
					},
					commands: []tfCommand{
						{
							command:           []string{"init"},
							expectedCmdOutput: `Terraform Cloud has been successfully initialized!`,
						},
						{
							command:         []string{"apply", "-auto-approve"},
							postInputOutput: []string{`Apply complete!`},
						},
					},
				},
				{
					prep: func(t *testing.T, orgName, dir string) {
						tfBlock := terraformConfigCloudBackendOmitConfig()
						writeMainTF(t, tfBlock, dir)
					},
					commands: []tfCommand{
						{
							command:           []string{"init"},
							expectedCmdOutput: `Terraform Cloud has been successfully initialized!`,
						},
						{
							command:           []string{"workspace", "show"},
							expectedCmdOutput: "cloud-workspace",
						},
					},
				},
			},
			validations: func(t *testing.T, orgName string) {
				expectedName := "cloud-workspace"
				ws, err := tfeClient.Workspaces.Read(ctx, org.Name, expectedName)
				if err != nil {
					t.Fatal(err)
				}
				if ws == nil {
					t.Fatalf("Expected workspace %s to be present, but is not.", expectedName)
				}
			},
		},
	}

	testRunner(t, cases, 1,
		fmt.Sprintf(`TF_ORGANIZATION=%s`, org.Name),
		fmt.Sprintf(`TF_HOSTNAME=%s`, tfeHostname),
		`TF_WORKSPACE=cloud-workspace`)
}
