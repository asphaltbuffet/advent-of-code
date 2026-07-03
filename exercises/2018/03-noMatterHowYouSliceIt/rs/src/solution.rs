// Solution for Advent of Code 2018 day 3.

struct Claim {
    id: i64,
    left: usize,
    top: usize,
    width: usize,
    height: usize,
}

/// Parse each claim by scanning its five integers, which is robust to
/// whitespace and delivery quirks in the input.
fn parse(input: &str) -> Vec<Claim> {
    input
        .lines()
        .filter(|l| !l.trim().is_empty())
        .map(|line| {
            let nums: Vec<i64> = line
                .split(|c: char| !c.is_ascii_digit() && c != '-')
                .filter(|s| !s.is_empty())
                .map(|s| s.parse().expect("integer"))
                .collect();
            assert_eq!(nums.len(), 5, "expected 5 numbers in {line:?}");
            Claim {
                id: nums[0],
                left: nums[1] as usize,
                top: nums[2] as usize,
                width: nums[3] as usize,
                height: nums[4] as usize,
            }
        })
        .collect()
}

const SIDE: usize = 1000;

/// Count how many claims cover each fabric cell in a flat grid.
fn coverage(claims: &[Claim]) -> Vec<u16> {
    let mut cover = vec![0u16; SIDE * SIDE];
    for c in claims {
        for y in c.top..c.top + c.height {
            for x in c.left..c.left + c.width {
                cover[y * SIDE + x] += 1;
            }
        }
    }
    cover
}

pub fn part_one(input: &str) -> String {
    let claims = parse(input);
    coverage(&claims)
        .iter()
        .filter(|&&n| n >= 2)
        .count()
        .to_string()
}

pub fn part_two(input: &str) -> String {
    let claims = parse(input);
    let cover = coverage(&claims);

    claims
        .iter()
        .find(|c| {
            (c.top..c.top + c.height).all(|y| {
                (c.left..c.left + c.width).all(|x| cover[y * SIDE + x] == 1)
            })
        })
        .map(|c| c.id.to_string())
        .expect("no non-overlapping claim found")
}
