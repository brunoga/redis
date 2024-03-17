-- lock_script.lua
-- KEYS[1] - reader count.
-- KEYS[2] - writer count.
-- ARGV[1] - expiration time in milliseconds.

-- Opportunistically set the expiration time for the writer count key.
redis.call('PEXPIRE', KEYS[2], ARGV[1])

-- Get reader/writer counts.
local reader_count = redis.call('GET', KEYS[1])
local writer_count = redis.call('GET', KEYS[2])

if writer_count == false then
    -- There is no active or waiting writer.
    if reader_count == false then
        -- There are no active readers either. Acquire the write lock.
        --redis.log(redis.LOG_NOTICE, 'Acquiring write lock.')
        redis.call('SET', KEYS[2], 1, 'PX', ARGV[1])
        return 1 -- 1 writer.
    else
        -- There is at least 1 reader. Record intention to acquire write lock.
        redis.call('SET', KEYS[2], 0, 'PX', ARGV[1])
        return -1 
    end
end

-- GET returns a string. Convert it to a number.
writer_count = tonumber(writer_count)

-- There is an active or waiting writer.
if writer_count == 0 and reader_count == false then
    -- There is a waiting writer which might as well be us and there are no
    -- readers. Acquire the write lock.
    --redis.log(redis.LOG_NOTICE, 'Acquiring write lock.')
    redis.call('SET', KEYS[2], 1, 'PX', ARGV[1])
    return 1 -- 1 writer.
end
    
 -- There is an active writer or an intent to write and active readers. We can
 -- not acquire the write lock.
return -1
