# Requires vendors to be built prior to docker run.
# docker build -f packaging/Dockerfile .
# docker tag XXX charlyyfon/ebpf:metadata
# docker push charlyyfon/ebpf:metadata
FROM golang:latest
RUN mkdir -p /go/src/github.com/CharlyF/cluster-monitoring
ADD . /go/src/github.com/CharlyF/cluster-monitoring/
WORKDIR /go/src/github.com/CharlyF/cluster-monitoring
RUN go build -o main ./cmd/metadata-collector


