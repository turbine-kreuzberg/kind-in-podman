package test

import (
	"os"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
	"sigs.k8s.io/e2e-framework/support/kind"
)

var (
	testenv        env.Environment
	namespace      string = envconf.RandomName("nginx-test", 16)
	deploymentName string = "nginx-deployment"
)

func TestMain(m *testing.M) {
	kindClusterName := "kind"
	testenv = env.New()

	testenv.Setup(
		envfuncs.CreateClusterWithConfig(
			kind.NewProvider(),
			kindClusterName,
			"kind-in-podman.yaml",
			kind.WithImage("kindest/node:v1.27.3"),
		),
		envfuncs.CreateNamespace(namespace),
	)

	testenv.Finish(
		envfuncs.DeleteNamespace(namespace),
		envfuncs.ExportClusterLogs(kindClusterName, "./logs"),
		envfuncs.DestroyCluster(kindClusterName),
	)

	os.Exit(testenv.Run(m))
}
