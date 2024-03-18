-- refresh.lua
-- KEYS[1] - writer count (may be empty).
-- KEYS[2] - reader count (may be empty).
-- ARGV[1] - expiration time in milliseconds.

-- Resets expiration of writer counts if it exists.
if KEYS[1] ~= '' then
    redis.call('PEXPIRE', KEYS[1], ARGV[1])
end

-- Resets expiration of reader counts if it exists.
if KEYS[2] ~= '' then
    redis.call('PEXPIRE', KEYS[2], ARGV[1])
end

return 0