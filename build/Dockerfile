FROM golang:1.21.3-bullseye as build

WORKDIR /build

COPY internal internal
COPY go.mod .
COPY go.sum .
COPY main.go .
RUN CGO_ENABLED=0 go build -o /svc ./main.go

FROM alpine

COPY --from=build /svc /svc

ENTRYPOINT [ "./svc" ]
