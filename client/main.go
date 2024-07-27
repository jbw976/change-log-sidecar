package main

import (
	"context"
	"fmt"
	"log"

	crv1alpha1 "github.com/crossplane/crossplane-runtime/apis/proto/v1alpha1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
)

func main() {
	socketPath := "/var/run/change-logs/change-logs.sock"

	conn, err := grpc.NewClient("unix://"+socketPath, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial server: %+v", err)
	}

	client := crv1alpha1.NewChangeLogServiceClient(conn)

	beforeState := "{\"apiVersion\":\"kubernetes.crossplane.io/v1alpha2\",\"kind\":\"Object\",\"metadata\":{\"annotations\":{\"crossplane.io/composition-resource-name\":\"resource-object-1\",\"crossplane.io/external-create-pending\":\"2024-07-26T06:20:47Z\",\"crossplane.io/external-name\":\"object-1\"},\"creationTimestamp\":\"2024-07-26T06:20:47Z\",\"finalizers\":[\"finalizer.managedresource.crossplane.io\"],\"generateName\":\"traceperf-tester-2bjgk-\",\"generation\":2,\"labels\":{\"crossplane.io/claim-name\":\"traceperf-tester\",\"crossplane.io/claim-namespace\":\"default\",\"crossplane.io/composite\":\"traceperf-tester-2bjgk\"},\"managedFields\":[{\"apiVersion\":\"kubernetes.crossplane.io/v1alpha2\",\"fieldsType\":\"FieldsV1\",\"fieldsV1\":{\"f:metadata\":{\"f:annotations\":{\"f:crossplane.io/composition-resource-name\":{}},\"f:generateName\":{},\"f:labels\":{\"f:crossplane.io/claim-name\":{},\"f:crossplane.io/claim-namespace\":{},\"f:crossplane.io/composite\":{}},\"f:ownerReferences\":{\"k:{\\\"uid\\\":\\\"c2177d42-a211-4e4b-bc56-7c255bb9d736\\\"}\":{}}},\"f:spec\":{\"f:forProvider\":{\"f:manifest\":{\"f:apiVersion\":{},\"f:data\":{\".\":{},\"f:key\":{}},\"f:kind\":{},\"f:metadata\":{\"f:namespace\":{}}}}}},\"manager\":\"apiextensions.crossplane.io/composed/ef08e4db35abd8b42eb51aba418b5b815767c572772cb6c2665857cad5848ba4\",\"operation\":\"Apply\",\"time\":\"2024-07-26T06:20:47Z\"},{\"apiVersion\":\"kubernetes.crossplane.io/v1alpha2\",\"fieldsType\":\"FieldsV1\",\"fieldsV1\":{\"f:metadata\":{\"f:annotations\":{\"f:crossplane.io/external-create-pending\":{},\"f:crossplane.io/external-name\":{}},\"f:finalizers\":{\".\":{},\"v:\\\"finalizer.managedresource.crossplane.io\\\"\":{}}},\"f:spec\":{\"f:readiness\":{\".\":{},\"f:policy\":{}}}},\"manager\":\"crossplane-kubernetes-provider\",\"operation\":\"Update\",\"time\":\"2024-07-26T06:20:47Z\"}],\"name\":\"object-1\",\"ownerReferences\":[{\"apiVersion\":\"trace-perf.crossplane.io/v1alpha1\",\"blockOwnerDeletion\":true,\"controller\":true,\"kind\":\"XTracePerf\",\"name\":\"traceperf-tester-2bjgk\",\"uid\":\"c2177d42-a211-4e4b-bc56-7c255bb9d736\"}],\"resourceVersion\":\"16605\",\"uid\":\"7935545a-c561-4dfa-9345-0555751489e0\"},\"spec\":{\"deletionPolicy\":\"Delete\",\"forProvider\":{\"manifest\":{\"apiVersion\":\"v1\",\"data\":{\"key\":\"value-1\"},\"kind\":\"ConfigMap\",\"metadata\":{\"namespace\":\"default\"}}},\"managementPolicies\":[\"*\"],\"providerConfigRef\":{\"name\":\"default\"},\"readiness\":{\"policy\":\"SuccessfulCreate\"}},\"status\":{\"atProvider\":{\"manifest\":{\"apiVersion\":\"v1\",\"data\":{\"key\":\"value-1\"},\"kind\":\"ConfigMap\",\"metadata\":{\"annotations\":{\"kubectl.kubernetes.io/last-applied-configuration\":\"{\\\"apiVersion\\\":\\\"v1\\\",\\\"data\\\":{\\\"key\\\":\\\"value-1\\\"},\\\"kind\\\":\\\"ConfigMap\\\",\\\"metadata\\\":{\\\"namespace\\\":\\\"default\\\"}}\"},\"creationTimestamp\":\"2024-07-26T06:20:47Z\",\"managedFields\":[{\"apiVersion\":\"v1\",\"fieldsType\":\"FieldsV1\",\"fieldsV1\":{\"f:data\":{\".\":{},\"f:key\":{}},\"f:metadata\":{\"f:annotations\":{\".\":{},\"f:kubectl.kubernetes.io/last-applied-configuration\":{}}}},\"manager\":\"crossplane-kubernetes-provider\",\"operation\":\"Update\",\"time\":\"2024-07-26T06:20:47Z\"}],\"name\":\"object-1\",\"namespace\":\"default\",\"resourceVersion\":\"16606\",\"uid\":\"d7d6fef6-e968-4ab0-9291-12cce7eecdb5\"}}}}}"
	beforeStateStruct := &structpb.Struct{}
	err = beforeStateStruct.UnmarshalJSON([]byte(beforeState))
	if err != nil {
		log.Fatalf("failed to unmarshal before state: %+v", err)
	}

	changeErrorMessage := "simulated change failure"

	entry := &crv1alpha1.ChangeLogEntry{
		Provider:          "provider-unknown:v9.99.999",
		Type:              "kubernetes.crossplane.io/v1alpha2, Kind=Object",
		Name:              "object-0",
		ExternalName:      "object-0",
		Operation:         crv1alpha1.OperationType_OPERATION_CREATE,
		BeforeState:       beforeStateStruct,
		AfterState:        &structpb.Struct{},
		ErrorMessage:      &changeErrorMessage,
		AdditionalDetails: &structpb.Struct{},
	}

	resp, err := client.SendChangeLog(context.TODO(), entry)
	if err != nil {
		log.Fatalf("failed to send change log entry: %+v", err)
	}

	fmt.Printf("received response: %v\n", resp)
}
