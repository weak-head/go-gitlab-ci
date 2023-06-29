##############################################
## Step 1: Build Gogin micro-service ------- #
##############################################

FROM golang:1.20-alpine AS build

ARG GRPC_HEALTH_PROBE_VERSION=v0.4.19

WORKDIR /go/src/gogin
COPY . .

# Compile Gogin micro-service
RUN CGO_ENABLED=0 go build -o /go/bin/gogin ./cmd/gogin

# Support for liveleness and readiness probes via gRPC
RUN wget -qO/go/bin/grpc_health_probe \
    https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /go/bin/grpc_health_probe

#############################################
## Step 2: Create Deployable Image ------- ##
#############################################

# Copy the binaries of our service the new lightweigh container
FROM scratch

# Copy the binaries that were built on the 'build' step to this container
COPY --from=build /go/bin/gogin /bin/gogin
COPY --from=build /go/bin/grpc_health_probe /bin/grpc_health_probe

CMD ["/bin/gogin"]