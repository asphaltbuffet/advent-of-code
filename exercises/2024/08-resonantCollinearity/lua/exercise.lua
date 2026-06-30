-- Solution for Advent of Code 2024 day 8.
--
-- part_one(input) and part_two(input) each receive the puzzle input as a
-- string and must return a value (any type; it will be tostring'd).

local function parse(input)
  local grid = {}
  for line in input:gmatch("[^\r\n]+") do
    grid[#grid + 1] = line
  end
  local rows = #grid
  local cols = rows > 0 and #grid[1] or 0
  local antennas = {}
  for r = 1, rows do
    for c = 1, cols do
      local ch = grid[r]:sub(c, c)
      if ch ~= "." then
        antennas[ch] = antennas[ch] or {}
        local list = antennas[ch]
        list[#list + 1] = { r = r, c = c }
      end
    end
  end
  return rows, cols, antennas
end

local function count(set)
  local n = 0
  for _ in pairs(set) do
    n = n + 1
  end
  return n
end

function part_one(input)
  local rows, cols, antennas = parse(input)
  local function in_bounds(r, c)
    return r >= 1 and r <= rows and c >= 1 and c <= cols
  end

  local nodes = {}
  for _, pts in pairs(antennas) do
    for i = 1, #pts do
      for j = i + 1, #pts do
        local a, b = pts[i], pts[j]
        local dr, dc = b.r - a.r, b.c - a.c
        local n1r, n1c = a.r - dr, a.c - dc
        local n2r, n2c = b.r + dr, b.c + dc
        if in_bounds(n1r, n1c) then
          nodes[n1r * cols + n1c] = true
        end
        if in_bounds(n2r, n2c) then
          nodes[n2r * cols + n2c] = true
        end
      end
    end
  end
  return count(nodes)
end

function part_two(input)
  local rows, cols, antennas = parse(input)
  local function in_bounds(r, c)
    return r >= 1 and r <= rows and c >= 1 and c <= cols
  end

  local nodes = {}
  for _, pts in pairs(antennas) do
    for i = 1, #pts do
      for j = i + 1, #pts do
        local a, b = pts[i], pts[j]
        local dr, dc = b.r - a.r, b.c - a.c
        for _, step in ipairs({ -1, 1 }) do
          local r, c = a.r, a.c
          while in_bounds(r, c) do
            nodes[r * cols + c] = true
            r, c = r + step * dr, c + step * dc
          end
        end
      end
    end
  end
  return count(nodes)
end
