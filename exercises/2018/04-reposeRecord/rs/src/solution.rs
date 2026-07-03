use std::collections::HashMap;

/// Sort the shuffled log (the timestamp format sorts lexically) and tally, per
/// guard, how many days they were asleep at each minute 0..59.
/// `asleep[guard][minute]` = number of days that guard was asleep then.
fn sleep_by_guard(input: &str) -> HashMap<u32, [u32; 60]> {
    let mut lines: Vec<&str> = input.lines().map(str::trim).filter(|l| !l.is_empty()).collect();
    lines.sort_unstable();

    let mut asleep: HashMap<u32, [u32; 60]> = HashMap::new();
    let mut guard = 0u32;
    let mut start = 0usize;

    for line in lines {
        // "[YYYY-MM-DD HH:MM] ..." — the minute is at bytes 15..17.
        let minute: usize = line[15..17].parse().unwrap();

        if let Some((_, rest)) = line.split_once('#') {
            // "Guard #NN begins shift" — the ID is the only # number.
            guard = rest.split_whitespace().next().unwrap().parse().unwrap();
        } else if line.ends_with("falls asleep") {
            start = minute;
        } else if line.ends_with("wakes up") {
            let row = asleep.entry(guard).or_insert([0; 60]);
            row[start..minute].iter_mut().for_each(|n| *n += 1);
        }
    }

    asleep
}

pub fn part_one(input: &str) -> String {
    let asleep = sleep_by_guard(input);

    // Strategy 1: guard with the most total minutes asleep.
    let (&guard, mins) = asleep
        .iter()
        .max_by_key(|(_, mins)| mins.iter().sum::<u32>())
        .unwrap();

    // The minute that guard was most often asleep.
    let minute = (0..60).max_by_key(|&m| mins[m]).unwrap();

    (guard as usize * minute).to_string()
}

pub fn part_two(input: &str) -> String {
    let asleep = sleep_by_guard(input);

    // Strategy 2: the single (guard, minute) with the highest sleep count.
    let (guard, minute, _) = asleep
        .iter()
        .flat_map(|(&g, mins)| (0..60).map(move |m| (g, m, mins[m])))
        .max_by_key(|&(_, _, count)| count)
        .unwrap();

    (guard as usize * minute).to_string()
}
