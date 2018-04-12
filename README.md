# Relative Deltas (reldel)

This is an experimental delta function for doing diffs.

It uses unique flanking substrings to find positions of added/deleted
text. 

Ideally this will decrease the size of the diff and it also allows
same-line edits to be merged without conflicts.

It is O(N^2) because of the alignment of texts using Needleman–Wunsch algorithm. (Beware of large texts).

# License

MIT