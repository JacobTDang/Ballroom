#include <vector>

// CanCompleteCircuit returns the starting gas station index from which
// the circuit can be completed once, or -1 if no such start exists.
int CanCompleteCircuit(std::vector<int>& gas, std::vector<int>& cost) {
    int totalGas = 0, totalCost = 0;
    for (size_t i = 0; i < gas.size(); i++) {
        totalGas += gas[i];
        totalCost += cost[i];
    }
    if (totalGas < totalCost) return -1;

    int tank = 0;
    int start = 0;
    for (size_t i = 0; i < gas.size(); i++) {
        tank += gas[i] - cost[i];
        if (tank < 0) {
            tank = 0;
            start = (int)i + 1;
        }
    }
    return start;
}
