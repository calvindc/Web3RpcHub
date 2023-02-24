#!/bin/sh



set -e

# ensure tools are installed
go get github.com/volatiletech/sqlboiler/v4
go get github.com/volatiletech/sqlboiler-sqlite3

# make sure we are in the correct directory
cd "$(dirname $0)"

# run the migrations (creates testrun/TestSchema/hubdb)
go test -run Schema

 make sure the sqlite file was created
test -f testrun/TestSchema/hubdb || {
    echo 'roomdb file missing'
    exit 1
}

# generate the models package
sqlboiler sqlite3 --wipe --no-tests

echo "all done. models updated!"
