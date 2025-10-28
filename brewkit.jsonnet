local project = import 'brewkit/project.libsonnet';

// TODO: appID поменять

local appIDs = [
    'microservicetemplate',
];

local proto = [
    'api/client/testinternal/testinternal.proto',
    'api/server/microservicetemplateinternal/microservicetemplateinternal.proto',
];

project.project(appIDs, proto)