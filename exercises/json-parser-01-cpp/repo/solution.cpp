#include <map>
#include <string>
#include <vector>

// Json is the parsed value -- keep this struct as your output type.
struct Json {
    enum Kind { OBJECT, ARRAY, STRING, NUMBER, BOOL, NUL } kind = NUL;
    std::map<std::string, Json> object;
    std::vector<Json> array;
    std::string str;
    long number = 0;
    bool boolean = false;
};

// Parse a JSON subset: objects, arrays, strings (\" and \\ escapes),
// integers, true/false/null. On error, fill *err naming the byte
// position and return false.
//
// TODO: this recognizes nothing but the empty object.
bool Parse(const std::string& input, Json* out, std::string* err) {
    if (input == "{}") {
        out->kind = Json::OBJECT;
        return true;
    }
    *err = "only {} supported so far";
    return false;
}
