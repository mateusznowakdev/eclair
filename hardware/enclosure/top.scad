module rsquare(dim, r=0, center=false) {
  x = dim[0];
  y = dim[1];

  ox = center ? -x/2 : 0;
  oy = center ? -y/2 : 0;

  translate([ox, oy]) {
    union() {
      translate([r, 0])
        square([x-2*r, y]);
      translate([0, r])
        square([x, y-2*r]);
      translate([r, r])
        circle(r);
      translate([x-r, r])
        circle(r);
      translate([r, y-r])
        circle(r);
      translate([x-r, y-r])
        circle(r);
    }
  }
}

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

// total height:
// - buttons: 1.8mm
// -     PCB: 0.8mm
// -   USB-C: 3.2mm
z = 5.8;

module Base() {
  translate([0, 0, -z])
    linear_extrude(z + wfr)
      rsquare([bx + 2*wlr + 2*tol, by + 2*wtb + 2*tol], r=3, center=true);
}

module Extrude() {
  translate([0, 0, -z])
    linear_extrude(z)
      rsquare([bx + 2*tol, by + 2*tol], r=2-tol, center=true);
}

module ExtrudeSlideSwitch() {
  translate([-bx/2 - 0.5 - tol, by/4, -z])
    linear_extrude(z)
      square([1, by/2 - 4 + 2*tol], center=true);
}

module CutoutDisplay() {
  translate([0, 12])
    linear_extrude(wfr + tol) // add extra length to prevent rendering issues
      rsquare([25, 8], r=1, center=true);
}

module CutoutButtons() {
  buttons = [
    [-2,    0],
    [-1,    0],
    [ 0,    0],
    [ 1,    0],
    [ 2,    0],
    [-2,   -1],
    [-1,   -1],
    [ 0,   -1],
    [ 1,   -1],
    [ 2,   -1],
    [-1.5, -2],
    [-0.5, -2],
    [ 0.5, -2],
    [ 1.5, -2],
  ];

  translate([0, 1.6])
    linear_extrude(wfr + tol)
      for (b = buttons)
        translate([b[0] * 8.89, b[1] * 7.62])
          circle(2);
}

module CutoutSlideSwitch() {
  // top and bottom positions are is 12mm and 10.5mm away from the center
  // the main body of the switch is 4.5mm x 3mm, with 1.5mm travel
  translate([-bx/2-tol, (12+10.5)/2, -z/2])
    rotate([-90, 0, 90])
      linear_extrude(wlr + tol)
        rsquare([4.5 + 1.5 + 2*tol, 3 + 2*tol], r=1.5, center=true);
}

module CutoutUSB() {
  // slightly larger than in the datasheet
  w = 9;
  h = 3.2;

  translate([-w/2, by/2, -z+h])
    rotate([-90, 0, 0])
      linear_extrude(wtb + tol)
        square([w, h]);
}

module CutoutCharm() {
  // pads are 9mm and 11mm away from the center, but this would be too fragile
  translate([bx/2+tol, -8.5, -z+2])
    rotate([0, 90, 0])
      linear_extrude(wlr + tol) {
        circle(0.75);
        translate([0, -3, 0]) circle(0.75);
      }
}

module Supports() {
  // Y dimension between buttons, see CutoutButtons()
  middle = 1.6-7.62/2;
    
  supports = [
    [-bx/2 + 2 - tol,  by/2 - 0.5 + tol], // top left
    [ bx/2 - 2 + tol,  by/2 - 0.5 + tol], // top right
    [-bx/2 + 2 - tol, -by/2 + 0.5 - tol], // bottom left
    [ bx/2 - 2 + tol, -by/2 + 0.5 - tol], // bottom right
    [-bx/2 + 2 - tol,            middle], // middle left
    [              0,            middle], // middle center
    [ bx/2 - 2 + tol,            middle], // middle right
  ];

  translate([0, 0, -1.8])
    linear_extrude(1.8)
      for (s = supports)
        translate([s[0], s[1], -1.8])
          square([4, 1], center=true);
}

difference() {
  Base();
  difference() {
    Extrude();
    Supports();
  }
  ExtrudeSlideSwitch();
  CutoutDisplay();
  CutoutButtons();
  CutoutSlideSwitch();
  CutoutUSB();
  CutoutCharm();
}

// module mSnapHoles() {
//   oZ = 2.5;
//
//   holes = [
//     // clockwise
//     [[0, 180, 90], (pX-22)/2, pY],
//     [[0, 180, 90], (pX+22)/2, pY],
//     [[0, 180, 0],  pX,        pY*1/2],
//     [[0, 180, 0],  pX,        pY*1/10],
//     [[0, 0, 90],   (pX+22)/2, 0],
//     [[0, 0, 90],   (pX-22)/2, 0],
//     [[0, 0, 0],    0,         pY*1/10],
//     [[0, 0, 0],    0,         pY*1/2],
//   ];
//
//   for (h = holes) {
//     translate([h[1], h[2], -z+oZ])
//       rotate(h[0])
//         mSnapHole();
//   }
// }
//
// module mSnapHole() {
//   d = 0.5 + c;
//   l = 3 + 2 * c;
//
//   rotate([-90, 0, 0])
//     linear_extrude(l, center=true)
//       polygon([[-d, 0], [0, d], [0, -d]]);
// }
