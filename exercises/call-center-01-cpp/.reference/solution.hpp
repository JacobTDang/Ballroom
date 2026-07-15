#pragma once

#include <algorithm>
#include <deque>
#include <map>
#include <string>
#include <vector>

// Call center with respondents, managers, and directors. Calls escalate
// respondent -> manager -> director, then queue FIFO.
class CallCenter {
public:
    CallCenter(int respondents, int managers, int directors) {
        free_["respondent"] = respondents;
        free_["manager"] = managers;
        free_["director"] = directors;
    }

    // Route a new call. Return the handling level ("respondent",
    // "manager", "director") or "queued" when everyone is busy.
    std::string dispatch(int call_id) {
        for (const std::string level : {"respondent", "manager", "director"}) {
            if (free_[level] > 0) {
                free_[level]--;
                active_[call_id] = level;
                return level;
            }
        }
        queue_.push_back(call_id);
        return "queued";
    }

    // Finish an active call (freeing its employee -- who takes the
    // longest-waiting queued call) or abandon a queued one. Return false
    // for unknown/already-ended calls.
    bool end_call(int call_id) {
        auto it = active_.find(call_id);
        if (it != active_.end()) {
            std::string level = it->second;
            active_.erase(it);
            if (!queue_.empty()) {
                int next = queue_.front();
                queue_.pop_front();
                active_[next] = level;
            } else {
                free_[level]++;
            }
            return true;
        }
        auto qit = std::find(queue_.begin(), queue_.end(), call_id);
        if (qit != queue_.end()) {
            queue_.erase(qit);
            return true;
        }
        return false;
    }

    // The level handling the call, "queued" if waiting, or "" if unknown
    // or ended.
    std::string handler_of(int call_id) {
        auto it = active_.find(call_id);
        if (it != active_.end()) return it->second;
        if (std::find(queue_.begin(), queue_.end(), call_id) != queue_.end()) {
            return "queued";
        }
        return "";
    }

private:
    std::map<std::string, int> free_;
    std::map<int, std::string> active_;
    std::deque<int> queue_;
};
