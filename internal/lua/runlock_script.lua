-- runlock_script.lua
-- KEYS[1] - reader count.

local result = redis.call('DECR', KEYS[1])

if result < 0 then
    -- Too many unlocks.
    redis.call('INCR', KEYS[1])
    return -1
end

if result == 0 then
    -- No more readers. Remove the reader count key.
    if redis.call('DEL', KEYS[1]) == 0 then
        -- An error happened. Return it.
        return redis.error_reply("Failed to remove reader count key.")
    end
end

-- Return number of readers after decrementing.
--redis.log(redis.LOG_NOTICE, 'Released read lock.')
return result
