# Relative Deltas (reldel)

This is an experimental delta function for doing diffs.

It uses unique flanking substrings to find positions of added/deleted
text. 

Ideally this will decrease the size of the diff and it also allows
same-line edits to merge with conflicts.
