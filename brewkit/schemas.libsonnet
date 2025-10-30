local images = import 'images.libsonnet';

local copy = std.native('copy');

local PROTOC_VERSION = "27.3";

{
    generateGRPC(protoFiles):: {
        from: "golang:1.23-bookworm",
        workdir: "/app",
        env: {
            PATH: "/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
        },
        copy: [copy(f, f) for f in protoFiles],
        command: std.join(" && ", [
            "apt-get update && apt-get install -y --no-install-recommends unzip",

            "curl -sSL https://github.com/protocolbuffers/protobuf/releases/download/v" + PROTOC_VERSION + "/protoc-" + PROTOC_VERSION + "-linux-x86_64.zip -o /tmp/protoc.zip",
            "unzip -q /tmp/protoc.zip -d /usr/local",
            "rm /tmp/protoc.zip",

            "go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2",
            "go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1",

            "find . -name '*.proto' -exec protoc --proto_path=. --proto_path=/usr/local/include --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative,require_unimplemented_servers=false:. {} \\;"
        ]),
        output: {
            artifact: "/app/api",
            "local": "./api"
        },
    },
}
