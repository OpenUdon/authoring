package operationlifecycle

import (
	"slices"
	"testing"

	"github.com/OpenUdon/authoring/promptcontext"
)

func TestExpandKubernetesCollectionItemLifecycle(t *testing.T) {
	ctx := promptcontext.Context{Operations: []promptcontext.OperationCandidate{
		op("k8s", "createCoreV1NamespacedConfigMap", "POST", "/api/v1/namespaces/{namespace}/configmaps"),
		op("k8s", "readCoreV1NamespacedConfigMap", "GET", "/api/v1/namespaces/{namespace}/configmaps/{name}"),
		op("k8s", "replaceCoreV1NamespacedConfigMap", "PUT", "/api/v1/namespaces/{namespace}/configmaps/{name}"),
		op("k8s", "deleteCoreV1NamespacedConfigMap", "DELETE", "/api/v1/namespaces/{namespace}/configmaps/{name}"),
	}}
	expanded := Expand(ctx, ctx.Operations[0], Options{Goal: "create, read, update, and delete configmaps", DesiredState: true})
	if got := roleIDs(expanded); !slices.Equal(got, []string{"create:createCoreV1NamespacedConfigMap", "read:readCoreV1NamespacedConfigMap", "update:replaceCoreV1NamespacedConfigMap", "delete:deleteCoreV1NamespacedConfigMap"}) {
		t.Fatalf("roles = %#v", got)
	}
}

func TestExpandGoogleDotOperationIDs(t *testing.T) {
	ctx := promptcontext.Context{Operations: []promptcontext.OperationCandidate{
		op("storage", "storage.buckets.insert", "POST", "/b"),
		op("storage", "storage.buckets.get", "GET", "/b/{bucket}"),
		op("storage", "storage.buckets.patch", "PATCH", "/b/{bucket}"),
		op("storage", "storage.buckets.update", "PUT", "/b/{bucket}"),
		op("storage", "storage.buckets.delete", "DELETE", "/b/{bucket}"),
	}}
	expanded := Expand(ctx, ctx.Operations[0], Options{Goal: "create, read, update, and delete buckets", DesiredState: true})
	if got := roleIDs(expanded); !slices.Equal(got, []string{"create:storage.buckets.insert", "read:storage.buckets.get", "update:storage.buckets.patch", "delete:storage.buckets.delete"}) {
		t.Fatalf("roles = %#v", got)
	}
}

func TestExpandGoogleUploadPathsMatchCanonicalObjectPaths(t *testing.T) {
	ctx := promptcontext.Context{Operations: []promptcontext.OperationCandidate{
		op("storage", "storage.objects.insert", "POST", "/upload/storage/v1/b/{bucket}/o"),
		op("storage", "storage.objects.get", "GET", "/storage/v1/b/{bucket}/o/{object}"),
		op("storage", "storage.objects.patch", "PATCH", "/storage/v1/b/{bucket}/o/{object}"),
		op("storage", "storage.objects.delete", "DELETE", "/storage/v1/b/{bucket}/o/{object}"),
	}}
	expanded := Expand(ctx, ctx.Operations[0], Options{Goal: "create, read, update, and delete objects", DesiredState: true})
	if got := roleIDs(expanded); !slices.Equal(got, []string{"create:storage.objects.insert", "read:storage.objects.get", "update:storage.objects.patch", "delete:storage.objects.delete"}) {
		t.Fatalf("roles = %#v", got)
	}
}

func TestNormalizePathOnlyStripsUploadSegment(t *testing.T) {
	if got := normalizePath("/upload/storage/v1/b/{bucket}/o"); got != "/storage/v1/b/{bucket}/o" {
		t.Fatalf("upload path = %q", got)
	}
	if got := normalizePath("/uploading/storage/v1/b/{bucket}/o"); got != "/uploading/storage/v1/b/{bucket}/o" {
		t.Fatalf("uploading path = %q", got)
	}
}

func TestExpandCloudflareHyphenOperationIDs(t *testing.T) {
	ctx := promptcontext.Context{Operations: []promptcontext.OperationCandidate{
		op("cloudflare", "r2-create-bucket", "POST", "/accounts/{account_id}/r2/buckets"),
		op("cloudflare", "r2-get-bucket", "GET", "/accounts/{account_id}/r2/buckets/{bucket_name}"),
		op("cloudflare", "r2-patch-bucket", "PATCH", "/accounts/{account_id}/r2/buckets/{bucket_name}"),
		op("cloudflare", "r2-delete-bucket", "DELETE", "/accounts/{account_id}/r2/buckets/{bucket_name}"),
	}}
	expanded := Expand(ctx, ctx.Operations[0], Options{Goal: "create, read, update, and delete buckets", DesiredState: true})
	if got := roleIDs(expanded); !slices.Equal(got, []string{"create:r2-create-bucket", "read:r2-get-bucket", "update:r2-patch-bucket", "delete:r2-delete-bucket"}) {
		t.Fatalf("roles = %#v", got)
	}
}

func TestExpandAzureCreateOrUpdateLifecycle(t *testing.T) {
	ctx := promptcontext.Context{Operations: []promptcontext.OperationCandidate{
		op("azure", "Databases_CreateOrUpdate", "PUT", "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/databases/{databaseName}"),
		op("azure", "Databases_Get", "GET", "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/databases/{databaseName}"),
		op("azure", "Databases_Delete", "DELETE", "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/databases/{databaseName}"),
	}}
	expanded := Expand(ctx, ctx.Operations[0], Options{DesiredState: true})
	if got := roleIDs(expanded); !slices.Equal(got, []string{"create:Databases_CreateOrUpdate", "read:Databases_Get", "delete:Databases_Delete"}) {
		t.Fatalf("roles = %#v", got)
	}
}

func TestExpandRejectsAmbiguousSibling(t *testing.T) {
	ctx := promptcontext.Context{Operations: []promptcontext.OperationCandidate{
		op("widgets", "createWidget", "POST", "/widgets"),
		op("widgets", "getWidget", "GET", "/widgets/{id}"),
		op("widgets", "readWidget", "GET", "/widgets/{id}"),
	}}
	expanded := Expand(ctx, ctx.Operations[0], Options{DesiredState: true})
	if got := roleIDs(expanded); !slices.Equal(got, []string{"post:createWidget"}) {
		t.Fatalf("roles = %#v", got)
	}
	if len(expanded.Diagnostics) != 1 || expanded.Diagnostics[0].Code != "operation_lifecycle.ambiguous_read" {
		t.Fatalf("diagnostics = %#v", expanded.Diagnostics)
	}
}

func TestExpandPreservesAPIFirstSingleOperationRoles(t *testing.T) {
	for _, tc := range []struct {
		id   string
		verb string
		role string
	}{
		{id: "createWidget", verb: "POST", role: "post"},
		{id: "putWidget", verb: "PUT", role: "put"},
		{id: "deleteWidget", verb: "DELETE", role: "delete"},
	} {
		ctx := promptcontext.Context{Operations: []promptcontext.OperationCandidate{op("api", tc.id, tc.verb, "/widgets/{id}")}}
		expanded := Expand(ctx, ctx.Operations[0], Options{})
		if got := roleIDs(expanded); !slices.Equal(got, []string{tc.role + ":" + tc.id}) {
			t.Fatalf("%s roles = %#v", tc.id, got)
		}
	}
}

func op(source, id, verb, path string) promptcontext.OperationCandidate {
	return promptcontext.OperationCandidate{ID: id, SourceID: source, OperationID: id, Verb: verb, Path: path}
}

func roleIDs(expanded Expansion) []string {
	var out []string
	for _, role := range expanded.Roles {
		out = append(out, role.Role+":"+role.Operation.OperationID)
	}
	return out
}
