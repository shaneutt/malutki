# Malutki

A tiny HTTP server used to validate HTTP traffic in testing scenarios.

## About

This HTTP server is meant to be used in integration, E2E and manual testing
scenarios to validate HTTP traffic routing. The initial use case for `malutki`
was in fact to be deployed to [Kubernetes][k8s] [Pods][pods] in order to
validate that an [ingress controller][ing] was properly routing HTTP traffic to
those `Pods`, and to test [service mesh][mesh] observability.

This server provides testing APIs that can be used to test certain kinds of HTTP
responses, including a `/status/{code}` endpoint which simply returns the HTTP
status it is provided. This can be helpful when your tests require a specific
response or when you're testing tools that trace and observe HTTP traffic over
a network (such as the [observability features][obs] of a service mesh like
[Istio][istio]).

It is a core design goal of this project to ensure that the binary and container
images for this tool are extremely small and portable. It is an eventual goal to
build using [TinyGo][tgo] once the project [supports net/http][tgosup], so this
tool intentionally uses minimal dependencies in preparation for that transition.
Container images are built using [distroless][dless] to minimize size.

Trivia: "malutki" means "tiny" in Polish.

[k8s]:https://kubernetes.io
[pods]:https://kubernetes.io/docs/concepts/workloads/pods/
[ing]:https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/
[mesh]:https://wikipedia.org/wiki/Service_mesh
[obs]:https://istio.io/latest/docs/concepts/observability/
[istio]:https://istio.io
[tgo]:https://tinygo.org
[tgosup]:https://tinygo.org/docs/reference/lang-support/stdlib
[dless]:https://github.com/GoogleContainerTools/distroless

## Quickstart (with Kubernetes)

You can deploy Malutki on your [Kubernetes] cluster using the example from this
repository:

```console
$ kubectl apply -f examples/deployment-with-loadbalancer-service.yaml
```

This will create a Kubernetes `Deployment` which will be exposed outside of the
cluster via a `LoadBalancer` type `Service`.

[Kubernetes]:https://github.com/kubernetes/kubernetes

## Usage

This testing server is broken up into different APIs which provide various
testing capabilities.

### /status/{code} API

This is a basic API that is used to return the [HTTP Status Code][status] that
you send it back to you. This is helpful if you're testing the machinery that
provisions the server and just want a simple "it's working" response, e.g.:

```console
$ curl -w '%{http_code}\n' ${MALUTKI_ADDR}/status/201
201
```

Only `2XX`, `4XX` and `5XX` status codes are currently supported.

[status]:https://developer.mozilla.org/docs/Web/HTTP/Status

## Contributing

Contributions are welcome! Please feel free to create [issues][iss],
[discussions][disc] and [pull requests][prs].

[iss]:https://github.com/shaneutt/malutki/issues
[disc]:https://github.com/shaneutt/malutki/discussions
[prs]:https://github.com/shaneutt/malutki/pulls
