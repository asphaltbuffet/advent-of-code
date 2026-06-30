-- Solution for Advent of Code 2024 day 9.
--
-- part_one(input) and part_two(input) each receive the puzzle input as a
-- string and must return a value (any type; it will be tostring'd).

-- Checksums exceed 2^32 and the runtime uses doubles, so format with %.0f to
-- avoid scientific notation in the result.
local function fmt(n)
  return string.format("%.0f", n)
end

local function expand(s)
  local blocks = {}
  local id = 0
  for i = 1, #s do
    local n = tonumber(s:sub(i, i))
    local val
    if (i - 1) % 2 == 0 then
      val = id
      id = id + 1
    else
      val = -1
    end
    for _ = 1, n do
      blocks[#blocks + 1] = val
    end
  end
  return blocks
end

function part_one(input)
  input = input:gsub("%s+$", "")
  local blocks = expand(input)
  local l, r = 1, #blocks
  while l < r do
    if blocks[l] ~= -1 then
      l = l + 1
    elseif blocks[r] == -1 then
      r = r - 1
    else
      blocks[l] = blocks[r]
      blocks[r] = -1
      l = l + 1
      r = r - 1
    end
  end

  local checksum = 0
  for i = 1, #blocks do
    if blocks[i] ~= -1 then
      checksum = checksum + (i - 1) * blocks[i]
    end
  end
  return fmt(checksum)
end

function part_two(input)
  input = input:gsub("%s+$", "")
  local s = input

  local files = {} -- { id, start, length } with 0-based start
  local frees = {} -- { start, length }
  local pos = 0
  local id = 0
  for i = 1, #s do
    local n = tonumber(s:sub(i, i))
    if (i - 1) % 2 == 0 then
      files[#files + 1] = { id = id, start = pos, length = n }
      id = id + 1
    elseif n > 0 then
      frees[#frees + 1] = { start = pos, length = n }
    end
    pos = pos + n
  end

  for i = #files, 1, -1 do
    local f = files[i]
    for _, g in ipairs(frees) do
      if g.start >= f.start then
        break
      end
      if g.length >= f.length then
        f.start = g.start
        g.start = g.start + f.length
        g.length = g.length - f.length
        break
      end
    end
  end

  local checksum = 0
  for _, f in ipairs(files) do
    for k = 0, f.length - 1 do
      checksum = checksum + (f.start + k) * f.id
    end
  end
  return fmt(checksum)
end
