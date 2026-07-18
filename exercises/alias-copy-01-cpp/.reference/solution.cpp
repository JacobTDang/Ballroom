#include <vector>

class Grid {
public:
    Grid(int rows, int cols) : cells_(rows, std::vector<int>(cols, 0)) {}

    int Get(int r, int c) const { return cells_[r][c]; }
    void Set(int r, int c, int v) { cells_[r][c] = v; }

    // Returns an independent copy of the grid's current cell values,
    // for a caller that wants to edit the live grid and later compare
    // against or restore this saved state.
    Grid Snapshot() { return *this; }

private:
    std::vector<std::vector<int>> cells_;
};

// Exercises Snapshot()'s contract: builds a rows x cols grid, seeds
// cell (r, c), takes a snapshot, edits that cell, and reports whether
// the snapshot still shows the pre-edit value.
bool SnapshotSurvivesEdit(int rows, int cols, int r, int c, int initial, int new_value) {
    Grid g(rows, cols);
    g.Set(r, c, initial);
    const Grid& snap = g.Snapshot();
    g.Set(r, c, new_value);
    return snap.Get(r, c) == initial;
}

// Exercises two snapshots taken at different points in time: each
// must keep reflecting its own moment, independent of edits made
// after it was taken (including edits made after the OTHER snapshot).
bool SequentialSnapshotsAreIndependent(int rows, int cols, int r, int c, int initial,
                                        int after_first_edit, int after_second_edit) {
    Grid g(rows, cols);
    g.Set(r, c, initial);
    const Grid& snap1 = g.Snapshot();
    g.Set(r, c, after_first_edit);
    const Grid& snap2 = g.Snapshot();
    g.Set(r, c, after_second_edit);
    return snap1.Get(r, c) == initial && snap2.Get(r, c) == after_first_edit;
}
