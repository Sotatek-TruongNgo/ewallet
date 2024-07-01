FROM golang:1.22.3 as builder
WORKDIR /workspace
RUN apt-get update && apt-get install -y --no-install-recommends make curl unzip && rm -rf /var/lib/apt/lists/*
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0
COPY Makefile Makefile
RUN make install-bin
COPY go.mod go.sum ./
RUN make dep
COPY . /workspace
RUN make doc build

FROM alpine
WORKDIR /workspace
RUN apk upgrade && apk add --no-cache ca-certificates
COPY --from=builder /workspace/build ./build
COPY --from=builder /workspace/config ./config
ENTRYPOINT ["/workspace/build/server/server"]
CMD ["--method=server"]