#include <mutex>

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

// Transfer: std::scoped_lock acquires both mutexes with a deadlock-
// avoidance algorithm -- the standard library's answer to exactly this
// inversion. (Ordering by id by hand is equally valid.)
bool Transfer(Account& from, Account& to, int amount) {
    std::scoped_lock lock(from.mu, to.mu);
    if (from.balance_ < amount) return false;
    from.balance_ -= amount;
    to.balance_ += amount;
    return true;
}
