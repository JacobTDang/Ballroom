#include <cassert>
#include <cstdio>

bool SnapshotSurvivesEdit(int rows, int cols, int r, int c, int initial, int new_value);
bool SequentialSnapshotsAreIndependent(int rows, int cols, int r, int c, int initial,
                                        int after_first_edit, int after_second_edit);

int main() {
    // Headline case: edit after the snapshot must not change it.
    assert(SnapshotSurvivesEdit(2, 2, 0, 0, 1, 999));

    // The bug is whole-grid aliasing, not row-specific, but exercise a
    // non-origin cell on a bigger grid too.
    assert(SnapshotSurvivesEdit(3, 2, 2, 1, 5, 777));
    assert(SnapshotSurvivesEdit(3, 2, 0, 0, 42, -1));

    // Two snapshots taken at different times must stay independent of
    // each other, not just of the live grid.
    assert(SequentialSnapshotsAreIndependent(2, 2, 0, 0, 1, 111, 222));

    printf("all assertions passed\n");
    return 0;
}
