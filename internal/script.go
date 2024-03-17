package internal

import (
	_ "embed"

	"github.com/redis/go-redis/v9"
)

//go:embed lua/lock_script.lua
var lockScript string

//go:embed lua/unlock_script.lua
var unlockScript string

//go:embed lua/rlock_script.lua
var rlockScript string

//go:embed lua/runlock_script.lua
var runlockScript string

//go:embed lua/refresh_script.lua
var refreshScript string

var (
	// LockScript is the Lua script for acquiring a lock.
	LockScript = redis.NewScript(lockScript)

	// UnlockScript is the Lua script for releasing a lock.
	UnlockScript = redis.NewScript(unlockScript)

	// RLockScript is the Lua script for acquiring a read lock.
	RLockScript = redis.NewScript(rlockScript)

	// RUnlockScript is the Lua script for releasing a read lock.
	RUnlockScript = redis.NewScript(runlockScript)

	// RefreshScript is the Lua script for refreshing a lock.
	RefreshScript = redis.NewScript(refreshScript)
)
