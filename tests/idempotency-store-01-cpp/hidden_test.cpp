#include "solution.cpp"

#include <cstdio>

#define CHECK(cond, msg)                       \
    if (!(cond)) {                              \
        fprintf(stderr, "FAILED: %s\n", msg);   \
        return 1;                               \
    }

static const long kTTL = 1000;

int main() {
    BeginResult res;
    std::string err;

    // Lifecycle matrix: execute -> in-flight -> replay.
    {
        IdempotencyStore store(kTTL);

        CHECK(store.BeginAt("k1", "fp-a", 0, &res, &err), "first BeginAt failed");
        CHECK(res.state == BeginState::kExecute, "first BeginAt wasn't kExecute");

        CHECK(store.BeginAt("k1", "fp-a", 10, &res, &err), "duplicate BeginAt failed");
        CHECK(res.state == BeginState::kInFlight, "duplicate BeginAt wasn't kInFlight");

        CHECK(store.CompleteAt("k1", "RESULT-1", 20, &err), "CompleteAt failed");

        CHECK(store.BeginAt("k1", "fp-a", 30, &res, &err), "replay BeginAt failed");
        CHECK(res.state == BeginState::kReplay && res.response == "RESULT-1",
              "replay BeginAt didn't return the stored response");
    }

    // Byte-identical replay.
    {
        IdempotencyStore store(kTTL);
        store.BeginAt("k1", "fp-a", 0, &res, &err);
        std::string payload = "{\"amount\": 4200, \"currency\": \"usd\", \"note\": \"caf\xc3\xa9\"}";
        CHECK(store.CompleteAt("k1", payload, 10, &err), "CompleteAt failed");

        CHECK(store.BeginAt("k1", "fp-a", 20, &res, &err), "BeginAt failed");
        CHECK(res.response == payload, "replay must return the stored response byte-for-byte");
    }

    // Conflict on a live key (in-flight, and completed-but-not-expired)
    // with a different fingerprint.
    {
        IdempotencyStore store(kTTL);
        store.BeginAt("k1", "fp-a", 0, &res, &err);
        CHECK(!store.BeginAt("k1", "fp-b", 5, &res, &err),
              "in-flight key with a mismatched fingerprint didn't error");

        IdempotencyStore store2(kTTL);
        store2.BeginAt("k2", "fp-a", 0, &res, &err);
        store2.CompleteAt("k2", "RESULT", 5, &err);
        CHECK(!store2.BeginAt("k2", "fp-b", 10, &res, &err),
              "completed key with a mismatched fingerprint didn't error");
    }

    // Exact TTL boundary.
    {
        IdempotencyStore store(kTTL);
        store.BeginAt("k1", "fp-a", 0, &res, &err);
        store.CompleteAt("k1", "RESULT", 0, &err); // deadline = 0 + ttl

        CHECK(store.BeginAt("k1", "fp-a", kTTL - 1, &res, &err), "BeginAt failed");
        CHECK(res.state == BeginState::kReplay && res.response == "RESULT",
              "just inside the retention window: expected kReplay/RESULT");

        CHECK(store.BeginAt("k1", "fp-a", kTTL, &res, &err), "BeginAt failed");
        CHECK(res.state == BeginState::kExecute,
              "deadline reached: expected kExecute (brand new key)");
    }

    // CompleteAt on unknown/expired -> error.
    {
        IdempotencyStore store(kTTL);
        CHECK(!store.CompleteAt("never-begun", "R", 5, &err),
              "CompleteAt on an unknown key didn't error");

        IdempotencyStore store2(kTTL);
        store2.BeginAt("k1", "fp-a", 0, &res, &err); // in-flight, deadline = ttl
        CHECK(!store2.CompleteAt("k1", "R", kTTL, &err),
              "CompleteAt on a key whose deadline already passed didn't error");
    }

    printf("all assertions passed\n");
    return 0;
}
