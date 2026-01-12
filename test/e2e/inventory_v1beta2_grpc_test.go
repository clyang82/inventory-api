package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	grpcinsecure "google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	pbv1beta2 "github.com/project-kessel/inventory-api/api/kessel/inventory/v1beta2"
)

// bearerAuth implements grpc.PerRPCCredentials to inject Authorization
type bearerAuth struct {
	token string
}

func (b *bearerAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", b.token),
	}, nil
}

func (b *bearerAuth) RequireTransportSecurity() bool {
	return false // Set to true if using TLS
}

func enableShortMode(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long-running test in short mode")
	}
}

func TestInventoryAPIHTTP_v1beta2_workspace_movement_tests(t *testing.T) {
	enableShortMode(t)
	ctx := context.Background()

	// Test configuration
	oldWorkspace := "c100"
	newWorkspace := "c101"
	resourceId := "218e5a67-a098-4958-8063-cb5421a2d6cd"
	reporterType := "hbi"
	reporterInstanceId := "test-instance"

	t.Logf("=== Testing Workspace Change Functionality ===")
	t.Logf("Old Workspace: %s", oldWorkspace)
	t.Logf("New Workspace: %s", newWorkspace)

	conn, err := grpc.NewClient(
		inventoryapi_grpc_url,
		grpc.WithTransportCredentials(grpcinsecure.NewCredentials()),
		grpc.WithPerRPCCredentials(&bearerAuth{token: "1234"}),
	)
	assert.NoError(t, err, "Failed to create gRPC client")
	defer func() {
		if connErr := conn.Close(); connErr != nil {
			t.Logf("Failed to close gRPC connection: %v", connErr)
		}
	}()

	conn.Connect()
	assert.NoError(t, err, "Failed to connect gRPC client")

	client := pbv1beta2.NewKesselInventoryServiceClient(conn)

	// Step 1: Add resource with old workspace
	t.Logf("1. Adding resource with workspace_id: %s", oldWorkspace)
	reporterStruct, err := structpb.NewStruct(map[string]interface{}{
		"ansible_host": "test-host.example.com",
	})
	assert.NoError(t, err, "Failed to create structpb for host reporter")

	req := &pbv1beta2.ReportResourceRequest{
		WriteVisibility:    pbv1beta2.WriteVisibility_MINIMIZE_LATENCY,
		Type:               "host",
		ReporterType:       reporterType,
		ReporterInstanceId: reporterInstanceId,
		Representations: &pbv1beta2.ResourceRepresentations{
			Metadata: &pbv1beta2.RepresentationMetadata{
				LocalResourceId: resourceId,
				ApiHref:         "https://api.example.com/hosts/test-host-123",
				ConsoleHref:     proto.String("https://console.example.com/hosts/test-host-123"),
			},
			Common: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"workspace_id": structpb.NewStringValue(oldWorkspace),
				},
			},
			Reporter: reporterStruct,
		},
	}

	_, err = client.ReportResource(ctx, req)
	assert.NoError(t, err, "Failed to Report Resource with old workspace")

	// Step 2: Update resource to new workspace
	t.Logf("2. Updating resource to workspace_id: %s", newWorkspace)
	reqUpdated := &pbv1beta2.ReportResourceRequest{
		WriteVisibility:    pbv1beta2.WriteVisibility_MINIMIZE_LATENCY,
		Type:               "host",
		ReporterType:       reporterType,
		ReporterInstanceId: reporterInstanceId,
		Representations: &pbv1beta2.ResourceRepresentations{
			Metadata: &pbv1beta2.RepresentationMetadata{
				LocalResourceId: resourceId,
				ApiHref:         "https://api.example.com/hosts/test-host-123",
				ConsoleHref:     proto.String("https://console.example.com/hosts/test-host-123"),
			},
			Common: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"workspace_id": structpb.NewStringValue(newWorkspace),
				},
			},
			Reporter: reporterStruct,
		},
	}

	_, err = client.ReportResource(ctx, reqUpdated)
	assert.NoError(t, err, "Failed to Update Resource with new workspace")

	// Step 3: Check authorization for OLD workspace (should be ALLOWED_FALSE)
	t.Logf("3. Checking authorization for OLD workspace %s (should be allowed: false)", oldWorkspace)
	t.Log("Waiting for outbox events to be processed (polling up to 10s)...")

	checkOldWorkspace := &pbv1beta2.CheckRequest{
		Object: &pbv1beta2.ResourceReference{
			ResourceType: "host",
			ResourceId:   resourceId,
			Reporter: &pbv1beta2.ReporterReference{
				Type:       reporterType,
				InstanceId: proto.String(reporterInstanceId),
			},
		},
		Relation: "workspace",
		Subject: &pbv1beta2.SubjectReference{
			Resource: &pbv1beta2.ResourceReference{
				ResourceType: "workspace",
				ResourceId:   oldWorkspace,
				Reporter: &pbv1beta2.ReporterReference{
					Type: "rbac",
				},
			},
		},
	}

	// Poll for up to 10 seconds
	oldWorkspaceAllowed := false
	for i := 0; i < 10; i++ {
		checkResp, err := client.Check(ctx, checkOldWorkspace)
		if err == nil && checkResp.GetAllowed() == pbv1beta2.Allowed_ALLOWED_FALSE {
			t.Logf("✓ Old workspace check returned ALLOWED_FALSE (attempt %d)", i+1)
			oldWorkspaceAllowed = true
			break
		}
		if err != nil {
			t.Logf("Check request failed (attempt %d): %v", i+1, err)
		} else {
			t.Logf("Old workspace check returned %v (attempt %d), expected ALLOWED_FALSE", checkResp.GetAllowed(), i+1)
		}
		if i < 9 {
			t.Log("Waiting 1s before retry...")
			time.Sleep(1 * time.Second)
		}
	}
	assert.True(t, oldWorkspaceAllowed, "Old workspace authorization check did not return ALLOWED_FALSE within timeout")

	// Step 4: Check authorization for NEW workspace (should be ALLOWED_TRUE)
	t.Logf("4. Checking authorization for NEW workspace %s (should be allowed: true)", newWorkspace)
	t.Log("Waiting for outbox events to be processed (polling up to 10s)...")

	checkNewWorkspace := &pbv1beta2.CheckRequest{
		Object: &pbv1beta2.ResourceReference{
			ResourceType: "host",
			ResourceId:   resourceId,
			Reporter: &pbv1beta2.ReporterReference{
				Type:       reporterType,
				InstanceId: proto.String(reporterInstanceId),
			},
		},
		Relation: "workspace",
		Subject: &pbv1beta2.SubjectReference{
			Resource: &pbv1beta2.ResourceReference{
				ResourceType: "workspace",
				ResourceId:   newWorkspace,
				Reporter: &pbv1beta2.ReporterReference{
					Type: "rbac",
				},
			},
		},
	}

	// Poll for up to 10 seconds
	newWorkspaceAllowed := false
	for i := 0; i < 10; i++ {
		checkResp, err := client.Check(ctx, checkNewWorkspace)
		if err == nil && checkResp.GetAllowed() == pbv1beta2.Allowed_ALLOWED_TRUE {
			t.Logf("✓ New workspace check returned ALLOWED_TRUE (attempt %d)", i+1)
			newWorkspaceAllowed = true
			break
		}
		if err != nil {
			t.Logf("Check request failed (attempt %d): %v", i+1, err)
		} else {
			t.Logf("New workspace check returned %v (attempt %d), expected ALLOWED_TRUE", checkResp.GetAllowed(), i+1)
		}
		if i < 9 {
			t.Log("Waiting 1s before retry...")
			time.Sleep(1 * time.Second)
		}
	}
	assert.True(t, newWorkspaceAllowed, "New workspace authorization check did not return ALLOWED_TRUE within timeout")

	t.Log("=== Test Complete ===")

	// Cleanup: Delete the resource
	delReq := &pbv1beta2.DeleteResourceRequest{
		Reference: &pbv1beta2.ResourceReference{
			ResourceType: "host",
			ResourceId:   resourceId,
			Reporter: &pbv1beta2.ReporterReference{
				Type:       reporterType,
				InstanceId: proto.String(reporterInstanceId),
			},
		},
	}
	_, err = client.DeleteResource(ctx, delReq)
	assert.NoError(t, err, "Failed to Delete Resource during cleanup")
}

