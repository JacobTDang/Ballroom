#include <map>
#include <string>

enum class BeginState { kExecute, kInFlight, kReplay };

// BeginResult is what BeginAt hands back. response is only meaningful
// when state == kReplay.
struct BeginResult {
    BeginState state;
    std::string response;
};

namespace {
struct Record {
    std::string fingerprint;
    bool in_flight;
    long deadline; // now_ms past this: treat the record as gone
    std::string response;
};
}  // namespace

// IdempotencyStore: one record per key -- fingerprint, in-flight or
// completed, a deadline that BeginAt sets and CompleteAt renews, and
// (once completed) the stored response. Past its deadline, a record
// is treated as if it never existed: BeginAt starts clean, and
// CompleteAt has nothing to attach to.
class IdempotencyStore {
public:
    explicit IdempotencyStore(long ttl_ms) : ttl_ms_(ttl_ms) {}

    bool BeginAt(const std::string& key, const std::string& fingerprint, long now_ms,
                 BeginResult* out, std::string* err) {
        auto it = records_.find(key);
        if (it == records_.end() || now_ms >= it->second.deadline) {
            records_[key] = Record{fingerprint, true, now_ms + ttl_ms_, ""};
            out->state = BeginState::kExecute;
            out->response.clear();
            return true;
        }

        Record& r = it->second;
        if (r.fingerprint != fingerprint) {
            *err = "idempotency: fingerprint conflict for key " + key;
            return false;
        }

        if (r.in_flight) {
            out->state = BeginState::kInFlight;
            out->response.clear();
        } else {
            out->state = BeginState::kReplay;
            out->response = r.response;
        }
        return true;
    }

    bool CompleteAt(const std::string& key, const std::string& response, long now_ms,
                     std::string* err) {
        auto it = records_.find(key);
        if (it == records_.end() || now_ms >= it->second.deadline || !it->second.in_flight) {
            *err = "idempotency: no in-flight request for key " + key;
            return false;
        }
        it->second.in_flight = false;
        it->second.response = response;
        it->second.deadline = now_ms + ttl_ms_;
        return true;
    }

private:
    long ttl_ms_;
    std::map<std::string, Record> records_;
};
