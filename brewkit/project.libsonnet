local images = import 'images.libsonnet';
local schemas = import 'schemas.libsonnet';

local cache = std.native('cache');
local copy = std.native('copy');
local copyFrom = std.native('copyFrom');

// External cache for go compiler, go mod, golangci-lint
local gocache = [
    cache("go-build", "/app/cache"),
    cache("go-mod", "/go/pkg/mod"),
];

// Sources which will be tracked for changes
local gosources = [
    "go.mod",
    "go.sum",
    "data",
    "cmd",
    "api",
    "pkg",
];

{
    // Function that generate project build definitions, including code generating, app compilation and e.t.c
    project(appIDs, protos):: {
        apiVersion: "brewkit/v1",

        targets: {
            all: ['build', 'test', 'check'],

            // build target to chain all build of apps
            build: [appID for appID in appIDs],

            gobase: {
                from: images.gobuilder,
                workdir: "/app",
                env: {
                    GOCACHE: "/app/cache/go-build",
                },
                copy: copyFrom(
                    'gosources',
                    '/app',
                    '/app'
                ),
            },
        } + {
            [appID]: {
                from: "gobase",
                workdir: "/app",
                cache: gocache,
                dependsOn: ['generate', 'modules'],
                command: 'go build \\
                        -trimpath -v \\
                        -o ./bin/' + appID + ' ./cmd/' + appID,
                output: {
                    artifact: "/app/bin/" + appID,
                    "local": "./bin",
                },
            }
            for appID in appIDs // expand build target for each appID
        } + {
            gosources: {
                from: "scratch",
                workdir: "/app",
                copy: [copy(source, source) for source in gosources]
            },

            generate: ['generategrpc'],

            generategrpc: schemas.generateGRPC(protos),

            modules: ["gotidy"],

            gotidy: {
              from: "gobase",
              workdir: "/app",
              cache: gocache,
              command: "go mod tidy",
              output: {
                artifact: "/app/go.*",
                "local": ".",
              },
            },

            test: {
                from: "gobase",
                workdir: "/app",
                cache: gocache,
                command: "go test ./...",
            },

            check: {
                from: images.golangcilint,
                workdir: "/app",
                env: {
                    GOCACHE: "/app/cache/go-build",
                    GOLANGCI_LINT_CACHE: "/app/cache/go-build",
                },
                cache: gocache,
                copy: [
                    copy('.golangci.yml', '.golangci.yml'),
                    copyFrom(
                        'gosources',
                        '/app',
                        '/app'
                    ),
                ],
                command: "golangci-lint run",
            },
        },
    },
}
