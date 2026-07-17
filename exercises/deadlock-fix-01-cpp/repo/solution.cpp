#include <chrono>
#include <mutex>
#include <thread>

// Account is a bank account with its own lock.
class Account {
public:
    Account(int id, int balance) : id(id), balance_(balance) {}

    int Balance() {
        std::lock_guard<std::mutex> g(mu);
        return balance_;
    }

    int id;
    std::mutex mu;
    int balance_;
};

// Transfer moves amount between accounts, locking both.
//
// TODO: locking from-then-to deadlocks the moment two transfers cross
// (A->B and B->A each hold one lock and wait for the other). Fix the
// ordering -- don't just wrap everything in one global lock.
bool Transfer(Account& from, Account& to, int amount) {
    from.mu.lock();
    std::this_thread::sleep_for(std::chrono::milliseconds(1)); // bookkeeping
    to.mu.lock();

    bool ok = false;
    if (from.balance_ >= amount) {
        from.balance_ -= amount;
        to.balance_ += amount;
        ok = true;
    }
    to.mu.unlock();
    from.mu.unlock();
    return ok;
}
