FROM golang:1.18-alpine AS build
WORKDIR /go/src/gogin
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/gogin .


# Copy the binaries of our service the new lightweigh container
FROM scratch
COPY --from=build /go/bin/gogin /bin/gogin
ENTRYPOINT ["/bin/gogin"]