# ------------------------------------------------------------------------------
# Build
# ------------------------------------------------------------------------------

FROM golang:1.18 AS build

RUN apt install make -qq

WORKDIR /workspace

COPY go.mod go.mod

COPY go.sum go.sum

COPY main.go main.go

RUN CGO_ENABLED=0 go build -a -o /workspace/malutki main.go

# ------------------------------------------------------------------------------
# Base
# ------------------------------------------------------------------------------

FROM gcr.io/distroless/static:nonroot AS distroless

LABEL org.opencontainers.image.source https://github.com/shaneutt/malutki

WORKDIR /

COPY --from=build /workspace/malutki /malutki

USER 65532:65532

EXPOSE 8080

ENTRYPOINT ["/malutki"]
