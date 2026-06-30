-- Solution for Advent of Code 2024 day 10.
--
-- part_one(input) and part_two(input) each receive the puzzle input as a
-- string and must return a value (any type; it will be tostring'd).

local STEPS = { { -1, 0 }, { 1, 0 }, { 0, -1 }, { 0, 1 } }

local function parse(input)
  local grid = {}
  for line in input:gmatch("[^\r\n]+") do
    local row = {}
    for i = 1, #line do
      local ch = line:sub(i, i)
      row[i] = tonumber(ch) or -1
    end
    grid[#grid + 1] = row
  end
  return grid
end

local function at(grid, r, c)
  local row = grid[r]
  if not row or c < 1 or c > #row then
    return -1
  end
  return row[c]
end

local function reachable_9s(grid, r, c, ends)
  if grid[r][c] == 9 then
    ends[r * 1000 + c] = true
    return
  end
  for _, d in ipairs(STEPS) do
    local nr, nc = r + d[1], c + d[2]
    if at(grid, nr, nc) == grid[r][c] + 1 then
      reachable_9s(grid, nr, nc, ends)
    end
  end
end

function part_one(input)
  local grid = parse(input)
  local score = 0
  for r = 1, #grid do
    for c = 1, #grid[r] do
      if grid[r][c] == 0 then
        local ends = {}
        reachable_9s(grid, r, c, ends)
        for _ in pairs(ends) do
          score = score + 1
        end
      end
    end
  end
  return score
end

function part_two(input)
  local grid = parse(input)
  local memo = {}

  local function ratings(r, c)
    if grid[r][c] == 9 then
      return 1
    end
    local key = r * 1000 + c
    if memo[key] ~= nil then
      return memo[key]
    end
    local total = 0
    for _, d in ipairs(STEPS) do
      local nr, nc = r + d[1], c + d[2]
      if at(grid, nr, nc) == grid[r][c] + 1 then
        total = total + ratings(nr, nc)
      end
    end
    memo[key] = total
    return total
  end

  local rating = 0
  for r = 1, #grid do
    for c = 1, #grid[r] do
      if grid[r][c] == 0 then
        rating = rating + ratings(r, c)
      end
    end
  end
  return rating
end
