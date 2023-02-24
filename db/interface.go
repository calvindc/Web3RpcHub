package db

import (
	"context"

	refs "github.com/ssbc/go-ssb-refs"
)

// note:一个表对应一个操作数据库的服务接口

// AliasesServiceg 管理alias:
type AliasesDBService interface {

	// 客户端通过验证后，注册alias
	Register(ctx context.Context, alias string, userFeed refs.FeedRef, signature []byte) error

	// 查询所有注册的aliases
	List(ctx context.Context) ([]Alias, error)

	// 删除某个alias
	Revoke(ctx context.Context, alias string) error

	// 通过alias查找所有相关信息
	Resolve(context.Context, string) (Alias, error)

	// 通过id查找alias
	GetByID(context.Context, int64) (Alias, error)
}
