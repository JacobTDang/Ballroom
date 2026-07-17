#include <condition_variable>
#include <mutex>

// Barrier makes n threads rendezvous: each Wait blocks until all n
// have arrived, then all proceed -- and the barrier must be reusable
// for the next round.
//
// TODO: the released_ flag is never reset for the next round -- after
// round one, Wait falls straight through and nobody rendezvouses.
class Barrier {
public:
    explicit Barrier(int n) : n_(n) {}

    void Wait() {
        std::unique_lock<std::mutex> lock(mu_);
        if (released_) return;
        arrived_++;
        if (arrived_ == n_) {
            arrived_ = 0;
            released_ = true;
            cv_.notify_all();
            return;
        }
        cv_.wait(lock, [this] { return released_; });
    }

private:
    int n_;
    int arrived_ = 0;
    bool released_ = false;
    std::mutex mu_;
    std::condition_variable cv_;
};
