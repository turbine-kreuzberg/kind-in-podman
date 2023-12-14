package test

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"os"
	"sigs.k8s.io/e2e-framework/klient/decoder"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	"testing"
	"time"
)

func TestKindCluster(t *testing.T) {
	testdata := os.DirFS("k8s-resources")
	pattern := "*"

	feature := features.
		New("Testing nginx deployment").
		Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			resource, err := resources.New(cfg.Client().RESTConfig())
			if err != nil {
				t.Fatal(err)
			}

			if err := decoder.DecodeEachFile(ctx, testdata, pattern,
				decoder.CreateHandler(resource),
				decoder.MutateNamespace(namespace),
			); err != nil {
				t.Fatal(err)
			}

			return ctx
		}).
		Assess(
			"Nginx was deployed",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				var deployment appsv1.Deployment
				client, err := cfg.NewClient()
				if err != nil {
					t.Fatal(err)
				}

				if err := client.Resources(namespace).Get(ctx, deploymentName, namespace, &deployment); err != nil {
					t.Fatal(err)
				}

				err = wait.For(conditions.New(cfg.Client().Resources()).DeploymentConditionMatch(
					&deployment,
					appsv1.DeploymentAvailable,
					corev1.ConditionTrue,
				),
					wait.WithTimeout(time.Minute*1),
				)
				if err != nil {
					t.Fatal(err)
				}
				return ctx
			}).
		Assess(
			"Pods became ready",
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				client, err := cfg.NewClient()
				podList := corev1.PodList{}

				err = client.Resources(namespace).List(context.TODO(), &podList, resources.WithLabelSelector(
					labels.FormatLabels(map[string]string{"app": "nginx"})),
				)

				if err != nil {
					t.Fatal(err)
				}

				if len(podList.Items) == 0 {
					t.Fatal("could not find any pod in the namespace ", namespace)
				}

				for i := range podList.Items {
					err = wait.For(
						conditions.New(client.Resources().WithNamespace(namespace)).PodReady(&podList.Items[i]), wait.WithTimeout(time.Minute*2),
					)
				}

				return ctx
			}).
		Teardown(
			func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				resource, err := resources.New(cfg.Client().RESTConfig())
				if err != nil {
					t.Fatal(err)
				}
				if err := decoder.DecodeEachFile(ctx, testdata, pattern,
					decoder.DeleteHandler(resource),
					decoder.MutateNamespace(namespace),
				); err != nil {
					t.Fatal(err)
				}

				return ctx
			}).Feature()
	testenv.Test(t, feature)
}
