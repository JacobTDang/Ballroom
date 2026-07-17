#include <deque>
#include <functional>
#include <mutex>
#include <condition_variable>
#include <thread>
#include <vector>

// Server runs handle on submitted jobs with a pool of workers. Stop
// must drain everything accepted, refuse new work, and only return
// when the workers are done.
//
// TODO: this Stop flips the flag and returns immediately -- queued
// jobs are abandoned and the workers are killed mid-drain by detach.
class Server {
public:
    Server(int workers, std::function<void(int)> handle) : handle_(handle) {
        for (int w = 0; w < workers; w++) {
            std::thread([this] { Worker(); }).detach();
        }
    }

    bool Submit(int v) {
        if (stopped_) return false;
        std::lock_guard<std::mutex> g(mu_);
        jobs_.push_back(v);
        cv_.notify_one();
        return true;
    }

    void Stop() {
        stopped_ = true;
    }

    ~Server() { Stop(); }

private:
    void Worker() {
        while (true) {
            int v;
            {
                std::unique_lock<std::mutex> lock(mu_);
                cv_.wait(lock, [this] { return !jobs_.empty(); });
                v = jobs_.front();
                jobs_.pop_front();
            }
            handle_(v);
        }
    }

    std::function<void(int)> handle_;
    std::deque<int> jobs_;
    bool stopped_ = false;
    std::mutex mu_;
    std::condition_variable cv_;
};
