use<_base.scad>

$fn = 25;

// switch dimensions
sx = 4.5;
sy = 3;
travel = 1.5;

// wall thickness
wlr = 2.0; // left-right

module Top() {
  linear_extrude(wlr)
    rsquare([sx, sy], r=sy/2, center=true);
}

module Bottom() {
  translate([0, 0, -wlr/2])
    linear_extrude(wlr/2+0.001) // prevent rendering issues
      difference() {
        square([sx+travel+1.5, sy+2], center=true);
        square([1.5, sy+2], center=true);
      }
}

Top();
Bottom();
