#include <map>
#include <string>

enum class BeginState { kExecute, kInFlight, kReplay };

// BeginResult is what BeginAt hands back. response is only meaningful
// when state == kReplay.
struct BeginResult {
    BeginState state;
    std::string response;
};

// IdempotencyStore: tracks one request per key through its lifecycle.
//
// TODO: no fingerprint tracking, no in-flight/completed distinction,
// no deadline at all -- everything after the first BeginAt just
// replays whatever was last stored, even if nothing ever completed.
// Every rule in the problem statement is still yours to build.
class IdempotencyStore {
public:
    explicit IdempotencyStore(long ttl_ms) : ttl_ms_(ttl_ms) {}

    bool BeginAt(const std::string& key, const std::string& fingerprint, long now_ms,
                 BeginResult* out, std::string* err) {
        auto it = seen_.find(key);
        if (it == seen_.end()) {
            seen_[key] = "";
            out->state = BeginState::kExecute;
            out->response.clear();
            return true;
        }
        out->state = BeginState::kReplay;
        out->response = it->second;
        return true;
    }

    bool CompleteAt(const std::string& key, const std::string& response, long now_ms,
                     std::string* err) {
        seen_[key] = response;
        return true;
    }

private:
    long ttl_ms_;
    std::map<std::string, std::string> seen_;
};
