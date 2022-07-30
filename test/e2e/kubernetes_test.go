//go:build e2e_tests
// +build e2e_tests

package integration

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/loadimage"
	"github.com/kong/kubernetes-testing-framework/pkg/clusters/addons/metallb"
	"github.com/kong/kubernetes-testing-framework/pkg/environments"
	"github.com/kong/kubernetes-testing-framework/pkg/utils/kubernetes/generators"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const port = 8080

var testImage = os.Getenv("MALUTKI_TEST_IMAGE")

func TestKubernetesBuild(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Log("verifying container image")
	require.NotEmpty(t, testImage, "a container image must be provided")
	imageBuilder, err := loadimage.NewBuilder().WithImage(testImage)
	require.NoError(t, err)

	t.Log("building kubernetes cluster")
	env, err := environments.NewBuilder().WithAddons(metallb.New(), imageBuilder.Build()).Build(ctx)
	require.NoError(t, err)
	require.NoError(t, <-env.WaitForReady(ctx))
	defer env.Cleanup(ctx)

	t.Log("creating a Deployment")
	container := generators.NewContainer("malutki", testImage, port)
	deployment := generators.NewDeploymentForContainer(container)
	deployment.Spec.Template.Spec.Containers[0].ImagePullPolicy = corev1.PullNever // use loadimage
	deployment, err = env.Cluster().Client().AppsV1().Deployments(corev1.NamespaceDefault).Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("exposing Deployment via LoadBalancer Service")
	service := generators.NewServiceForDeployment(deployment, corev1.ServiceTypeLoadBalancer)
	service, err = env.Cluster().Client().CoreV1().Services(corev1.NamespaceDefault).Create(ctx, service, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Logf("waiting for LoadBalancer Service to receive an IP address")
	var ip string
	require.Eventually(t, func() bool {
		service, err = env.Cluster().Client().CoreV1().Services(corev1.NamespaceDefault).Get(ctx, service.Name, metav1.GetOptions{})
		require.NoError(t, err)
		if len(service.Status.LoadBalancer.Ingress) > 0 {
			if service.Status.LoadBalancer.Ingress[0].IP != "" {
				ip = service.Status.LoadBalancer.Ingress[0].IP
				return true
			}
		}
		return false
	}, time.Minute, time.Second)

	t.Log("verifying connectivity to malutki")
	httpc := &http.Client{Timeout: time.Second * 10}
	require.Eventually(t, func() bool {
		resp, err := httpc.Get(fmt.Sprintf("http://%s:%d/status/200", ip, port))
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, time.Minute, time.Second)

	t.Log("verifying 2XX and 5XX status code routing")
	for code := 200; code < 600; code++ {
		if code == 300 {
			code = 500 // jump from 2XX to 5XX
		}
		if http.StatusText(code) != "" {
			resp, err := httpc.Get(fmt.Sprintf("http://%s:%d/status/%d", ip, port, code))
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, code, resp.StatusCode)
		}
	}
}
