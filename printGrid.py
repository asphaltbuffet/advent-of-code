import argparse


def print_grid(x, y, s="\t"):
    for row in range(y):
        for col in range(x):
            print(f"[{col}, {row}]{s}", end="")
        print()


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Print a grid with specified dimensions."
    )
    parser.add_argument("x", type=int, help="Width of the grid.")
    parser.add_argument("y", type=int, help="Height of the grid.")
    parser.add_argument(
        "-s", "--string", type=str, default="\t", help="An optional separator."
    )
    args = parser.parse_args()

    print_grid(args.x, args.y, args.string)
