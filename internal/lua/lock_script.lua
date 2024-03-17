-- lock_script.lua
-- KEYS[1] - reader count.
-- KEYS[2] - writer count.
-- ARGV[1] - expiration time in milliseconds.

-- Get reader/writer counts.
local writer_count = redis.call('GET', KEYS[1])
local reader_count = redis.call('GET', KEYS[2])


if writer_count == nil then
    -- There is no active or waiting writer.
    if reader_count == nil then
        -- There are no active readers either. Acquire the write lock.
        redis.call('SET', KEYS[2], 1, 'PX', ARGV[1])
        return 1 -- 1 writer.
    else
        -- There is at least 1 reader. Record intention to acquire write lock.
        redis.call('SET', KEYS[2], 0, 'PX', ARGV[1])
        return -1 
    end
end

-- There is an active or waiting writer.
if writer_count == 0 then
    -- There is a waiting writer which might as well be us. Acquire the
    -- write lock.
    redis.call('SET', KEYS[2], 1, 'PX', ARGV[1])
    return 1 -- 1 writer.
end
    
-- There is an active writer. We can not acquire the write lock.
return -1
