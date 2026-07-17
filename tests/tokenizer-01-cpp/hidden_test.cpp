#include "solution.cpp"

#include <cstdio>

static bool TokenEq(const Token& t, const char* kind, const char* text, int pos) {
    return t.kind == kind && t.text == text && t.pos == pos;
}

int main() {
    {
        std::string err;
        auto toks = Tokenize("3+4.5*x", &err);
        if (!err.empty() || toks.size() != 5 ||
            !TokenEq(toks[0], "number", "3", 0) ||
            !TokenEq(toks[1], "op", "+", 1) ||
            !TokenEq(toks[2], "number", "4.5", 2) ||
            !TokenEq(toks[3], "op", "*", 5) ||
            !TokenEq(toks[4], "ident", "x", 6)) {
            fprintf(stderr, "dense expression tokens wrong (err=%s, %zu tokens)\n", err.c_str(), toks.size());
            return 1;
        }
    }
    {
        std::string err;
        auto toks = Tokenize("price * (1 + tax_rate2)", &err);
        if (!err.empty() || toks.size() != 7 ||
            !TokenEq(toks[0], "ident", "price", 0) ||
            !TokenEq(toks[2], "lparen", "(", 8) ||
            !TokenEq(toks[5], "ident", "tax_rate2", 13) ||
            !TokenEq(toks[6], "rparen", ")", 22)) {
            fprintf(stderr, "parens/ident tokens wrong\n");
            return 1;
        }
    }
    {
        std::string err;
        Tokenize("12..3", &err);
        if (err.find("3") == std::string::npos) {
            fprintf(stderr, "12..3 should error naming position 3, got %s\n", err.c_str());
            return 1;
        }
    }
    {
        std::string err;
        Tokenize("a @ b", &err);
        if (err.find("2") == std::string::npos) {
            fprintf(stderr, "@ should error naming position 2, got %s\n", err.c_str());
            return 1;
        }
    }
    {
        std::string err;
        auto toks = Tokenize("", &err);
        if (!err.empty() || !toks.empty()) {
            fprintf(stderr, "empty input should be empty tokens, no error\n");
            return 1;
        }
    }
    printf("all assertions passed\n");
    return 0;
}
