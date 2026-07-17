#include <functional>

// Lazy computes a value on first use -- the init function is expensive
// and must run exactly once, no matter how many threads call Get
// concurrently.
//
// TODO: the check below isn't atomic with the assignment -- two
// threads can both see done_ == false and both run init.
class Lazy {
public:
    explicit Lazy(std::function<int()> init) : init_(init) {}

    int Get() {
        if (!done_) {
            value_ = init_();
            done_ = true;
        }
        return value_;
    }

private:
    std::function<int()> init_;
    int value_ = 0;
    bool done_ = false;
};
