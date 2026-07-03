// Solution for Advent of Code 2025 day 8.
//
// Junctions are 3D points. We take every pairwise squared distance as a candidate
// wire and sort shortest first. Part one connects the k shortest wires (k depends
// on the count of junctions) and multiplies the three largest resulting circuit
// sizes. Part two runs Kruskal's MST: connect the shortest wires that merge
// components until the graph is one tree, and report the product of the last
// merging junctions' x-coordinates. A union-find drives both.

struct Dsu {
    parent: Vec<usize>,
}

impl Dsu {
    fn new(n: usize) -> Self {
        Dsu { parent: (0..n).collect() }
    }

    fn find(&mut self, mut x: usize) -> usize {
        while self.parent[x] != x {
            self.parent[x] = self.parent[self.parent[x]];
            x = self.parent[x];
        }
        x
    }

    /// Merge the sets of `a` and `b`; returns whether they were distinct.
    fn union(&mut self, a: usize, b: usize) -> bool {
        let (ra, rb) = (self.find(a), self.find(b));
        if ra == rb {
            return false;
        }
        self.parent[ra] = rb;
        true
    }
}

fn junctions(input: &str) -> Vec<[i64; 3]> {
    input
        .lines()
        .map(|line| {
            let mut it = line.split(',').map(|t| t.parse().unwrap());
            [it.next().unwrap(), it.next().unwrap(), it.next().unwrap()]
        })
        .collect()
}

/// (squared distance, i, j) for every pair, sorted shortest first.
fn edges(pts: &[[i64; 3]]) -> Vec<(i64, usize, usize)> {
    let mut edges = Vec::with_capacity(pts.len() * (pts.len() - 1) / 2);
    for i in 0..pts.len() {
        for j in 0..i {
            let d = (0..3).map(|k| (pts[i][k] - pts[j][k]).pow(2)).sum();
            edges.push((d, i, j));
        }
    }
    edges.sort_unstable();
    edges
}

pub fn part_one(input: &str) -> String {
    let pts = junctions(input);
    let wires = if pts.len() < 100 { 10 } else { 1000 };

    let mut dsu = Dsu::new(pts.len());
    for &(_, a, b) in edges(&pts).iter().take(wires) {
        dsu.union(a, b);
    }

    let mut sizes = vec![0usize; pts.len()];
    for i in 0..pts.len() {
        let root = dsu.find(i);
        sizes[root] += 1;
    }
    sizes.sort_unstable_by(|a, b| b.cmp(a));
    (sizes[0] * sizes[1] * sizes[2]).to_string()
}

pub fn part_two(input: &str) -> String {
    let pts = junctions(input);
    let mut dsu = Dsu::new(pts.len());

    let mut connections = 0;
    for (_, a, b) in edges(&pts) {
        if dsu.union(a, b) {
            connections += 1;
            if connections == pts.len() - 1 {
                return (pts[a][0] * pts[b][0]).to_string();
            }
        }
    }
    "-1".to_string()
}
