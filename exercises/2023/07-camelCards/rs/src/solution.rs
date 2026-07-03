// Solution for Advent of Code 2023 day 7.
//
// Camel Cards ranks hands first by type (five of a kind down to high card) and
// then card by card. We turn each hand into a sortable key
// `(type_rank, [card_rank; 5])`; Rust derives the lexicographic ordering on that
// tuple, so ranking the hands is a single `sort`. Part two makes 'J' a joker
// that counts as the most frequent other card for typing but sorts weakest.

/// Type rank 0..=6: high card up to five of a kind, from the card counts.
fn type_rank(hand: &[u8; 5], joker: bool) -> u8 {
    let mut counts = [0u8; 15];
    for &c in hand {
        counts[c as usize] += 1;
    }

    let jokers = if joker { std::mem::take(&mut counts[0]) } else { 0 };

    // Sort the non-zero counts descending; jokers reinforce the largest group.
    let mut groups: Vec<u8> = counts.iter().copied().filter(|&n| n > 0).collect();
    groups.sort_unstable_by(|a, b| b.cmp(a));
    if let Some(top) = groups.first_mut() {
        *top += jokers;
    } else {
        groups.push(jokers); // all jokers -> five of a kind
    }

    match (groups.first(), groups.get(1)) {
        (Some(5), _) => 6,
        (Some(4), _) => 5,
        (Some(3), Some(2)) => 4,
        (Some(3), _) => 3,
        (Some(2), Some(2)) => 2,
        (Some(2), _) => 1,
        _ => 0,
    }
}

/// Map a card to its strength; with jokers, 'J' becomes the weakest (0).
fn card_rank(c: char, joker: bool) -> u8 {
    match c {
        'A' => 14,
        'K' => 13,
        'Q' => 12,
        'J' => if joker { 0 } else { 11 },
        'T' => 10,
        d => d as u8 - b'0',
    }
}

fn winnings(input: &str, joker: bool) -> u64 {
    let mut hands: Vec<((u8, [u8; 5]), u64)> = input
        .lines()
        .map(|line| {
            let (cards, bid) = line.split_once(' ').unwrap();
            let ranks: [u8; 5] = std::array::from_fn(|i| {
                card_rank(cards.as_bytes()[i] as char, joker)
            });
            ((type_rank(&ranks, joker), ranks), bid.parse().unwrap())
        })
        .collect();

    hands.sort_by(|a, b| a.0.cmp(&b.0));
    hands
        .iter()
        .enumerate()
        .map(|(i, (_, bid))| (i as u64 + 1) * bid)
        .sum()
}

pub fn part_one(input: &str) -> String {
    winnings(input, false).to_string()
}

pub fn part_two(input: &str) -> String {
    winnings(input, true).to_string()
}
