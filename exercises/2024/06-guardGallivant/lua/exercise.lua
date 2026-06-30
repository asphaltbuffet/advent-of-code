-- Solution for Advent of Code 2024 day 6.
--
-- part_one(input) and part_two(input) each receive the puzzle input as a
-- string and must return a value (any type; it will be tostring'd).

-- Facing directions, ordered so (dir % 4) + 1 is a right turn (1-based dir).
local DR = { -1, 0, 1, 0 } -- up, right, down, left
local DC = { 0, 1, 0, -1 }

local function parse(input)
  local grid = {}
  for line in input:gmatch("[^\r\n]+") do
    grid[#grid + 1] = line
  end
  return grid
end

local function find_start(grid)
  for r = 1, #grid do
    local c = grid[r]:find("%^")
    if c then
      return r, c
    end
  end
end

-- walk simulates the guard. extra_key (or nil) is an added obstruction encoded
-- as r * 10000 + c. Returns (visited table keyed by r*10000+c, looped boolean).
local function walk(grid, sr, sc, extra_key)
  local rows = #grid
  local visited = {}
  local seen = {}
  local r, c, d = sr, sc, 1
  while true do
    visited[r * 10000 + c] = true
    local skey = (r * 10000 + c) * 4 + d
    if seen[skey] then
      return visited, true
    end
    seen[skey] = true

    local nr, nc = r + DR[d], c + DC[d]
    if nr < 1 or nr > rows or nc < 1 or nc > #grid[nr] then
      return visited, false
    end
    local ch = grid[nr]:sub(nc, nc)
    if ch == "#" or (nr * 10000 + nc) == extra_key then
      d = (d % 4) + 1
    else
      r, c = nr, nc
    end
  end
end

function part_one(input)
  local grid = parse(input)
  local sr, sc = find_start(grid)
  local visited = walk(grid, sr, sc, nil)
  local count = 0
  for _ in pairs(visited) do
    count = count + 1
  end
  return count
end

function part_two(input)
  local grid = parse(input)
  local sr, sc = find_start(grid)
  local path = walk(grid, sr, sc, nil)
  path[sr * 10000 + sc] = nil

  local count = 0
  for key in pairs(path) do
    local _, looped = walk(grid, sr, sc, key)
    if looped then
      count = count + 1
    end
  end
  return count
end
