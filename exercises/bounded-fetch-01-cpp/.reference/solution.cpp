#include <condition_variable>
#include <functional>
#include <mutex>
#include <thread>
#include <vector>

// A counting semaphore from a mutex + condition variable (C++17 has no
// std::counting_semaphore). Acquired inside each thread immediately
// around the task body, so the bound covers execution, not launch.
class Semaphore {
public:
    explicit Semaphore(int count) : count_(count) {}

    void Acquire() {
        std::unique_lock<std::mutex> lock(mu_);
        cv_.wait(lock, [this] { return count_ > 0; });
        count_--;
    }

    void Release() {
        std::lock_guard<std::mutex> g(mu_);
        count_++;
        cv_.notify_one();
    }

private:
    int count_;
    std::mutex mu_;
    std::condition_variable cv_;
};

void RunLimited(const std::vector<std::function<void()>>& tasks, int limit) {
    Semaphore sem(limit);
    std::vector<std::thread> threads;
    for (const auto& task : tasks) {
        threads.emplace_back([&sem, task] {
            sem.Acquire();
            task();
            sem.Release();
        });
    }
    for (auto& t : threads) t.join();
}
