# EclairM0

[Project description](https://mateusznowak.dev/eclair) and [build instructions](https://mateusznowak.dev/eclair/build) can be found on my website.

KiCad files (board) can be previewed in the web browser: https://kicanvas.org<br/>OpenSCAD files (enclosure) can be previewed in the web browser: https://openscad.cloud/openscad/

License terms for software and hardware are present in the [LICENSE.md](./LICENSE.md) file.

## Updating font data

Font data is stored in the `tools/font.pbm` image file, with `tools/font.txt` containing font widths. The image file can be edited in GIMP.

These source files can be converted into working Go code, using the `tools/convert.py` script, which should work with any modern version of Python:

```bash
cd tools
python convert.py > ../firmware/hal/display/font.go
```
