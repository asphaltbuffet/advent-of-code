-- Solution for Advent of Code 2024 day 12.
--
-- part_one(input) and part_two(input) each receive the puzzle input as a
-- string and must return a value (any type; it will be tostring'd).

local NEIGHBORS = { { -1, 0 }, { 1, 0 }, { 0, -1 }, { 0, 1 } }
local CORNERS = { { -1, -1 }, { -1, 1 }, { 1, -1 }, { 1, 1 } }

-- regions returns a list of { area, perimeter, sides } for each connected
-- same-letter region.
local function regions(input)
  local grid = {}
  for line in input:gmatch("[^\r\n]+") do
    grid[#grid + 1] = line
  end
  local rows = #grid

  local function at(r, c)
    local row = grid[r]
    if not row or c < 1 or c > #row then
      return nil
    end
    return row:sub(c, c)
  end

  local seen = {}
  local out = {}
  for sr = 1, rows do
    for sc = 1, #grid[sr] do
      local skey = sr * 1000 + sc
      if not seen[skey] then
        local letter = grid[sr]:sub(sc, sc)
        local area, perimeter, sides = 0, 0, 0
        local stack = { { sr, sc } }
        seen[skey] = true
        while #stack > 0 do
          local cur = stack[#stack]
          stack[#stack] = nil
          local r, c = cur[1], cur[2]
          area = area + 1

          local function same(nr, nc)
            return at(nr, nc) == letter
          end

          for _, d in ipairs(NEIGHBORS) do
            local nr, nc = r + d[1], c + d[2]
            if not same(nr, nc) then
              perimeter = perimeter + 1
            else
              local key = nr * 1000 + nc
              if not seen[key] then
                seen[key] = true
                stack[#stack + 1] = { nr, nc }
              end
            end
          end

          for _, q in ipairs(CORNERS) do
            local vert = same(r + q[1], c)
            local horiz = same(r, c + q[2])
            local diag = same(r + q[1], c + q[2])
            if not vert and not horiz then
              sides = sides + 1
            elseif vert and horiz and not diag then
              sides = sides + 1
            end
          end
        end
        out[#out + 1] = { area = area, perimeter = perimeter, sides = sides }
      end
    end
  end
  return out
end

function part_one(input)
  local total = 0
  for _, reg in ipairs(regions(input)) do
    total = total + reg.area * reg.perimeter
  end
  return total
end

function part_two(input)
  local total = 0
  for _, reg in ipairs(regions(input)) do
    total = total + reg.area * reg.sides
  end
  return total
end
