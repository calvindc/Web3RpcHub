#!/bin/sh

# generate /db/privacymode_string.go
~/godev/src/golang.org/x/tools/cmd/stringer/stringer --type=PrivacyMode

# generate /db/role_string.go
~/godev/src/golang.org/x/tools/cmd/stringer/stringer --type=Role