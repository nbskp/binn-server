FROM golang:1.19.2-bullseye AS build
ENV GO111MODULE=on
WORKDIR /go/src/github.com/nubesk/binn
COPY ./ ./
RUN CGO_ENABLED=0 go build -o /go/src/github.com/nubesk/binn/build/binn main.go

FROM gcr.io/distroless/static-debian11
COPY --from=build /go/src/github.com/nubesk/binn/build/binn /
CMD ["/binn"]