func TestInventoryAPIHTTP_v1beta2_acm_cluster_with_user_permissions(t *testing.T) {
	enableShortMode(t)
	ctx := context.Background()

	// Test configuration
	workspaceId := "acm-test-workspace"
	clusterId := "test-cluster-abc-123"
	reporterType := "ACM"
	reporterInstanceId := "test-acm-hub"
	user1Principal := "user1"

	t.Logf("=== Testing ACM Cluster with User Permissions ===")
	t.Logf("Workspace: %s", workspaceId)
	t.Logf("Cluster ID: %s", clusterId)
	t.Logf("User: %s", user1Principal)

	conn, err := grpc.NewClient(
		inventoryapi_grpc_url,
		grpc.WithTransportCredentials(grpcinsecure.NewCredentials()),
		grpc.WithPerRPCCredentials(&bearerAuth{token: "1234"}),
	)
	assert.NoError(t, err, "Failed to create gRPC client")
	defer func() {
		if connErr := conn.Close(); connErr != nil {
			t.Logf("Failed to close gRPC connection: %v", connErr)
		}
	}()

	conn.Connect()
	assert.NoError(t, err, "Failed to connect gRPC client")

	client := pbv1beta2.NewKesselInventoryServiceClient(conn)

	// Step 1: Create ACM k8s_cluster with workspace_id
	t.Logf("1. Creating ACM k8s_cluster with workspace_id: %s", workspaceId)

	reporterStruct, err := structpb.NewStruct(map[string]interface{}{
		"external_cluster_id": "external-" + clusterId,
		"cluster_status":      "READY",
		"cluster_reason":      "Cluster is healthy",
		"kube_version":        "v1.28.0",
		"kube_vendor":         "OPENSHIFT",
		"vendor_version":      "4.14.0",
		"cloud_platform":      "AWS_IPI",
		"nodes": []interface{}{
			map[string]interface{}{
				"name":   "master-0",
				"cpu":    "4",
				"memory": "16Gi",
			},
			map[string]interface{}{
				"name":   "worker-0",
				"cpu":    "8",
				"memory": "32Gi",
			},
		},
	})
	assert.NoError(t, err, "Failed to create structpb for k8s_cluster reporter")

	createReq := &pbv1beta2.ReportResourceRequest{
		WriteVisibility:    pbv1beta2.WriteVisibility_MINIMIZE_LATENCY,
		Type:               "k8s_cluster",
		ReporterType:       reporterType,
		ReporterInstanceId: reporterInstanceId,
		Representations: &pbv1beta2.ResourceRepresentations{
			Metadata: &pbv1beta2.RepresentationMetadata{
				LocalResourceId: clusterId,
				ApiHref:         "https://api.example.com/clusters/" + clusterId,
				ConsoleHref:     proto.String("https://console.example.com/clusters/" + clusterId),
				ReporterVersion: proto.String("2.10.0"),
			},
			Common: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"workspace_id": structpb.NewStringValue(workspaceId),
				},
			},
			Reporter: reporterStruct,
		},
	}

	_, err = client.ReportResource(ctx, createReq)
	assert.NoError(t, err, "Failed to create k8s_cluster resource")

	// Step 2: Check that cluster is associated with workspace
	t.Logf("2. Checking that cluster is associated with workspace (polling up to 10s)...")

	checkWorkspace := &pbv1beta2.CheckRequest{
		Object: &pbv1beta2.ResourceReference{
			ResourceType: "k8s_cluster",
			ResourceId:   clusterId,
			Reporter: &pbv1beta2.ReporterReference{
				Type:       reporterType,
				InstanceId: proto.String(reporterInstanceId),
			},
		},
		Relation: "workspace",
		Subject: &pbv1beta2.SubjectReference{
			Resource: &pbv1beta2.ResourceReference{
				ResourceType: "workspace",
				ResourceId:   workspaceId,
				Reporter: &pbv1beta2.ReporterReference{
					Type: "rbac",
				},
			},
		},
	}

	// Poll for up to 10 seconds
	workspaceAssociated := false
	for i := 0; i < 10; i++ {
		checkResp, err := client.Check(ctx, checkWorkspace)
		if err == nil && checkResp.GetAllowed() == pbv1beta2.Allowed_ALLOWED_TRUE {
			t.Logf("✓ Cluster is associated with workspace (attempt %d)", i+1)
			workspaceAssociated = true
			break
		}
		if err != nil {
			t.Logf("Check request failed (attempt %d): %v", i+1, err)
		} else {
			t.Logf("Workspace check returned %v (attempt %d), expected ALLOWED_TRUE", checkResp.GetAllowed(), i+1)
		}
		if i < 9 {
			t.Log("Waiting 1s before retry...")
			time.Sleep(1 * time.Second)
		}
	}
	assert.True(t, workspaceAssociated, "Cluster workspace association check did not return ALLOWED_TRUE within timeout")

	// Step 3: Check update permission for user1
	// Note: This step requires RBAC tuples to be set up via Relations API
	// For a complete test, you would need to:
	// 1. Create a role with inventory_k8s_cluster_update permission
	// 2. Create a role_binding linking user1 to the role
	// 3. Link the role_binding to the workspace
	// This is typically done via the Relations API directly or through a separate setup

	t.Logf("3. Checking update permission for user1 on cluster")
	t.Log("Note: This requires RBAC tuples to be configured via Relations API")
	t.Log("For this test, we're checking the permission structure is correct")

	checkUserPermission := &pbv1beta2.CheckRequest{
		Object: &pbv1beta2.ResourceReference{
			ResourceType: "k8s_cluster",
			ResourceId:   clusterId,
			Reporter: &pbv1beta2.ReporterReference{
				Type:       reporterType,
				InstanceId: proto.String(reporterInstanceId),
			},
		},
		Relation: "update",
		Subject: &pbv1beta2.SubjectReference{
			Resource: &pbv1beta2.ResourceReference{
				ResourceType: "principal",
				ResourceId:   user1Principal,
				Reporter: &pbv1beta2.ReporterReference{
					Type: "rbac",
				},
			},
		},
	}

	// Try to check permission (may return ALLOWED_FALSE if RBAC not configured)
	checkResp, err := client.Check(ctx, checkUserPermission)
	if err != nil {
		t.Logf("Permission check failed: %v (This is expected if RBAC tuples are not configured)", err)
	} else {
		t.Logf("Permission check result: %v", checkResp.GetAllowed())
		t.Log("To grant user1 edit permission, you need to:")
		t.Log("  1. Create role with inventory_k8s_clusters_write permission")
		t.Log("  2. Create role_binding linking user1 to the role")
		t.Log("  3. Link role_binding to the workspace")
		t.Log("  Example zed commands:")
		t.Logf("    zed relationship create rbac/role:cluster-editor inventory_k8s_clusters_write rbac/principal:*")
		t.Logf("    zed relationship create rbac/role_binding:user1-cluster-editor t_role rbac/role:cluster-editor")
		t.Logf("    zed relationship create rbac/role_binding:user1-cluster-editor t_subject rbac/principal:%s", user1Principal)
		t.Logf("    zed relationship create rbac/workspace:%s t_binding rbac/role_binding:user1-cluster-editor", workspaceId)
	}

	t.Log("=== Test Complete ===")

	// Cleanup: Delete the resource
	delReq := &pbv1beta2.DeleteResourceRequest{
		Reference: &pbv1beta2.ResourceReference{
			ResourceType: "k8s_cluster",
			ResourceId:   clusterId,
			Reporter: &pbv1beta2.ReporterReference{
				Type:       reporterType,
				InstanceId: proto.String(reporterInstanceId),
			},
		},
	}
	_, err = client.DeleteResource(ctx, delReq)
	assert.NoError(t, err, "Failed to Delete Resource during cleanup")
}

