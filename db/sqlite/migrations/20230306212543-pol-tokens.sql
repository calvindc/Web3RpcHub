-- +migrate Up

-- SIWWEB3R stands for sign-in with WEB3R
CREATE TABLE SIWWEB3R_sessions (
  id            INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  token         TEXT UNIQUE NOT NULL,
  member_id     INTEGER NOT NULL,
  created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY ( member_id ) REFERENCES members( "id" ) ON DELETE CASCADE
);
CREATE UNIQUE INDEX SIWWEB3R_by_token ON SIWWEB3R_sessions(token);
CREATE INDEX SIWWEB3R_by_member ON SIWWEB3R_sessions(member_id);



-- +migrate Down
DROP INDEX SIWWEB3R_by_token;
DROP INDEX SIWWEB3R_by_member;
DROP TABLE SIWWEB3R_sessions;