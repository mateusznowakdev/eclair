use<_base.scad>

$fn = 25;

// board dimensions
bx = 44.45;
by = 34.93;

// wall thickness
wlr = 2.0; // left-right
wtb = 1.0; // top-bottom
wfr = 0.8; // front-rear

// tolerance
tol = 0.15;

module Base() {
  translate([0, 0, -wfr])
    linear_extrude(wfr)
      rsquare([bx+2*wlr+2*tol, by+2*wtb+2*tol], r=3, center=true);
}

module Snaps() {
  sx = 2;
  sy = 1;
  sz = 3.3;

  snaps = [
    [[-11.4,  by/2], 90],
    [[ 11.4,  by/2], 90],
    [[ bx/2,   3.5], 0],
    [[ bx/2,  -5.5], 0],
    [[ 11.4, -by/2], 270],
    [[-11.4, -by/2], 270],
    [[-bx/2,  -5.5], 180],
    [[-bx/2,   3.5], 180],
    // these two are joined to form a single solid support
    [[    0,   3.5], 0],
    [[    0,   3.5], 180],
  ];

  for (s = snaps)
    translate(s[0])
      rotate([0, 0, s[1]])
        // single snap
        translate([0, sx/2])
          rotate([90, 0, 0])
            linear_extrude(sx)
              polygon([[0, 0], [-sy, 0], [-sy, sz], [0, sz], [0.5, sz-0.5], [0, sz-1.0]]);
}

Base();
Snaps();
