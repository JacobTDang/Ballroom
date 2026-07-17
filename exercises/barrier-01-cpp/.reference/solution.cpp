#include <condition_variable>
#include <mutex>

// Barrier: a generation counter makes reuse safe -- each round waits
// on its own generation, so a round-2 arrival can never consume a
// round-1 release. The last arriver flips the generation and notifies
// all; waiters wait until THEIR generation has passed.
class Barrier {
public:
    explicit Barrier(int n) : n_(n) {}

    void Wait() {
        std::unique_lock<std::mutex> lock(mu_);
        long gen = generation_;
        arrived_++;
        if (arrived_ == n_) {
            arrived_ = 0;
            generation_++;
            cv_.notify_all();
            return;
        }
        cv_.wait(lock, [this, gen] { return gen != generation_; });
    }

private:
    int n_;
    int arrived_ = 0;
    long generation_ = 0;
    std::mutex mu_;
    std::condition_variable cv_;
};
