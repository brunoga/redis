-- refresh.lua
-- KEYS[1] - writer count.
-- KEYS[2] - reader count.
-- ARGV[1] - expiration time in milliseconds.

-- Resets expiration of writer and reader counts if they exist.
redis.call('PEXPIRE', KEYS[1], ARGV[1])
redis.call('PEXPIRE', KEYS[2], ARGV[1])

return 0