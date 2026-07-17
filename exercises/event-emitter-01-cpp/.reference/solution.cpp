#include <algorithm>
#include <functional>
#include <map>
#include <string>
#include <vector>

// Every subscription gets a unique id in a single ordered list per
// event. Emit walks a snapshot of ids and re-checks liveness before
// each call -- that's what makes removal-during-emit safe. Once
// unhooks itself BEFORE calling, so re-entrant emits can't
// double-fire it.
class Emitter {
public:
    int On(const std::string& event, std::function<void(int)> fn) {
        return Add(event, fn, false);
    }

    int Once(const std::string& event, std::function<void(int)> fn) {
        return Add(event, fn, true);
    }

    void Off(int id) {
        auto it = by_id_.find(id);
        if (it == by_id_.end()) return;
        auto& list = subs_[it->second];
        list.erase(std::remove_if(list.begin(), list.end(),
                                  [id](const Sub& s) { return s.id == id; }),
                   list.end());
        by_id_.erase(it);
    }

    void Emit(const std::string& event, int v) {
        std::vector<int> ids;
        for (const auto& s : subs_[event]) ids.push_back(s.id);
        for (int id : ids) {
            auto owner = by_id_.find(id);
            if (owner == by_id_.end() || owner->second != event) continue;
            Sub* sub = nullptr;
            for (auto& s : subs_[event]) {
                if (s.id == id) { sub = &s; break; }
            }
            if (!sub) continue;
            auto fn = sub->fn;
            if (sub->once) Off(id);
            fn(v);
        }
    }

private:
    struct Sub {
        int id;
        std::function<void(int)> fn;
        bool once;
    };

    int Add(const std::string& event, std::function<void(int)> fn, bool once) {
        int id = ++next_id_;
        subs_[event].push_back({id, fn, once});
        by_id_[id] = event;
        return id;
    }

    std::map<std::string, std::vector<Sub>> subs_;
    std::map<int, std::string> by_id_;
    int next_id_ = 0;
};
