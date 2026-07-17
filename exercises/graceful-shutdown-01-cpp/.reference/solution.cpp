#include <deque>
#include <functional>
#include <mutex>
#include <condition_variable>
#include <thread>
#include <vector>

// Server: `stopping_` under the same mutex as the queue is the "no
// more work is coming" signal. Workers drain until the queue is empty
// AND stopping is set; Stop joins every worker, so when it returns the
// drain is complete. Submit checks stopping under the lock, so no job
// sneaks in during shutdown.
class Server {
public:
    Server(int workers, std::function<void(int)> handle) : handle_(handle) {
        for (int w = 0; w < workers; w++) {
            threads_.emplace_back([this] { Worker(); });
        }
    }

    bool Submit(int v) {
        std::lock_guard<std::mutex> g(mu_);
        if (stopping_) return false;
        jobs_.push_back(v);
        cv_.notify_one();
        return true;
    }

    void Stop() {
        {
            std::lock_guard<std::mutex> g(mu_);
            if (stopping_) return;
            stopping_ = true;
            cv_.notify_all();
        }
        for (auto& t : threads_) {
            if (t.joinable()) t.join();
        }
    }

    ~Server() { Stop(); }

private:
    void Worker() {
        while (true) {
            int v;
            {
                std::unique_lock<std::mutex> lock(mu_);
                cv_.wait(lock, [this] { return !jobs_.empty() || stopping_; });
                if (jobs_.empty()) return; // stopping and drained
                v = jobs_.front();
                jobs_.pop_front();
            }
            handle_(v);
        }
    }

    std::function<void(int)> handle_;
    std::deque<int> jobs_;
    bool stopping_ = false;
    std::mutex mu_;
    std::condition_variable cv_;
    std::vector<std::thread> threads_;
};
