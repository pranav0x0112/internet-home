#!/usr/bin/env sh
# Post-build: copy collections/<name>.html -> <name>/index.html
set -e
RENDERED_DIR="site/rendered"
COL_DIR="$RENDERED_DIR/collections"

if [ ! -d "$COL_DIR" ]; then
  echo "post_build_copy: no collections directory found at $COL_DIR"
  exit 0
fi

for src in "$COL_DIR"/*.html; do
  [ -e "$src" ] || continue
  filename=$(basename -- "$src")
  name=${filename%.html}
  dest_dir="$RENDERED_DIR/$name"
  dest_file="$dest_dir/index.html"
  mkdir -p "$dest_dir"
  cp -f "$src" "$dest_file"
  echo "post_build_copy: copied $src -> $dest_file"
done

exit 0
