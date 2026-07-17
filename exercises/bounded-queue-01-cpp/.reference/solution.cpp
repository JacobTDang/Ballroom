#include <condition_variable>
#include <deque>
#include <mutex>

// BoundedQueue: one mutex, two condition variables (not-full for
// producers, not-empty for consumers). The while-loops around wait()
// matter: a woken thread must re-check, because another thread may
// have taken the slot first.
class BoundedQueue {
public:
    explicit BoundedQueue(int capacity) : capacity_(capacity) {}

    void Put(int v) {
        std::unique_lock<std::mutex> lock(mu_);
        not_full_.wait(lock, [this] { return (int)items_.size() < capacity_; });
        items_.push_back(v);
        not_empty_.notify_one();
    }

    int Get() {
        std::unique_lock<std::mutex> lock(mu_);
        not_empty_.wait(lock, [this] { return !items_.empty(); });
        int v = items_.front();
        items_.pop_front();
        not_full_.notify_one();
        return v;
    }

private:
    int capacity_;
    std::deque<int> items_;
    std::mutex mu_;
    std::condition_variable not_full_;
    std::condition_variable not_empty_;
};
