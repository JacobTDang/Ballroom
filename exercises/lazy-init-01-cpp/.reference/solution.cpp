#include <functional>
#include <mutex>

// Lazy: std::call_once IS this problem, solved -- exactly-once
// execution with the result visible to every returning caller.
class Lazy {
public:
    explicit Lazy(std::function<int()> init) : init_(init) {}

    int Get() {
        std::call_once(flag_, [this] { value_ = init_(); });
        return value_;
    }

private:
    std::function<int()> init_;
    std::once_flag flag_;
    int value_ = 0;
};
