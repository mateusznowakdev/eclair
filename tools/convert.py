#!/usr/bin/env python3

import sys

if len(sys.argv) != 3:
    print("usage: source.pgm source.txt", file=sys.stderr)
    sys.exit(1)

# parse metadata

with open(sys.argv[1], "rt") as sprite_file:
    lines = sprite_file.readlines()

assert lines[0].strip() == "P2"
if lines[1].startswith("#"):
    lines.pop(1)

t_width, t_height = lines[1].strip().split(" ")
t_width = int(t_width)
t_height = int(t_height)
assert t_height in (8, 16)

t_max_color = int(lines[2].strip())

# parse pixel data

sprite_data = [0] * t_width
mask_data = [0] * t_width
x = 0
y = 0

for line in lines[3:]:
    for value in line.strip().split():
        value = int(value)
        if value <= t_max_color * 1 / 3:
            sprite_data[x] |= 1 << y
        if value <= t_max_color * 2 / 3:
            mask_data[x] |= 1 << y

        x += 1
        if x == t_width:
            x = 0
            y += 1

# parse sprite widths

with open(sys.argv[2], "rt") as metadata_file:
    sprite_widths = [int(w) for w in metadata_file.read().splitlines()]

# handle tall (16px) spritesheets

if t_height == 16:
    sprite_widths *= 2
    sprite_data = [(s & 0xFF) for s in sprite_data] + [(s >> 8) for s in sprite_data]
    mask_data = [(m & 0xFF) for m in mask_data] + [(m >> 8) for m in mask_data]

# build Go data for spritesheet

print("var spritesheet = [][]uint8{")

d = 0
for width in sprite_widths:
    print("{", end="")
    for n in sprite_data[d : d + width]:
        print(f"0x{n:02X}, ", end="")
    print("},")
    d += width

print("}")

# build Go data for mask sheet

print("var masksheet = [][]uint8{")

d = 0
for width in sprite_widths:
    print("{", end="")
    for n in mask_data[d : d + width]:
        print(f"0x{n:02X}, ", end="")
    print("},")
    d += width

print("}")
