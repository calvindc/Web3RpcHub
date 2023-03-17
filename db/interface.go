package db

import (
	"context"

	"github.com/calvindc/Web3RpcHub/internal/refs"
	"go.mindeco.de/http/auth"
)

// note:一个表对应一个操作数据库的服务接口

type HubConfig interface {
	GetPrivacyMode(context.Context) (PrivacyMode, error)
	SetPrivacyMode(context.Context, PrivacyMode) error
	GetDefaultLanguage(context.Context) (string, error)
	SetDefaultLanguage(context.Context, string) error
}

type AuthFallbackService interface {
	// 对接收到的用户名和密文明文进行检查登录名可能是别名或者ssb-id,有效返回用户id
	auth.Auther
	// SetPassword 创建或更新回滚此用户的登录密码
	SetPassword(_ context.Context, memberID int64, password string) error
	// CreateResetToken 返回一个token，该token可以通过SetPasswordWithToken来重置成员的密码。.
	CreateResetToken(_ context.Context, createdByMember, forMember int64) (string, error)
	// SetPasswordWithToken 使用CreateResetToken创建的token,并相应地更新该成员的密码.
	SetPasswordWithToken(_ context.Context, resetToken string, password string) error
}

// AliasesService 管理alias:
type AliasesService interface {
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

// AuthWithTokenService 对本网络协议的challenge/response的管理接口
type AuthWithTokenService interface {
	// 用于生成存储在cookie中的令牌
	CreateToken(ctx context.Context, memberID int64) (string, error)
	// CheckToken 检查成员的token是否合法, 返回成员id
	CheckToken(ctx context.Context, token string) (int64, error)
	// RemoveToken 删除一个token
	RemoveToken(ctx context.Context, token string) error
	// WipeTokensForMember 删除某个成员所持有的所所有token
	WipeTokensForMember(ctx context.Context, memberID int64) error
}

type MembersService interface {
	// Add 添加一个成员
	Add(_ context.Context, pubKey refs.FeedRef, r Role) (int64, error)
	// GetByID id查询成员
	GetByID(context.Context, int64) (Member, error)
	// GetByFeed feed查询成员
	GetByFeed(context.Context, refs.FeedRef) (Member, error)
	// List 查询所有成员.
	List(context.Context) ([]Member, error)
	// Count 人数统计.
	Count(context.Context) (uint, error)
	// RemoveFeed 删除 by feed.
	RemoveFeed(context.Context, refs.FeedRef) error
	// RemoveID 删除by id.
	RemoveID(context.Context, int64) error
	// SetRole 改变已经通过验证的成员的role状态，注意最后一个管理员的删除会导致问题.
	SetRole(context.Context, int64, Role) error
}

// DeniedKeysService 不允许进入hub的用户名单
type DeniedKeysService interface {
	// Add feed和其他成员的comment一起写入
	Add(ctx context.Context, ref refs.FeedRef, comment string) error
	// HasFeed 某feed是否在此名单中.
	HasFeed(context.Context, refs.FeedRef) bool
	// HasID 某用户id是否在此名单中.
	HasID(context.Context, int64) bool
	// GetByID
	GetByID(context.Context, int64) (ListEntry, error)
	// List
	List(context.Context) ([]ListEntry, error)
	// Count
	Count(context.Context) (uint, error)
	// RemoveFeed removes the feed from the list.
	RemoveFeed(context.Context, refs.FeedRef) error
	// RemoveID removes the feed for the ID from the list.
	RemoveID(context.Context, int64) error
}

//InvitesService 邀请码的管理(生成,使用)
type InvitesService interface {
	// Create邀请码,createdBy是管理员/主持人/如果隐私模式是open,则允许1号成员创建
	Create(ctx context.Context, createdBy int64) (string, error)
	// Consume 检查邀请码有效性,将newMember添加到聊天室,同时使邀请码无效.
	Consume(ctx context.Context, token string, newMember refs.FeedRef) (Invite, error)
	// GetByToken returns the Invite if one for that token exists, or an error
	GetByToken(ctx context.Context, token string) (Invite, error)
	// GetByToken returns the Invite if one for that ID exists, or an error
	GetByID(ctx context.Context, id int64) (Invite, error)
	// List returns a list of all the valid invites
	List(ctx context.Context) ([]Invite, error)
	// Count returns the total number of invites, optionally excluding inactive invites
	Count(ctx context.Context, onlyActive bool) (uint, error)
	// Revoke 作废一个邀请码.
	Revoke(ctx context.Context, id int64) error
}

// PinnedNoticesService 管理员下发通知给功能页面:软件更新/隐私政策/行为准则/
type PinnedNoticesService interface {
	// List 返回所有固定通知及其相应通知和语言的列表
	List(context.Context) (PinnedNotices, error)
	// Set 创建一个通知，一个通知对应多语言版本(多条记录)
	Set(ctx context.Context, name PinnedNoticeName, id int64) error
	// Get 通过PinnedNoticeName和language获取notice内容
	Get(ctx context.Context, name PinnedNoticeName, language string) (*Notice, error)
}

// NoticesService
type NoticesService interface {
	// GetByID id返回通知页面
	GetByID(context.Context, int64) (Notice, error)
	// Save 更新或在id为0时创建
	Save(context.Context, *Notice) error
	// RemoveID 删除page.
	RemoveID(context.Context, int64) error
}
