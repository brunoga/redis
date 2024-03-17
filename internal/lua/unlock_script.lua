-- unlock_script.lua
-- KEYS[1] - writer count.

local writer_count = redis.call('GET', KEYS[1])
if writer_count == false or tonumber(writer_count) == 0 then
    -- There is no active writer.
    return -1
end

if redis.call("DEL", KEYS[1]) == 0 then
    -- An error happened. Return it.
    return redis.error_reply("Failed to remove writer count key.")
end

return 0 -- 0 writers.
