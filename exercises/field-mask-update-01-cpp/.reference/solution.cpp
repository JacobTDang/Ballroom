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

namespace {

std::vector<std::string> SplitDots(const std::string& path) {
    std::vector<std::string> segs;
    std::string cur;
    for (char c : path) {
        if (c == '.') {
            segs.push_back(cur);
            cur.clear();
        } else {
            cur.push_back(c);
        }
    }
    segs.push_back(cur);
    return segs;
}

// Walk segments through source. Returns false if the walk falls off
// the source -- a legitimate "clear", not an error.
bool SourceLookup(const Value& source, const std::vector<std::string>& segments, Value* out) {
    const Value* node = &source;
    for (auto& seg : segments) {
        if (!node->IsObject()) return false;
        auto it = node->object.find(seg);
        if (it == node->object.end()) return false;
        node = &it->second;
    }
    *out = *node;
    return true;
}

// Walk every segment except the last through target -- each one must
// already exist as an object. On success, *parent points at the
// object the final segment lives in.
bool TargetParent(Value* target, const std::vector<std::string>& segments,
                   const std::string& full_path, Value** parent, std::string* err) {
    Value* node = target;
    for (size_t i = 0; i + 1 < segments.size(); i++) {
        const std::string& seg = segments[i];
        if (!node->IsObject()) {
            *err = "unknown path \"" + full_path + "\": \"" + seg + "\" is not an object";
            return false;
        }
        auto it = node->object.find(seg);
        if (it == node->object.end()) {
            *err = "unknown path \"" + full_path + "\": \"" + seg + "\" does not exist";
            return false;
        }
        node = &it->second;
    }
    if (!node->IsObject()) {
        *err = "unknown path \"" + full_path + "\": not an object";
        return false;
    }
    *parent = node;
    return true;
}

struct MaskOp {
    Value* parent;
    std::string leaf;
    bool found;
    Value value;
};

}  // namespace

// Update: two passes on purpose. Validate every path's target-side
// intermediates first, THEN apply. A bad path anywhere in the mask
// must leave target completely untouched, not partially patched.
bool Update(Value* target, const Value& source, const std::vector<std::string>& mask,
            std::string* err) {
    if (mask.empty()) {
        *err = "update_mask must not be empty";
        return false;
    }

    std::vector<MaskOp> ops;
    for (auto& path : mask) {
        std::vector<std::string> segments = SplitDots(path);
        Value* parent = nullptr;
        if (!TargetParent(target, segments, path, &parent, err)) return false;
        std::string leaf = segments.back();
        Value value;
        bool found = SourceLookup(source, segments, &value);
        ops.push_back(MaskOp{parent, leaf, found, value});
    }

    for (auto& op : ops) {
        if (op.found) {
            op.parent->object[op.leaf] = op.value;
        } else {
            op.parent->object.erase(op.leaf);
        }
    }
    return true;
}
