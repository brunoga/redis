-- rlock_script.lua
-- KEYS[1] - reader count.
-- KEYS[2] - writer flag.
-- ARGV[1] - expiration time in milliseconds.

-- Opportunistically set the expiration time for the reader count key.
redis.call('PEXPIRE', KEYS[1], ARGV[1])

if redis.call('EXISTS', KEYS[2]) == 1 then
    -- A writer is either waiting or running. We can not have the lock.
    return -1
end

-- Increment the reader count and return the new value.
return redis.call('INCR', KEYS[1])
