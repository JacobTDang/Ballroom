package orchestrator

import "github.com/JacobTDang/Ballroom/internal/config"

// dockerBuildRoot resolves the directory to use as docker build's -f path
// base and build context: cfgRoot itself if it looks like a real ballroom
// checkout, else the last root that did (cached from a previous
// successful run). Without this, `ballroom` launched from PATH outside
// the checkout (its usual location once installed via `go install`)
// fails with a raw, confusing error straight from docker — "lstat
// docker: no such file or directory" — instead of either finding the
// real checkout or explaining clearly what's wrong.
//
// A thin delegate to config.ResolveRoot (issue #255): this package and
// config.Load used to each carry their own copy of the same
// looks-like-a-checkout / cached-root-fallback logic — config.Load's
// copy applied no validation at all, which is what let an installed
// binary launched outside the checkout silently resolve
// ExercisesDir/TestsDir/DataDir under the wrong directory (empty
// picker). Now there's exactly one implementation, in config (which
// this package already imports; config cannot import back).
func dockerBuildRoot(cfgRoot string) (string, error) {
	return config.ResolveRoot(cfgRoot)
}
