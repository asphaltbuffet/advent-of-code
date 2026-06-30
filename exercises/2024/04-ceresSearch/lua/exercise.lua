-- Solution for Advent of Code 2024 day 4.
--
-- part_one(input) and part_two(input) each receive the puzzle input as a
-- string and must return a value (any type; it will be tostring'd).

local function grid(input)
  local g = {}
  for line in input:gmatch("[^\r\n]+") do
    g[#g + 1] = line
  end
  return g
end

-- 1-based char lookup; returns nil if out of bounds.
local function at(g, r, c)
  local row = g[r]
  if not row or c < 1 or c > #row then
    return nil
  end
  return row:sub(c, c)
end

local DIRS = {
  { 0, 1 }, { 0, -1 }, { 1, 0 }, { -1, 0 },
  { 1, 1 }, { 1, -1 }, { -1, 1 }, { -1, -1 },
}
local WORD = "XMAS"

function part_one(input)
  local g = grid(input)
  local count = 0
  for r = 1, #g do
    for c = 1, #g[r] do
      if at(g, r, c) == "X" then
        for _, d in ipairs(DIRS) do
          local ok = true
          for k = 1, #WORD - 1 do
            if at(g, r + d[1] * k, c + d[2] * k) ~= WORD:sub(k + 1, k + 1) then
              ok = false
              break
            end
          end
          if ok then
            count = count + 1
          end
        end
      end
    end
  end
  return count
end

local function is_mas(a, b)
  return (a == "M" and b == "S") or (a == "S" and b == "M")
end

function part_two(input)
  local g = grid(input)
  local count = 0
  for r = 2, #g - 1 do
    for c = 2, #g[r] - 1 do
      if at(g, r, c) == "A"
        and is_mas(at(g, r - 1, c - 1), at(g, r + 1, c + 1))
        and is_mas(at(g, r - 1, c + 1), at(g, r + 1, c - 1)) then
        count = count + 1
      end
    end
  end
  return count
end
