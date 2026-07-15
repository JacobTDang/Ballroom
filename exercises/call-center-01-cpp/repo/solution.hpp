#pragma once

#include <string>

// Call center with respondents, managers, and directors. Calls escalate
// respondent -> manager -> director, then queue FIFO.
class CallCenter {
public:
    CallCenter(int respondents, int managers, int directors) {}

    // Route a new call. Return the handling level ("respondent",
    // "manager", "director") or "queued" when everyone is busy.
    std::string dispatch(int call_id) {
        // TODO: implement
        return "";
    }

    // Finish an active call (freeing its employee -- who takes the
    // longest-waiting queued call) or abandon a queued one. Return false
    // for unknown/already-ended calls.
    bool end_call(int call_id) {
        // TODO: implement
        return false;
    }

    // The level handling the call, "queued" if waiting, or "" if unknown
    // or ended.
    std::string handler_of(int call_id) {
        // TODO: implement
        return "";
    }
};
