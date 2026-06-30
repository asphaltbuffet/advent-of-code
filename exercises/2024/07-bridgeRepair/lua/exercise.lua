-- Solution for Advent of Code 2024 day 7.
--
-- part_one(input) and part_two(input) each receive the puzzle input as a
-- string and must return a value (any type; it will be tostring'd).

local function parse(input)
  local eqs = {}
  for line in input:gmatch("[^\r\n]+") do
    local target, rest = line:match("^(%d+):%s*(.+)$")
    local nums = {}
    for n in rest:gmatch("%d+") do
      nums[#nums + 1] = tonumber(n)
    end
    eqs[#eqs + 1] = { target = tonumber(target), nums = nums }
  end
  return eqs
end

-- concat returns a || b, e.g. 12 || 345 = 12345.
local function concat(a, b)
  local pow = 1
  local t = b
  while t > 0 do
    pow = pow * 10
    t = math.floor(t / 10)
  end
  return a * pow + b
end

-- solvable: can some operator placement over nums[i..] reach target from acc?
local function solvable(target, acc, nums, i, concat_op)
  if i > #nums then
    return acc == target
  end
  if acc > target then
    return false
  end
  local n = nums[i]
  if solvable(target, acc + n, nums, i + 1, concat_op) then
    return true
  end
  if solvable(target, acc * n, nums, i + 1, concat_op) then
    return true
  end
  if concat_op and solvable(target, concat(acc, n), nums, i + 1, concat_op) then
    return true
  end
  return false
end

local function calibrate(input, concat_op)
  local sum = 0
  for _, eq in ipairs(parse(input)) do
    if solvable(eq.target, eq.nums[1], eq.nums, 2, concat_op) then
      sum = sum + eq.target
    end
  end
  -- Values exceed 2^32 and the runtime uses doubles, so format explicitly to
  -- avoid tostring() emitting scientific notation (e.g. 2.48e+14).
  return string.format("%.0f", sum)
end

function part_one(input)
  return calibrate(input, false)
end

function part_two(input)
  return calibrate(input, true)
end
