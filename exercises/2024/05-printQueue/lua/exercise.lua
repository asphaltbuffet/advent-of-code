-- Solution for Advent of Code 2024 day 5.
--
-- part_one(input) and part_two(input) each receive the puzzle input as a
-- string and must return a value (any type; it will be tostring'd).

-- parse returns rules[a][b] = true (a must precede b) and a list of updates,
-- each update being a list of page numbers.
local function parse(input)
  input = input:gsub("\n+$", "")
  local rules_block, updates_block = input:match("^(.-)\n\n(.*)$")

  local rules = {}
  for a, b in rules_block:gmatch("(%d+)|(%d+)") do
    a, b = tonumber(a), tonumber(b)
    rules[a] = rules[a] or {}
    rules[a][b] = true
  end

  local updates = {}
  for line in updates_block:gmatch("[^\n]+") do
    local pages = {}
    for p in line:gmatch("%d+") do
      pages[#pages + 1] = tonumber(p)
    end
    updates[#updates + 1] = pages
  end

  return rules, updates
end

local function before(rules, a, b)
  return rules[a] ~= nil and rules[a][b] == true
end

local function ordered(rules, pages)
  for i = 1, #pages do
    for j = i + 1, #pages do
      if before(rules, pages[j], pages[i]) then
        return false
      end
    end
  end
  return true
end

function part_one(input)
  local rules, updates = parse(input)
  local sum = 0
  for _, pages in ipairs(updates) do
    if ordered(rules, pages) then
      sum = sum + pages[math.floor(#pages / 2) + 1]
    end
  end
  return sum
end

function part_two(input)
  local rules, updates = parse(input)
  local sum = 0
  for _, pages in ipairs(updates) do
    if not ordered(rules, pages) then
      table.sort(pages, function(a, b)
        return before(rules, a, b)
      end)
      sum = sum + pages[math.floor(#pages / 2) + 1]
    end
  end
  return sum
end
