module rsquare(dim, r=0) {
  x = dim[0];
  y = dim[1];

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
