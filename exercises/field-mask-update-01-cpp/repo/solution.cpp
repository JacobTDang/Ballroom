#include <map>
#include <string>
#include <vector>

// Value is either an object (a map of named fields) or a scalar (a
// leaf string). Keep this as your value representation -- building
// yet another JSON-ish variant isn't the exercise.
struct Value {
    enum Kind { OBJECT, SCALAR } kind = SCALAR;
    std::map<std::string, Value> object;
    std::string scalar;

    bool IsObject() const { return kind == OBJECT; }

    static Value Obj(std::map<std::string, Value> o) {
        Value v;
        v.kind = OBJECT;
        v.object = std::move(o);
        return v;
    }
    static Value Str(std::string s) {
        Value v;
        v.kind = SCALAR;
        v.scalar = std::move(s);
        return v;
    }
};

// Update applies mask's dotted paths, copying values from source into
// target in place.
//
// TODO: ignores the mask entirely -- just shallow-merges source's
// top-level fields into target. No clearing, no path validation, no
// recursive per-path merge. Every rule in the problem statement is
// still yours to build.
bool Update(Value* target, const Value& source, const std::vector<std::string>& mask,
            std::string* err) {
    for (auto& [k, v] : source.object) {
        target->object[k] = v;
    }
    return true;
}
