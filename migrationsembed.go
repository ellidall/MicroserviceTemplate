package microservicetemplate

import "embed"

//go:embed data/mysql/migrations/*.sql
var Migrations embed.FS
