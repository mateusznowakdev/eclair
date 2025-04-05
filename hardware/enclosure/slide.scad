use<_base.scad>

$fn = 25;

sx = 4.5;
sy = 3;

bx = sx + 4;
by = sy + 2.5;

h1 = 0.75;
h2 = 2.0;

module mBase() {
  translate([-2, -0.5])
    linear_extrude(h1)
      square([bx, by]);
}

module mStick() {
  translate([0, 0, h1])
    linear_extrude(h2)
      rsquare([sx, sy], sy/2);
}

module mCutTop() {
  w = sx/3;

  translate([(sx-w)/2, -10, (h2+h1)-0.25])
    linear_extrude(h1)
      square([w, 20]);
}

module mCutBottom() {
  w = 1.6;

  translate([(sx-w)/2, -10])
    linear_extrude(h1)
      square([w, 20]);
}

difference() {
  union() {
    mBase();
    mStick();
  }
  //mCutTop();
  mCutBottom();
}
