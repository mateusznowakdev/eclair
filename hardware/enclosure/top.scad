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

// total height + 0.2mm for soldering inaccuracies:
// - buttons: 1.8mm
// -     PCB: 0.8mm
// -   USB-C: 3.2mm
z = 6;

module Base() {
  translate([0, 0, -z])
    linear_extrude(z+wfr)
      rsquare([bx+2*wlr+2*tol, by+2*wtb+2*tol], r=3, center=true);
}

module Extrude() {
  translate([0, 0, -z])
    linear_extrude(z)
      rsquare([bx+2*tol, by+2*tol], r=1.5-tol, center=true);
}

module ExtrudeSlideSwitch() {
  // see CutoutSlideSwitch() for math info
  translate([-bx/2-0.5-tol, (12+10.5)/2, -z])
    linear_extrude(z)
      square([1, 9+4*tol], center=true);
}

module CutoutDisplay() {
  translate([0, 12])
    linear_extrude(wfr+0.001) // prevent rendering issues
      rsquare([25, 8], r=1, center=true);
}

module CutoutButtons() {
  buttons = [
    [-2.0,  0.0],
    [-1.0,  0.0],
    [ 0.0,  0.0],
    [ 1.0,  0.0],
    [ 2.0,  0.0],
    [-2.0, -1.0],
    [-1.0, -1.0],
    [ 0.0, -1.0],
    [ 1.0, -1.0],
    [ 2.0, -1.0],
    [-1.5, -2.0],
    [-0.5, -2.0],
    [ 0.5, -2.0],
    [ 1.5, -2.0],
  ];

  translate([0, 1.6])
    linear_extrude(wfr+0.001)
      for (b = buttons)
        translate([b[0]*8.89, b[1]*7.62])
          circle(2);
}

module CutoutSlideSwitch() {
  // top and bottom positions are is 12mm and 10.5mm away from the center
  // the main body of the switch is 4.5mm x 3mm, with 1.5mm travel
  translate([-bx/2-tol, (12+10.5)/2, -z/2])
    rotate([-90, 0, 90])
      linear_extrude(wlr)
        rsquare([4.5+1.5+2*tol, 3+2*tol], r=1.5, center=true);
}

module CutoutUSB() {
  w = 8.89;
  h = 3.2;

  translate([-w/2, by/2+tol-0.001, -z+h])
    rotate([-90, 0, 0])
      linear_extrude(wtb+0.001)
        square([w, h]);
}

module CutoutCharm() {
  // pads are 9mm and 11mm away from the center, but this would be too fragile
  // make holes at -8.5mm and -11.5mm instead
  translate([bx/2+tol, -8.5, -z+2])
    rotate([0, 90, 0])
      linear_extrude(wlr) {
        circle(0.75);
        translate([0, -3, 0]) circle(0.75);
      }
}

module Supports() {
  middle = 1.6-7.62/2; // Y dimension between buttons, see CutoutButtons()
    
  supports = [
    [-bx/2+2-tol,  by/2-0.5+tol], // top left
    [ bx/2-2+tol,  by/2-0.5+tol], // top right
    [-bx/2+2-tol, -by/2+0.5-tol], // bottom left
    [ bx/2-2+tol, -by/2+0.5-tol], // bottom right
    [-bx/2+2-tol,        middle], // middle left
    [          0,        middle], // middle center
    [ bx/2-2+tol,        middle], // middle right
  ];

  translate([0, 0, -1.9])
    linear_extrude(1.9)
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
