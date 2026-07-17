#include <cctype>
#include <map>
#include <string>
#include <vector>

// Json is the parsed value.
struct Json {
    enum Kind { OBJECT, ARRAY, STRING, NUMBER, BOOL, NUL } kind = NUL;
    std::map<std::string, Json> object;
    std::vector<Json> array;
    std::string str;
    long number = 0;
    bool boolean = false;
};

// Recursive descent: one function per grammar rule, each consuming
// exactly its production and leaving i just past it. Every failure
// names the byte position.
namespace {

struct Parser {
    const std::string& s;
    size_t i = 0;
    std::string* err;

    void SkipSpace() {
        while (i < s.size() && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r')) i++;
    }

    bool Fail(const std::string& msg, size_t pos) {
        *err = msg + " at position " + std::to_string(pos);
        return false;
    }

    bool ParseValue(Json* out) {
        if (i >= s.size()) return Fail("unexpected end of input", i);
        char c = s[i];
        if (c == '{') return ParseObject(out);
        if (c == '[') return ParseArray(out);
        if (c == '"') return ParseString(out);
        if (c == '-' || isdigit((unsigned char)c)) return ParseNumber(out);
        return ParseLiteral(out);
    }

    bool ParseObject(Json* out) {
        out->kind = Json::OBJECT;
        i++;
        SkipSpace();
        if (i < s.size() && s[i] == '}') { i++; return true; }
        while (true) {
            SkipSpace();
            if (i >= s.size() || s[i] != '"') return Fail("expected object key", i);
            Json key;
            if (!ParseString(&key)) return false;
            SkipSpace();
            if (i >= s.size() || s[i] != ':') return Fail("expected ':'", i);
            i++;
            SkipSpace();
            Json value;
            if (!ParseValue(&value)) return false;
            out->object[key.str] = value;
            SkipSpace();
            if (i >= s.size()) return Fail("unterminated object", i);
            if (s[i] == ',') { i++; continue; }
            if (s[i] == '}') { i++; return true; }
            return Fail("expected ',' or '}'", i);
        }
    }

    bool ParseArray(Json* out) {
        out->kind = Json::ARRAY;
        i++;
        SkipSpace();
        if (i < s.size() && s[i] == ']') { i++; return true; }
        while (true) {
            SkipSpace();
            Json value;
            if (!ParseValue(&value)) return false;
            out->array.push_back(value);
            SkipSpace();
            if (i >= s.size()) return Fail("unterminated array", i);
            if (s[i] == ',') { i++; continue; }
            if (s[i] == ']') { i++; return true; }
            return Fail("expected ',' or ']'", i);
        }
    }

    bool ParseString(Json* out) {
        out->kind = Json::STRING;
        size_t start = i;
        i++;
        std::string result;
        while (i < s.size()) {
            char c = s[i];
            if (c == '"') { i++; out->str = result; return true; }
            if (c == '\\') {
                if (i + 1 >= s.size()) break;
                char nxt = s[i + 1];
                if (nxt == '"') result += '"';
                else if (nxt == '\\') result += '\\';
                else return Fail("unsupported escape", i);
                i += 2;
                continue;
            }
            result += c;
            i++;
        }
        return Fail("unterminated string starting", start);
    }

    bool ParseNumber(Json* out) {
        out->kind = Json::NUMBER;
        size_t start = i;
        bool neg = false;
        if (s[i] == '-') { neg = true; i++; }
        long n = 0;
        int digits = 0;
        while (i < s.size() && isdigit((unsigned char)s[i])) {
            n = n * 10 + (s[i] - '0');
            i++;
            digits++;
        }
        if (digits == 0) return Fail("malformed number", start);
        out->number = neg ? -n : n;
        return true;
    }

    bool ParseLiteral(Json* out) {
        if (s.compare(i, 4, "true") == 0) { out->kind = Json::BOOL; out->boolean = true; i += 4; return true; }
        if (s.compare(i, 5, "false") == 0) { out->kind = Json::BOOL; out->boolean = false; i += 5; return true; }
        if (s.compare(i, 4, "null") == 0) { out->kind = Json::NUL; i += 4; return true; }
        return Fail("unexpected value", i);
    }
};

} // namespace

bool Parse(const std::string& input, Json* out, std::string* err) {
    Parser p{input, 0, err};
    p.SkipSpace();
    if (!p.ParseValue(out)) return false;
    p.SkipSpace();
    if (p.i != input.size()) {
        *err = "trailing garbage at position " + std::to_string(p.i);
        return false;
    }
    return true;
}
