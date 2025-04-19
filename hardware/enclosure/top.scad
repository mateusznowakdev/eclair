use<_base.scad>

$fn = 25;

// board dimensions
pX = 44.45;
pY = 34.93;
pZ = 0.8;

// wall thickness
tX = 2.0;
tY = 1.5;
tZ = 0.75;

// tolerance/clearance
c = 0.15;

// total height
z = 5.25 + pZ + tZ;

module mOuter() {
  oX = tX + c;
  oY = tY + c;
  r = 3.5;

  x = pX + 2 * oX;
  y = pY + 2 * oY;

  translate([-oX, -oY, -z])
    linear_extrude(z)
      rsquare([x, y], r);
}

module mInner() {
  oX = c;
  oY = c;
  r = 2;

  x = pX + 2 * oX;
  y = pY + 2 * oY;

  sx = tX / 2;
  sy = 14;

  translate([-oX, -oY, -z])
    linear_extrude(z-tZ)
      union() {
        // main
        rsquare([x, y], r);
        // slide switch
        translate([-sx, y-sy])
          rsquare([r*3, sy], r);
      }
}

module mButtonGrid() {
  dX = 8.89;
  dY = 7.62;

  oY = 1.6;

  buttons = [
    [[-2*dX,    0],    "QW"],
    [[-1*dX,    0],    "ER"],
    [[ 0,       0],    "TY"],
    [[ 1*dX,    0],    "UI"],
    [[ 2*dX,    0],    "OP"],
    [[-2*dX,   -1*dY], "AS"],
    [[-1*dX,   -1*dY], "DF"],
    [[ 0,      -1*dY], "GH"],
    [[ 1*dX,   -1*dY], "JK"],
    [[ 2*dX,   -1*dY], "L-"],
    [[-1.5*dX, -2*dY], "ZX"],
    [[-0.5*dX, -2*dY], "CV"],
    [[ 0.5*dX, -2*dY], "BN"],
    [[ 1.5*dX, -2*dY], "M."],
  ];

  translate([pX/2, pY/2+oY])
    for(b = buttons)
      mButton(b[0], b[1]);
}

module mButton(b, t) {
  x = b[0];
  y = b[1];
  r = 2.0;
  s = 2.5;

  translate([x, y, -tZ])
    linear_extrude(tZ)
      circle(r);

  // translate([x+0.5, y+s, -tZ/2])
  //   linear_extrude(tZ/2)
  //     offset(delta=0.1)
  //       text(t, font="Bungee Hairline", size=s, halign="center", spacing=1.5);
}

module mDisplay() {
  x = 24.5;
  y = 7.5;
  r = 2;

  oY = 8.5;

  translate([(pX-x)/2, pY/2+oY, -tZ])
    linear_extrude(tZ)
      rsquare([x, y], r);
}

module mUSB() {
  x = 8.89;
  y = 3.25;

  oZ = tZ + 1.75 + pZ;

  translate([(pX-x)/2, pY, -oZ])
    rotate([-90, 0, 0])
      linear_extrude(tY+c)
        union() {
          // main
          rsquare([x, y], y/2);
          // remove material under the receptacle
          translate([0, y/2]) square([x, y]);
        }
}

module mSlide() {
  x = 6 + 2 * c;
  y = 3 + 2 * c;
  r = y / 2 - c;

  oY = 3;
  oZ = tZ + 0.75;

  translate([0, pY-oY, -y-oZ])
    rotate([90, 0, -90])
      linear_extrude(tX+c)
        rsquare([x, y], r);
}

module mSnapHoles() {
  oZ = 2.5;

  holes = [
    // clockwise
    [[0, 180, 90], (pX-22)/2, pY],
    [[0, 180, 90], (pX+22)/2, pY],
    [[0, 180, 0],  pX,        pY*1/2],
    [[0, 180, 0],  pX,        pY*1/10],
    [[0, 0, 90],   (pX+22)/2, 0],
    [[0, 0, 90],   (pX-22)/2, 0],
    [[0, 0, 0],    0,         pY*1/10],
    [[0, 0, 0],    0,         pY*1/2],
  ];
  
  for (h = holes) {
    translate([h[1], h[2], -z+oZ])
      rotate(h[0])
        mSnapHole();
  }
}

module mSnapHole() {
  d = 0.5 + c;
  l = 3 + 2 * c;
  
  rotate([-90, 0, 0])
    linear_extrude(l, center=true)
      polygon([[-d, 0], [0, d], [0, -d]]);
}

module mCharmHole() {
  r = 1.5 / 2;
  
  translate([pX, 7.35, -z+2.5])
    rotate([90, 0, 90])
      linear_extrude(tX+c)
        union() {
          translate([-2*r, 0]) circle(r);
          translate([2*r, 0])  circle(r);
        }
}

difference() {
  mOuter();
  union() {
    mInner();
    mButtonGrid();
    mDisplay();
    mUSB();
    mSlide();
    mSnapHoles();
    mCharmHole();
  }
}
