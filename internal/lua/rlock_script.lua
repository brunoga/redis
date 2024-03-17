-- rlock_script.lua
-- KEYS[1] - reader count.
-- KEYS[2] - writer flag.
-- ARGV[1] - expiration time in milliseconds.

if redis.call('EXISTS', KEYS[2]) == 1 then
    -- A writer is either waiting or running. We can not have the lock.
    return -1
end

local reader_result = redis.call('INCR', KEYS[1])

if redis.call('PEXPIRE', KEYS[1], ARGV[1]) == 0 then
    -- An error happened. Return it.
    return redis.error_reply(
        "Failed to set expiration time for reader count key.")
end

-- We acquired the lock.
return reader_result