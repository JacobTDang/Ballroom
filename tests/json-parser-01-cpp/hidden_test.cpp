#include "solution.cpp"

#include <cstdio>

int main() {
    {
        Json v;
        std::string err;
        std::string doc = "{\"name\": \"ada\", \"age\": -3, \"tags\": [\"a\", \"b\"], "
                          "\"meta\": {\"ok\": true, \"note\": null}, \"empty\": []}";
        if (!Parse(doc, &v, &err)) {
            fprintf(stderr, "Parse failed: %s\n", err.c_str());
            return 1;
        }
        if (v.kind != Json::OBJECT ||
            v.object["name"].str != "ada" ||
            v.object["age"].number != -3 ||
            v.object["tags"].array.size() != 2 ||
            v.object["tags"].array[1].str != "b" ||
            v.object["meta"].object["ok"].boolean != true ||
            v.object["meta"].object["note"].kind != Json::NUL ||
            !v.object["empty"].array.empty()) {
            fprintf(stderr, "nested structure parsed wrong\n");
            return 1;
        }
    }
    {
        Json v;
        std::string err;
        if (!Parse("\"say \\\"hi\\\" and \\\\\"", &v, &err) ||
            v.str != "say \"hi\" and \\") {
            fprintf(stderr, "escape handling wrong: %s / %s\n", err.c_str(), v.str.c_str());
            return 1;
        }
    }
    {
        Json v;
        std::string err;
        if (!Parse("  { \"a\" :  [ 1 , 2 ]  }  ", &v, &err) ||
            v.object["a"].array.size() != 2 || v.object["a"].array[0].number != 1) {
            fprintf(stderr, "whitespace handling wrong\n");
            return 1;
        }
    }
    {
        struct Case { const char* doc; const char* pos; } cases[] = {
            {"{\"a\" 1}", "5"},
            {"\"unterminated", "0"},
            {"tru", "0"},
            {"{\"a\": 1} extra", "9"},
        };
        for (const auto& c : cases) {
            Json v;
            std::string err;
            if (Parse(c.doc, &v, &err)) {
                fprintf(stderr, "Parse(%s) succeeded, want error\n", c.doc);
                return 1;
            }
            if (err.find(c.pos) == std::string::npos) {
                fprintf(stderr, "Parse(%s) error '%s' should name position %s\n", c.doc, err.c_str(), c.pos);
                return 1;
            }
        }
    }
    printf("all assertions passed\n");
    return 0;
}
