-- Solution for Advent of Code 2024 day 11.
--
-- part_one(input) and part_two(input) each receive the puzzle input as a
-- string and must return a value (any type; it will be tostring'd).

-- blink applies one transformation, returning a (b may be nil for one stone).
local function blink(n)
  if n == 0 then
    return 1, nil
  end
  local s = string.format("%.0f", n)
  local len = #s
  if len % 2 == 0 then
    local half = len / 2
    return tonumber(s:sub(1, half)), tonumber(s:sub(half + 1))
  end
  return n * 2024, nil
end

local function count_after(input, blinks)
  local stones = {}
  for tok in input:gmatch("%d+") do
    local n = tonumber(tok)
    stones[n] = (stones[n] or 0) + 1
  end

  for _ = 1, blinks do
    local nxt = {}
    for n, cnt in pairs(stones) do
      local a, b = blink(n)
      nxt[a] = (nxt[a] or 0) + cnt
      if b ~= nil then
        nxt[b] = (nxt[b] or 0) + cnt
      end
    end
    stones = nxt
  end

  local total = 0
  for _, cnt in pairs(stones) do
    total = total + cnt
  end
  -- Total exceeds 2^32 and the runtime uses doubles; format to avoid
  -- scientific notation in the result.
  return string.format("%.0f", total)
end

function part_one(input)
  return count_after(input, 25)
end

function part_two(input)
  return count_after(input, 75)
end
