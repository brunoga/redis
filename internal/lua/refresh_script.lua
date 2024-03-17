-- refresh.lua
-- KEYS[1] - writer flag.
-- KEYS[2] - reader count.
-- ARGV[1] - expiration time in milliseconds.

-- Resets expiration of writer flag and reader count if they exist.
redis.call('PEXPIRE', KEYS[1], ARGV[1])
redis.call('PEXPIRE', KEYS[2], ARGV[1])

return 1
