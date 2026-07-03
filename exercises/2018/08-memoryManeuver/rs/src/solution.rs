//! 2018 Day 8: Memory Maneuver

/// Reads one node starting at the front of `nums`, returning its metadata sum.
fn meta_sum(nums: &mut impl Iterator<Item = usize>) -> usize {
    let children = nums.next().unwrap();
    let meta = nums.next().unwrap();
    let child_sum: usize = (0..children).map(|_| meta_sum(nums)).sum();
    child_sum + (0..meta).map(|_| nums.next().unwrap()).sum::<usize>()
}

/// Reads one node, returning its "value": leaves sum their metadata, internal
/// nodes sum the values of the children their (1-based) metadata point at.
fn node_value(nums: &mut impl Iterator<Item = usize>) -> usize {
    let children = nums.next().unwrap();
    let meta = nums.next().unwrap();
    let child_values: Vec<usize> = (0..children).map(|_| node_value(nums)).collect();

    (0..meta)
        .map(|_| nums.next().unwrap())
        .map(|r| {
            if children == 0 {
                r
            } else {
                child_values.get(r.wrapping_sub(1)).copied().unwrap_or(0)
            }
        })
        .sum()
}

fn parse(input: &str) -> impl Iterator<Item = usize> + '_ {
    input.split_whitespace().map(|t| t.parse().unwrap())
}

pub fn part_one(input: &str) -> String {
    meta_sum(&mut parse(input)).to_string()
}

pub fn part_two(input: &str) -> String {
    node_value(&mut parse(input)).to_string()
}
