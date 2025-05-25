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

module Supports() {
  sx = 2.5;
  sy = 1;

  supports = [
    [[     -11.4,  by/2-sy/2], 0],
    [[      11.4,  by/2-sy/2], 0],
    [[-bx/2+sy/2,        3.5], 270],
    [[         0,        3.5], 270],
    [[ bx/2-sy/2,        3.5], 90],
    [[-bx/2+sy/2,       -5.5], 270],
    [[ bx/2-sy/2,       -5.5], 90],
    [[-bx/2+sy/2,      -14.0], 270],
    [[ bx/2-sy/2,      -14.0], 90],
    [[     -11.4, -by/2+sy/2], 180],
    [[       0.0, -by/2+sy/2], 180],
    [[      11.4, -by/2+sy/2], 180],
  ];

  for (s = supports)
    linear_extrude(3.3)
      translate(s[0])
        rotate([0, 0, s[1]])
          square([sx, sy], center=true);
}

Base();
Supports();

//module mSnap() {
//  d = 0.5;
//  l = 3;
//  h = 2.5;
//
//  rotate([90, 0, 0])
//    linear_extrude(l, center=true)
//      polygon([[0, 0], [tZ, 0], [tZ, h+d+c], [0, h+d+c], [0, h+d], [-d, h], [0, h-d]]);
//}
