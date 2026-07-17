#include "solution.cpp"

#include <cstdio>

int main() {
    {
        std::map<std::string, std::map<std::string, std::string>> got;
        std::string err;
        std::string doc =
            "top = level\n# a comment\n; another\n\n[server]\nhost = localhost\n"
            "port = 8080\nhost = example.com\n\n[client]\nretries=3\n";
        if (!Parse(doc, &got, &err)) {
            fprintf(stderr, "Parse failed: %s\n", err.c_str());
            return 1;
        }
        if (got[""]["top"] != "level" || got["server"]["host"] != "example.com" ||
            got["server"]["port"] != "8080" || got["client"]["retries"] != "3") {
            fprintf(stderr, "parsed structure wrong\n");
            return 1;
        }
    }
    {
        std::map<std::string, std::map<std::string, std::string>> got;
        std::string err;
        if (!Parse("  spaced key   =   spaced value  ", &got, &err) ||
            got[""]["spaced key"] != "spaced value") {
            fprintf(stderr, "whitespace trimming wrong\n");
            return 1;
        }
    }
    {
        std::map<std::string, std::map<std::string, std::string>> got;
        std::string err;
        if (Parse("ok = 1\nnot a valid line\nok2 = 2", &got, &err)) {
            fprintf(stderr, "malformed line accepted\n");
            return 1;
        }
        if (err.find("2") == std::string::npos) {
            fprintf(stderr, "error should name line 2: %s\n", err.c_str());
            return 1;
        }
    }
    {
        std::map<std::string, std::map<std::string, std::string>> got;
        std::string err;
        if (Parse("[server]\nkey = v\n[broken", &got, &err) ||
            err.find("3") == std::string::npos) {
            fprintf(stderr, "unclosed header should error naming line 3: %s\n", err.c_str());
            return 1;
        }
    }
    printf("all assertions passed\n");
    return 0;
}
