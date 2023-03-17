#!/bin/sh

set -e



# 1、数据库使用sql-migrate制作迁移文件和管理 https://github.com/rubenv/sql-migrate
# sql-migrate new base,根据dbconfig.yml的配置将在migrations路径下生成[当前时间-base.sql]
# 在文件中编写-- +migrate Up 和-- +migrate Down 内容
# sql-migrate status up /down 详细见--help
# 执行sql-migrate up将生成数据库文件hub.db(路径自定义)，作为库调用详见github...

# 2、SQLBoiler用于生成数据库优先的ORM工具,而不是code-first(gorm/gorp)
# 根据sqlboiler.toml(只能这种格式),执行sqlboiler sqlite3 --no-tests --wipe (--help),将生成对应的数据库ORM代码
# 确保sqlboiler安装 https://github.com/volatiletech/sqlboiler
go get github.com/volatiletech/sqlboiler/v4
go get github.com/volatiletech/sqlboiler-sqlite3


# 确保当前路径
cd "$(dirname $0)"

# run the migrations (creates testrun/TestSchema/hubdb)
go test -run='HubDB'

# make sure the sqlite file was created
test -f testdata/TestHubDB/hubdb || {
    echo 'file hub.db is missing'
    exit 1
}

# generate the models package
sqlboiler sqlite3 --wipe --no-tests

echo "all done. models updated!"
