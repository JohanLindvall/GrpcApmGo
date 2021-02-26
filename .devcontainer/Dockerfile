FROM golang:1.16.0-buster
RUN apt update && \
    apt install -y protobuf-compiler && \
    go get golang.org/x/tools/gopls@latest && \
    go get github.com/ramya-rao-a/go-outline@latest && \
    go get github.com/go-delve/delve/cmd/dlv@latest && \
    go get github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest && \
    go get golang.org/x/lint/golint@latest && \
    go get github.com/jstemmer/go-junit-report@latest && \ 
    go get google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest