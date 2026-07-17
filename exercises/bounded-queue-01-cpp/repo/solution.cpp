#include <vector>

// BoundedQueue: a fixed-capacity FIFO shared between producer and
// consumer threads. Put must block while full; Get must block while
// empty.
//
// TODO: this version has no synchronization, no bound, and no
// blocking -- Get returns 0 when empty.
class BoundedQueue {
public:
    explicit BoundedQueue(int capacity) : capacity_(capacity) {}

    void Put(int v) {
        items_.push_back(v);
    }

    int Get() {
        if (items_.empty()) return 0;
        int v = items_.front();
        items_.erase(items_.begin());
        return v;
    }

private:
    int capacity_;
    std::vector<int> items_;
};
