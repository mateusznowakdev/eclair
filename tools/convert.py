import sys

if len(sys.argv) != 3:
    print("usage: source.pbm source.txt", file=sys.stderr)
    sys.exit(1)

# get and process glyph data

with open(sys.argv[1], "rt") as glyph_file:
    lines = glyph_file.readlines()

assert lines[0].strip() == "P1"
if lines[1].startswith("#"):
    lines.pop(1)

t_width, t_height = lines[1].split(" ")
t_width = int(t_width)
t_height = int(t_height)

data = [0] * t_width
x = 0
y = 0

for line in lines[2:]:
    for value in line.strip():
        if value == "1":
            data[x] |= (1 << y)
        x += 1
        if x == t_width:
            x = 0
            y += 1

# get glyph widths

with open(sys.argv[2], "rt") as metadata_file:
    g_widths = [int(w) for w in metadata_file.read().splitlines()]

# build Go data

print("package display")
print()
print("var font = [][]uint16{")

d = 0
g = 32
for width in g_widths:
    print(f"/*  {chr(g)}  */ {{0x0000, ", end="")
    for n in data[d:d+width]:
        print(f"0x{n:04X}, ", end="")
    print(f"}},")
    d += width
    g += 1

print("}")
