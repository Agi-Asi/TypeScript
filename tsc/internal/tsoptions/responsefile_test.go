package tsoptions_test

import (
	"testing"

	"github.com/microsoft/TypeScript/tsc/internal/tsoptions"
	"github.com/microsoft/TypeScript/tsc/internal/tsoptions/tsoptionstest"
	"gotest.tools/v3/assert"
)

// GH#62373: `tsc -b @scoped/package/src` must treat "@"-prefixed project paths
// as paths, not as response files, when the argument does not name a file.
func TestParseBuildCommandLineScopedPackagePath(t *testing.T) {
	t.Parallel()
	host := tsoptionstest.NewVFSParseConfigHost(map[string]string{
		"/home/src/@scoped/package/src/tsconfig.json": `{"compilerOptions":{"composite":true}}`,
	}, "/home/src", true)
	parsed := tsoptions.ParseBuildCommandLine([]string{"@scoped/package/src"}, host)
	assert.Equal(t, len(parsed.Errors), 0)
	assert.DeepEqual(t, parsed.Projects, []string{"@scoped/package/src"})
}

// Response files must still be expanded when the referenced file exists.
func TestParseBuildCommandLineResponseFileStillExpands(t *testing.T) {
	t.Parallel()
	host := tsoptionstest.NewVFSParseConfigHost(map[string]string{
		"/home/src/buildargs.txt": "--verbose --force",
	}, "/home/src", true)
	parsed := tsoptions.ParseBuildCommandLine([]string{"@buildargs.txt"}, host)
	assert.Equal(t, len(parsed.Errors), 0)
	assert.Assert(t, parsed.BuildOptions.Verbose.IsTrue())
	assert.Assert(t, parsed.BuildOptions.Force.IsTrue())
}

// A missing slash-containing "@" path is treated as an ordinary command-line
// argument rather than a response-file error; this is what "@"-prefixed scoped
// project paths need.
func TestParseCommandLineScopedPathBecomesArgument(t *testing.T) {
	t.Parallel()
	host := tsoptionstest.NewVFSParseConfigHost(map[string]string{}, "/home/src", true)
	parsed := tsoptions.ParseCommandLine([]string{"@scoped/not-a-response-file.ts"}, host)
	assert.Equal(t, len(parsed.Errors), 0)
	assert.DeepEqual(t, parsed.FileNames(), []string{"@scoped/not-a-response-file.ts"})
}

// Response files in normal (non-build) mode must still be expanded.
func TestParseCommandLineResponseFileStillExpands(t *testing.T) {
	t.Parallel()
	host := tsoptionstest.NewVFSParseConfigHost(map[string]string{
		"/home/src/args.txt":   "--strict @nested.txt",
		"/home/src/nested.txt": "0.ts",
	}, "/home/src", true)
	parsed := tsoptions.ParseCommandLine([]string{"@args.txt"}, host)
	assert.Equal(t, len(parsed.Errors), 0)
	assert.Assert(t, parsed.CompilerOptions().Strict.IsTrue())
	assert.DeepEqual(t, parsed.FileNames(), []string{"0.ts"})
}
