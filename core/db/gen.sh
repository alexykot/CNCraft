#!/usr/bin/env bash

DB_URL=$1

if [[ $DB_URL == "" ]]; then
  DB_URL=$(echo $CNCRAFT_TEST_DB_URL)
fi

if [[ $DB_URL == "" ]]; then
  echo "DB URL must be provided, as an argument or via CNCRAFT_TEST_DB_URL envar"
  exit
fi


DB_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
pushd "$DB_ROOT" >/dev/null || exit 1

# parse the DB_URL into variables for sqlboiler config
protocol="$(echo $DB_URL | grep :// | sed -e's,^\(.*://\).*,\1,g')"
url="$(echo ${DB_URL/$protocol/})"
userpass="$(echo $url | grep @ | cut -d@ -f1)"
user="$(echo $userpass | sed -e 's,:.*,,g')"
pass="$(echo ${userpass/$user:/})"
hostport="$(echo ${url/$userpass@/} | cut -d/ -f1)"
host="$(echo $hostport | sed -e 's,:.*,,g')"
port="$(echo $hostport | sed -e 's,^.*:,:,g' -e 's,.*:\([0-9]*\).*,\1,g' -e 's,[^0-9],,g')"
path="$(echo $url | grep / | cut -d/ -f2-)"
dbname="$(echo $path | sed -e 's,?.*,,g')"
pkgname=$(echo $SERVICE | tr - _)

if [[ $user == "" ]]; then
  user="postgres"
fi

# generate fresh schema migrations bindata file
go-bindata -nometadata -o ./schema.go -prefix ./schema -pkg db -ignore=".*\\.go|BUILD.bazel|.DS_Store" ./schema

# run migrations before generating the ORM, make sure we generate for the latest schema
go run ./migrator/migrator.go -db-url "$(echo $DB_URL)"

# output TOML config from template
function output_sqlboiler_config() {
  cat <<EOF
pkgname = "orm"
output = "$DB_ROOT/orm"
[psql]
  dbname = "$dbname"
  host   = "$host"
  port   = $port
  user   = "$user"
  pass   = "$pass"
  schema = "cncraft"
  sslmode = "disable"
  blacklist = ["cncraft_schema_migrations"]
[[types]]
  [types.match]
    type = "types.Decimal"
  [types.replace]
    type = "apd.Decimal"
  [types.imports]
    third_party = ['"github.com/cockroachdb/apd/v2"']
[[types]]
  [types.match]
    type = "types.NullDecimal"
    nullable = true
  [types.replace]
    type = "apd.NullDecimal"
  [types.imports]
    third_party = ['"github.com/cockroachdb/apd/v2"']
# Uncomment this if we will want to use UUID natively in the ORM.
# Not using it now as strings need less overhead and we are not using any
# UUID-specific functions.
# [[types]]
#   [types.match]
#     db_type = "uuid"
#   [types.replace]
#     type = "uuid.UUID"
#   [types.imports]
#     third_party = ['"github.com/google/uuid"']
EOF
}

# save temp config file
output_sqlboiler_config > ./sqlboiler.tmp.toml

# produce the generated ORM
sqlboiler --no-tests --no-hooks -c ./sqlboiler.tmp.toml psql

# remove temp config file
rm ./sqlboiler.tmp.toml

popd || exit 1
