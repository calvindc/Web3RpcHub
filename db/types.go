package db

import (
	"errors"

	"fmt"

	"database/sql/driver"

	"time"

	"sort"

	"github.com/calvindc/Web3RpcHub/refs"
)

// ErrNotFound
var ErrNotFound = errors.New("db: object not found")

//================================================================================================
// Alias db存储格式.
type Alias struct {
	ID        int64
	Name      string
	Feed      refs.FeedRef
	Signature []byte
}

// ErrAliasTaken
type ErrAliasTaken struct {
	Name string
}

func (e ErrAliasTaken) Error() string {
	return fmt.Sprintf("alias (%q) is already taken", e.Name)
}

//================================================================================================
// Member holds all the information an internal user of the hub has.
type Member struct {
	ID      int64
	Role    Role
	PubKey  refs.FeedRef
	Aliases []Alias
}

//================================================================================================
// PrivacyMode 设定的隐私模式
type PrivacyMode uint

// PrivacyMode describes the access mode the hub server is currently running under.
// ModeOpen 允许任何人创建聊天室邀请
// ModeCommunity 只能邀请已经存在的hub成员 (i.e. "internal users")
// ModeRestricted 只允许管理员和支持人创建邀请
const (
	ModeUnknown PrivacyMode = iota
	ModeOpen
	ModeCommunity
	ModeRestricted
)

// IsValid
func (pm PrivacyMode) IsValid() error {
	if pm == ModeUnknown || pm > ModeRestricted {
		return errors.New("No such privacy mode")
	}
	return nil
}

// PrivacyMode 对PrivacyMode的SQL封装处理接口
func (pm *PrivacyMode) Scan(src interface{}) error {
	dbValue, ok := src.(int64)
	if !ok {
		return fmt.Errorf("unexpected type: %T", src)
	}

	privacyMode := PrivacyMode(dbValue)
	err := privacyMode.IsValid()
	if err != nil {
		return err
	}

	*pm = privacyMode
	return nil
}

func (pm PrivacyMode) Value() (driver.Value, error) {
	return driver.Value(int64(pm)), nil
}

//================================================================================================
//角色描述内部用户（或成员）的授权级别, 有效角色为成员/主持人/管理员
type Role uint

const (
	RoleUnknown Role = iota
	RoleMember
	RoleModerator
	RoleAdmin
)

func (r Role) IsValid() error {
	if r == RoleUnknown {
		return errors.New("unknown member role")
	}
	if r > RoleAdmin {
		return errors.New("invalid member role")
	}
	return nil
}

//使用stringer来转换类型 golang.org/x/tools/cmd/stringer/stringer --type=Role

var (
	roleAdminString  = RoleAdmin.String()
	roleModString    = RoleModerator.String()
	roleMemberString = RoleMember.String()
)

// UnmarshalText checks if a string is a valid role
func (r *Role) UnmarshalText(text []byte) error {
	roleStr := string(text)
	switch roleStr {
	case roleAdminString:
		*r = RoleAdmin
	case roleModString:
		*r = RoleModerator
	case roleMemberString:
		*r = RoleMember
	default:
		return fmt.Errorf("unknown member role: %q", roleStr)
	}

	return nil
}

//================================================================================================
// ErrAlreadyAdded
type ErrAlreadyAdded struct {
	Ref refs.FeedRef
}

// Error
func (aa ErrAlreadyAdded) Error() string {
	return fmt.Sprintf("db: the item (%s) is already on the list", aa.Ref.PubKey())
}

//================================================================================================
// Invite
type Invite struct {
	ID        int64
	CreatedBy Member
	CreatedAt time.Time
}

// ListEntry values are returned by the DenyListServices
type ListEntry struct {
	ID     int64
	PubKey refs.FeedRef

	CreatedAt time.Time
	Comment   string
}

//================================================================================================
// DBFeedRef wraps a feed reference and implements the SQL marshaling interfaces.
type DBFeedRef struct{ refs.FeedRef }

// Scan
func (r *DBFeedRef) Scan(src interface{}) error {
	str, ok := src.(string)
	if !ok {
		return fmt.Errorf("unexpected type: %T", src)
	}

	fr, err := refs.ParseFeedRef(str)
	if err != nil {
		return err
	}

	r.FeedRef = fr
	return nil
}

// Value
func (r DBFeedRef) Value() (driver.Value, error) {
	return driver.Value(r.String()), nil
}

//================================================================================================
// PinnedNoticeName
type PinnedNoticeName string

func (n PinnedNoticeName) String() string {
	return string(n)
}

// These are the well known names that the room page will display
const (
	NoticeDescription   PinnedNoticeName = "NoticeDescription"
	NoticeNews          PinnedNoticeName = "NoticeNews"
	NoticePrivacyPolicy PinnedNoticeName = "NoticePrivacyPolicy"
	NoticeCodeOfConduct PinnedNoticeName = "NoticeCodeOfConduct"
)

// Valid returns true if the page name is well known.
func (n PinnedNoticeName) Valid() bool {
	return n == NoticeNews ||
		n == NoticeDescription ||
		n == NoticePrivacyPolicy ||
		n == NoticeCodeOfConduct
}

type PinnedNotices map[PinnedNoticeName][]Notice

// Notice holds the title and content of a page that is user generated
type Notice struct {
	Time     int64
	ID       int64
	Title    string
	Content  string
	Language string
	Picture  []byte
}

type PinnedNotice struct {
	Name    PinnedNoticeName
	Notices []Notice
}

type SortedPinnedNotices []PinnedNotice

// Sorted returns a sorted list of the map, by the key names
func (pn PinnedNotices) Sorted() SortedPinnedNotices {
	lst := make(SortedPinnedNotices, 0, len(pn))
	for name, notices := range pn {
		lst = append(lst, PinnedNotice{
			Name:    name,
			Notices: notices,
		})
	}
	sort.Sort(lst)
	return lst
}

var _ sort.Interface = (SortedPinnedNotices)(nil)

func (byName SortedPinnedNotices) Len() int { return len(byName) }

func (byName SortedPinnedNotices) Less(i, j int) bool {
	return byName[i].Name < byName[j].Name
}

func (byName SortedPinnedNotices) Swap(i, j int) {
	byName[i], byName[j] = byName[j], byName[i]
}
