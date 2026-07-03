// Solution for Advent of Code 2025 day 6.
//
// The same grid is read two ways. Part one treats whitespace-delimited tokens as
// a column table — the last row is the operators, and each column's numbers are
// combined by that column's operator. Part two instead reads digits stacked
// vertically in each character column, right to left, with operator columns
// delimiting problems and a blank column separating them.

/// Combine a problem's numbers with its operator.
fn apply(op: u8, nums: &[i64]) -> i64 {
    match op {
        b'+' => nums.iter().sum(),
        _ => nums.iter().product(),
    }
}

pub fn part_one(input: &str) -> String {
    let mut lines: Vec<&str> = input.trim_end_matches('\n').lines().collect();
    let op_row = lines.pop().unwrap();
    let ops: Vec<u8> = op_row.split_whitespace().map(|s| s.as_bytes()[0]).collect();

    let mut columns: Vec<Vec<i64>> = vec![Vec::new(); ops.len()];
    for row in lines {
        for (i, tok) in row.split_whitespace().enumerate() {
            columns[i].push(tok.parse().unwrap());
        }
    }

    ops.iter()
        .zip(&columns)
        .map(|(&op, col)| apply(op, col))
        .sum::<i64>()
        .to_string()
}

pub fn part_two(input: &str) -> String {
    let lines: Vec<&[u8]> = input.trim_end_matches('\n').lines().map(str::as_bytes).collect();
    let (op_row, number_rows) = lines.split_last().unwrap();
    let width = lines.iter().map(|l| l.len()).max().unwrap();

    // Byte at (row, col), treating short (right-trimmed) rows as space-padded.
    let at = |row: &[u8], col: usize| -> u8 { row.get(col).copied().unwrap_or(b' ') };

    let mut total: i64 = 0;
    let mut nums: Vec<i64> = Vec::new();
    for col in (0..width).rev() {
        let mut n: i64 = 0;
        let mut has_digit = false;
        for row in number_rows {
            let c = at(row, col);
            if c != b' ' {
                n = n * 10 + (c - b'0') as i64;
                has_digit = true;
            }
        }
        if has_digit {
            nums.push(n);
        }
        let op = at(op_row, col);
        if op != b' ' {
            total += apply(op, &nums);
            nums.clear();
        }
    }
    total.to_string()
}
