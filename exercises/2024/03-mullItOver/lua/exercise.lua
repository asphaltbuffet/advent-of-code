-- Solution for Advent of Code 2024 day 3.
--
-- part_one(input) and part_two(input) each receive the puzzle input as a
-- string and must return a value (any type; it will be tostring'd).

-- Lua patterns lack alternation and bounded quantifiers, so we bound the
-- digit count manually with %d%d?%d? (1-3 digits).
local MUL = "mul%((%d%d?%d?),(%d%d?%d?)%)"

function part_one(input)
  local sum = 0
  for a, b in input:gmatch(MUL) do
    sum = sum + tonumber(a) * tonumber(b)
  end
  return sum
end

function part_two(input)
  local sum = 0
  local enabled = true
  local i = 1
  local n = #input
  while i <= n do
    if input:sub(i, i + 3) == "do()" then
      enabled = true
      i = i + 4
    elseif input:sub(i, i + 6) == "don't()" then
      enabled = false
      i = i + 7
    else
      local s, e, a, b = input:find("^" .. MUL, i)
      if s then
        if enabled then
          sum = sum + tonumber(a) * tonumber(b)
        end
        i = e + 1
      else
        i = i + 1
      end
    end
  end
  return sum
end
