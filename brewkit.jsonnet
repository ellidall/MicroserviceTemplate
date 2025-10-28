local project = import 'brewkit/project.libsonnet';

// TODO: appID поменять

local appIDs = [
    'microservicetemplate',
];

local proto = [
    'api/client/testclientinternal/test.proto',
    'api/server/microservicetemplateinternal/test.proto',
];

project.project(appIDs, proto)