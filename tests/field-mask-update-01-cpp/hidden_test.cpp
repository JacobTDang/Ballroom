#include "solution.cpp"

#include <cstdio>

#define CHECK(cond, msg)                       \
    if (!(cond)) {                              \
        fprintf(stderr, "FAILED: %s\n", msg);   \
        return 1;                               \
    }

static bool Has(const Value& v, const std::string& key) {
    return v.IsObject() && v.object.count(key) > 0;
}

static std::string Scalar(const Value& v, const std::string& key) {
    auto it = v.object.find(key);
    return it == v.object.end() ? "<missing>" : it->second.scalar;
}

int main() {
    // Sibling untouched: updating one leaf never disturbs siblings.
    {
        Value target = Value::Obj({
            {"name", Value::Str("Alice")},
            {"address", Value::Obj({{"city", Value::Str("SF")}, {"zip", Value::Str("94107")}})},
            {"tags", Value::Obj({{"vip", Value::Str("true")}})},
        });
        Value source = Value::Obj({{"address", Value::Obj({{"city", Value::Str("NYC")}})}});
        std::string err;

        CHECK(Update(&target, source, {"address.city"}, &err), "Update failed");
        CHECK(Scalar(target.object["address"], "city") == "NYC", "city wasn't updated");
        CHECK(Scalar(target.object["address"], "zip") == "94107", "sibling zip was disturbed");
        CHECK(Scalar(target, "name") == "Alice", "an untouched top-level sibling changed");
    }

    // Multi-path: several mask entries, including ones sharing a parent.
    {
        Value target = Value::Obj({
            {"name", Value::Str("Alice")},
            {"address", Value::Obj({{"city", Value::Str("SF")}, {"zip", Value::Str("94107")}})},
            {"tags", Value::Obj({{"vip", Value::Str("true")}})},
        });
        Value source = Value::Obj({
            {"name", Value::Str("Bob")},
            {"address", Value::Obj({{"zip", Value::Str("10001")}})},
            {"tags", Value::Obj({{"vip", Value::Str("false")}})},
        });
        std::string err;

        CHECK(Update(&target, source, {"name", "address.zip", "tags.vip"}, &err), "Update failed");
        CHECK(Scalar(target, "name") == "Bob", "name wasn't updated");
        CHECK(Scalar(target.object["address"], "zip") == "10001", "zip wasn't updated");
        CHECK(Scalar(target.object["address"], "city") == "SF", "city sibling shouldn't move");
        CHECK(Scalar(target.object["tags"], "vip") == "false", "tags.vip wasn't updated");
    }

    // Clear via omission: a masked path absent from source deletes the field.
    {
        Value target = Value::Obj({
            {"name", Value::Str("Alice")},
            {"address", Value::Obj({{"city", Value::Str("SF")}, {"zip", Value::Str("94107")}})},
        });
        Value source = Value::Obj({{"address", Value::Obj({})}});
        std::string err;

        CHECK(Update(&target, source, {"address.zip", "name"}, &err), "Update failed");
        CHECK(!Has(target.object["address"], "zip"), "address.zip should have been cleared");
        CHECK(Scalar(target.object["address"], "city") == "SF", "city sibling shouldn't move");
        CHECK(!Has(target, "name"), "name should have been cleared");
    }

    // Missing intermediate: target has no "address" at all.
    {
        Value target = Value::Obj({{"name", Value::Str("Alice")}});
        Value source = Value::Obj({{"address", Value::Obj({{"city", Value::Str("NYC")}})}});
        std::string err;

        CHECK(!Update(&target, source, {"address.city"}, &err),
              "expected an error for a missing intermediate");
        CHECK(err.find("address") != std::string::npos, "error doesn't name the offending path");
        CHECK(!Has(target, "address"), "target changed despite the error");
        CHECK(Scalar(target, "name") == "Alice", "target changed despite the error");
    }

    // Scalar intermediate: target.address exists but isn't an object.
    {
        Value target = Value::Obj({{"name", Value::Str("Alice")}, {"address", Value::Str("not-an-object")}});
        Value source = Value::Obj({{"address", Value::Obj({{"city", Value::Str("NYC")}})}});
        std::string err;

        CHECK(!Update(&target, source, {"address.city"}, &err),
              "expected an error for a scalar intermediate");
        CHECK(err.find("address") != std::string::npos, "error doesn't name the offending path");
        CHECK(target.object["address"].scalar == "not-an-object", "target changed despite the error");
    }

    // Empty mask is an error.
    {
        Value target = Value::Obj({{"name", Value::Str("Alice")}});
        Value source = Value::Obj({});
        std::string err;

        CHECK(!Update(&target, source, {}, &err), "expected an error for an empty mask");
        CHECK(Scalar(target, "name") == "Alice", "target changed despite the error");
    }

    printf("all assertions passed\n");
    return 0;
}
