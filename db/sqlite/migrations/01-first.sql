-- +migrate Up
-- 内部用户表（使用邀请的人）
CREATE TABLE members (
  id            INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  role          INTEGER NOT NULL, --一个hub中的成员,主持人或管理员
  pub_key       TEXT    NOT NULL UNIQUE,

  CHECK(role > 0)
);
CREATE INDEX members_pubkeys ON members(pub_key);

--成员的密码登录（以防他们因任何原因无法使用ssb登录）member_id级联members中的id（删除）
CREATE TABLE fallback_passwords (
  id            INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  login         TEXT    NOT NULL UNIQUE,
  password_hash BLOB    NOT NULL,

  member_id     INTEGER NOT NULL,

  FOREIGN KEY ( member_id ) REFERENCES members( "id" )  ON DELETE CASCADE
);
CREATE INDEX fallback_passwords_by_login ON fallback_passwords(login);

-- 邀请码，每个成员的邀请码是一次性的
CREATE TABLE invites (
  id               INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  hashed_token     TEXT UNIQUE NOT NULL,
  created_by       INTEGER NOT NULL,
  created_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  active boolean   NOT NULL DEFAULT TRUE,

  FOREIGN KEY ( created_by ) REFERENCES members( "id" )  ON DELETE CASCADE
);
CREATE INDEX invite_active_ids ON invites(id) WHERE active=TRUE;
CREATE UNIQUE INDEX invite_active_tokens ON invites(hashed_token) WHERE active=TRUE;
CREATE INDEX invite_inactive ON invites(active);

-- 别名表
CREATE TABLE aliases (
  id            INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  name          TEXT UNIQUE NOT NULL,
  member_id     INTEGER NOT NULL,
  signature     BLOB NOT NULL,

  FOREIGN KEY ( member_id ) REFERENCES members( "id" )  ON DELETE CASCADE
);
CREATE UNIQUE INDEX aliases_ids ON aliases(id);
CREATE UNIQUE INDEX aliases_names ON aliases(name);

-- 被block不能进入hub的公钥(账户id)
CREATE TABLE denied_keys (
  id          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  pub_key     TEXT NOT NULL UNIQUE,
  comment     TEXT NOT NULL,
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX denied_keys_by_pubkey ON invites(active);



-- +migrate Down
DROP TABLE members;

DROP INDEX fallback_passwords_by_login;
DROP TABLE fallback_passwords;

DROP INDEX invite_active_ids;
DROP INDEX invite_active_tokens;
DROP INDEX invite_inactive;
DROP TABLE invites;

DROP INDEX aliases_ids;
DROP INDEX aliases_names;
DROP TABLE aliases;

DROP INDEX denied_keys_by_pubkey;
DROP TABLE denied_keys;