#include "solution.cpp"

#include <cstdio>

#define CHECK(cond, msg)                          \
    if (!(cond)) {                                \
        fprintf(stderr, "FAILED: %s\n", msg);    \
        return 1;                                 \
    }

int main() {
    {
        Emitter e;
        std::vector<std::string> got;
        e.On("a", [&](int) { got.push_back("first"); });
        e.On("a", [&](int) { got.push_back("second"); });
        e.On("b", [&](int) { got.push_back("other"); });
        e.Emit("a", 1);
        CHECK((got == std::vector<std::string>{"first", "second"}), "registration order / isolation");
    }
    {
        Emitter e;
        int calls = 0;
        e.Once("a", [&](int) { calls++; });
        e.Emit("a", 1);
        e.Emit("a", 2);
        CHECK(calls == 1, "once fired more than once");
    }
    {
        Emitter e;
        std::vector<std::string> got;
        int id = e.On("a", [&](int) { got.push_back("removed"); });
        e.On("a", [&](int) { got.push_back("kept"); });
        e.Off(id);
        e.Emit("a", 1);
        CHECK((got == std::vector<std::string>{"kept"}), "Off didn't remove the subscription");
    }
    {
        Emitter e;
        std::vector<std::string> got;
        int victim = 0;
        e.On("a", [&](int) {
            got.push_back("assassin");
            e.Off(victim);
        });
        victim = e.On("a", [&](int) { got.push_back("victim"); });
        e.On("a", [&](int) { got.push_back("bystander"); });
        e.Emit("a", 1);
        CHECK((got == std::vector<std::string>{"assassin", "bystander"}), "removed-during-emit handler fired");
    }
    {
        Emitter e;
        e.Emit("nobody", 42); // must not crash
        int got = 0;
        e.On("a", [&](int v) { got = v; });
        e.Emit("a", 99);
        CHECK(got == 99, "handler value");
    }
    printf("all assertions passed\n");
    return 0;
}
