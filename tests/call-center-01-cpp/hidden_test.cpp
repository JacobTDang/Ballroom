#include <cassert>
#include <cstdio>

#include "solution.hpp"

int main() {
    {
        CallCenter c(1, 1, 1);
        assert(c.dispatch(1) == "respondent");        // respondent first
    }
    {
        CallCenter c(1, 1, 1);
        c.dispatch(1);
        assert(c.dispatch(2) == "manager");           // escalate to manager
    }
    {
        CallCenter c(1, 1, 1);
        c.dispatch(1);
        c.dispatch(2);
        assert(c.dispatch(3) == "director");          // escalate to director
    }
    {
        CallCenter c(1, 0, 0);
        c.dispatch(1);
        assert(c.dispatch(2) == "queued");            // queue when all busy
        assert(c.handler_of(2) == "queued");
    }
    {
        CallCenter c(1, 1, 0);                        // freed employee takes queued call
        c.dispatch(1);
        c.dispatch(2);
        c.dispatch(3);
        assert(c.end_call(1) == true);
        assert(c.handler_of(3) == "respondent");
    }
    {
        CallCenter c(2, 0, 0);                        // FIFO assignment
        c.dispatch(1);
        c.dispatch(2);
        c.dispatch(3);
        c.dispatch(4);
        c.end_call(2);
        assert(c.handler_of(3) == "respondent");
        assert(c.handler_of(4) == "queued");
    }
    {
        CallCenter c(1, 0, 0);                        // freed slot reusable
        c.dispatch(1);
        c.end_call(1);
        assert(c.dispatch(2) == "respondent");
    }
    {
        CallCenter c(1, 1, 1);
        assert(c.end_call(99) == false);              // unknown call
    }
    {
        CallCenter c(1, 0, 0);                        // double end
        c.dispatch(1);
        assert(c.end_call(1) == true);
        assert(c.end_call(1) == false);
    }
    {
        CallCenter c(1, 0, 0);                        // abandon queued call
        c.dispatch(1);
        c.dispatch(2);
        c.dispatch(3);
        assert(c.end_call(2) == true);
        c.end_call(1);
        assert(c.handler_of(3) == "respondent");
        assert(c.handler_of(2) == "");
    }
    {
        CallCenter c(1, 1, 1);
        assert(c.handler_of(42) == "");               // unknown handler
    }
    {
        CallCenter c(1, 1, 1);                        // active levels reported
        c.dispatch(1);
        c.dispatch(2);
        assert(c.handler_of(1) == "respondent");
        assert(c.handler_of(2) == "manager");
    }
    printf("all assertions passed\n");
    return 0;
}
