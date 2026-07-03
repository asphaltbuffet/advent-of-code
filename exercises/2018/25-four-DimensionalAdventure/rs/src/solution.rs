// Solution for Advent of Code 2018 day 25.

type Point = [i32; 4];

fn parse(input: &str) -> Vec<Point> {
    input
        .lines()
        .filter_map(|line| {
            let mut it = line
                .split(',')
                .filter_map(|s| s.trim().parse::<i32>().ok());
            let p = [it.next()?, it.next()?, it.next()?, it.next()?];
            Some(p)
        })
        .collect()
}

fn manhattan(a: &Point, b: &Point) -> i32 {
    (0..4).map(|i| (a[i] - b[i]).abs()).sum()
}

fn find(parent: &mut [usize], mut x: usize) -> usize {
    while parent[x] != x {
        parent[x] = parent[parent[x]]; // path compression
        x = parent[x];
    }
    x
}

// part_one counts the constellations: connected components of points where any two
// within Manhattan distance 3 join the same group, via union-find.
pub fn part_one(input: &str) -> String {
    let pts = parse(input);
    let mut parent: Vec<usize> = (0..pts.len()).collect();

    for i in 0..pts.len() {
        for j in (i + 1)..pts.len() {
            if manhattan(&pts[i], &pts[j]) <= 3 {
                let (ri, rj) = (find(&mut parent, i), find(&mut parent, j));
                parent[ri] = rj;
            }
        }
    }

    let mut roots = std::collections::HashSet::new();
    for i in 0..pts.len() {
        let r = find(&mut parent, i);
        roots.insert(r);
    }
    roots.len().to_string()
}

// part_two is the day 25 finale: no second puzzle, only the closing message.
pub fn part_two(_input: &str) -> String {
    "Merry Christmas!".to_string()
}
