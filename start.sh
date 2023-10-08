#!/bin/sh

set -e

echo "run db migration"
/app/migrate -path /app/migrations -database "mysql://root:root456789@tcp(mysql:3306)/jungle?multiStatements=true" -verbose up

echo "start the app"
exec "$@"