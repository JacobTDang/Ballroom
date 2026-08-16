# Capital One OA roadmap

Capital One screens software engineers with CodeSignal's General Coding
Assessment: one proctored 70-minute sitting, 4 questions, camera/mic/
screen recorded, scored 200–600. The difficulty ramps: two easy
warm-ups, an implementation-heavy third question, an algorithmic
medium fourth. Dynamic programming is rare; speed and clean
implementation matter more than clever algorithms. No public pass
cutoff exists — circulating "750+" numbers are from a scale CodeSignal
retired in 2023.

## The pools

The capital-one track mirrors that shape with three pools:

- **Warm-ups** (`c1-warmup-*`, easy, 10 min) — array/string
  manipulation. The real Q1/Q2 should cost you well under 10 minutes
  each; they exist to bank time for the back half.
- **Grid** (`c1-grid-*`, medium, 25 min) — 2-D matrix walks,
  simulations, floods, rotations. The Q3 slot is rarely conceptually
  hard, but it's fiddly: off-grid checks, simultaneous updates,
  non-square grids. Most lost points live here.
- **Algo** (`c1-algo-*`, medium, 20 min) — hashmap counting,
  stacks, greedy, parsing. The Q4 slot is a standard LeetCode-medium
  shape.

## Suggested order

1. Warm-ups until each takes under 8 minutes.
2. Grid problems until a clean pass inside 25 minutes feels routine.
3. Algo problems for the Q4 archetypes.
4. Then **Mock** from the main menu: 2 warm-ups + 1 grid + 1 algo
   drawn fresh, against one 70-minute clock that never pauses —
   pacing across questions is the skill the single-problem drills
   can't teach.

The mock is forward-only (no returning to an earlier question), which
is stricter than the real CodeSignal UI — treat it as pacing training,
not a faithful simulator.
