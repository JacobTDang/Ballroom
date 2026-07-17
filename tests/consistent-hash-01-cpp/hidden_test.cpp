#include "solution.cpp"

#include <cstdio>
#include <map>

int main() {
    // Deterministic + balance.
    {
        Ring r(100);
        r.AddNode("node-a");
        r.AddNode("node-b");
        r.AddNode("node-c");

        std::map<std::string, int> counts;
        const int keys = 10000;
        for (int i = 0; i < keys; i++) {
            std::string k = "key-" + std::to_string(i);
            std::string owner = r.Lookup(k);
            if (owner.empty() || r.Lookup(k) != owner) {
                fprintf(stderr, "unstable or empty lookup for %s\n", k.c_str());
                return 1;
            }
            counts[owner]++;
        }
        for (const char* node : {"node-a", "node-b", "node-c"}) {
            double share = (double)counts[node] / keys;
            if (share < 0.10 || share > 0.60) {
                fprintf(stderr, "%s owns %.0f%%, want 10-60%% with 100 vnodes\n", node, share * 100);
                return 1;
            }
        }

        // Minimal remap + exact restore.
        std::map<std::string, std::string> before;
        for (int i = 0; i < keys; i++) {
            std::string k = "key-" + std::to_string(i);
            before[k] = r.Lookup(k);
        }
        r.AddNode("node-d");
        int moved = 0;
        for (auto& [k, owner] : before) {
            if (r.Lookup(k) != owner) moved++;
        }
        if (moved * 2 >= keys) {
            fprintf(stderr, "adding one node moved %d/%d keys -- %%N rehashing\n", moved, keys);
            return 1;
        }
        if (moved * 20 < keys) {
            fprintf(stderr, "adding a node moved only %d/%d keys\n", moved, keys);
            return 1;
        }
        r.RemoveNode("node-d");
        for (auto& [k, owner] : before) {
            if (r.Lookup(k) != owner) {
                fprintf(stderr, "removal did not restore the original mapping for %s\n", k.c_str());
                return 1;
            }
        }
    }

    // Empty ring.
    {
        Ring r(100);
        if (!r.Lookup("anything").empty()) {
            fprintf(stderr, "empty ring lookup returned a node\n");
            return 1;
        }
    }

    printf("all assertions passed\n");
    return 0;
}
