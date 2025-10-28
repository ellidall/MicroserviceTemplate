local images = import 'images.libsonnet';

local copy = std.native('copy');

{
    generateGRPC(protoFiles):: {
        local mappedFiles = [copy(protoFile, protoFile) for protoFile in protoFiles],

        from: images.protoc,
        workdir: "/app",
        copy: mappedFiles,
        command: std.join(' && ', [
                'mkdir -p pkg',
                'find . -name "*.proto" -exec protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative,require_unimplemented_servers=false:. {} \\;'
        ]),
        output: {
            artifact: "/app/api",
            "local": "./api"
        },
    },
}