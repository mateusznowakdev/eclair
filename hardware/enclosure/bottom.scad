use<_base.scad>

$fn = 25;

// board dimensions
pX = 44.45;
pY = 34.93;

// wall thickness
tX = 2.0;
tY = 1.5;
tZ = 1.0;

// tolerance/clearance
c = 0.25;

// total height
z = 3.25 + tZ;

module mBase() {
  x = pX + 2 * tX;
  y = pY + 2 * tY;
  r = 3.5;

  translate([-tX, -tY])
    linear_extrude(tZ)
      rsquare([x, y], r);
}

module mSnaps() {
  snaps = [
    // clockwise
    [[0, 0, -90], (pX-22)/2, pY],
    [[0, 0, -90], (pX+22)/2, pY],
    [[0, 0, 180], pX,        pY*1/2],
    [[0, 0, 180], pX,        pY*1/10],
    [[0, 0, 90],  (pX+22)/2, 0],
    [[0, 0, 90],  (pX-22)/2, 0],
    [[0, 0, 0],   0,         pY*1/10],
    [[0, 0, 0],   0,         pY*1/2],
  ];

  for (s = snaps) {
    translate([s[1], s[2], tZ])
      rotate(s[0])
        mSnap();
  }
}

module mSnap() {
  d = 0.5;
  l = 3;
  h = 2.5;

  rotate([90, 0, 0])
    linear_extrude(l, center=true)
      polygon([[0, 0], [tZ, 0], [tZ, h+d+c], [0, h+d+c], [0, h+d], [-d, h], [0, h-d]]);
}

union() {
  mBase();
  mSnaps();
}
