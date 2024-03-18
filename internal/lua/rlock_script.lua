-- rlock_script.lua
-- KEYS[1] - reader count.
-- KEYS[2] - writer flag.
-- ARGV[1] - expiration time in milliseconds.

if redis.call('EXISTS', KEYS[2]) == 1 then
    -- A writer is either waiting or running. We can not have the lock.
    return -1
end

-- Increment the reader count and return the new value.
--redis.log(redis.LOG_NOTICE, 'Acquiring read lock.')
return redis.call('INCR', KEYS[1])
