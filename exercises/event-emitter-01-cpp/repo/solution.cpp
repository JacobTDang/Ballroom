#include <functional>
#include <map>
#include <string>
#include <vector>

// Emitter: On/Once subscribe (returning an id), Off unsubscribes,
// Emit calls the event's handlers in registration order.
//
// TODO: no ids (always 0), Off does nothing, and Once is just On --
// it never unhooks itself.
class Emitter {
public:
    int On(const std::string& event, std::function<void(int)> fn) {
        handlers_[event].push_back(fn);
        return 0;
    }

    int Once(const std::string& event, std::function<void(int)> fn) {
        return On(event, fn);
    }

    void Off(int id) {}

    void Emit(const std::string& event, int v) {
        for (auto& fn : handlers_[event]) fn(v);
    }

private:
    std::map<std::string, std::vector<std::function<void(int)>>> handlers_;
};
