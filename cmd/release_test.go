package cmd

import (
	"testing"

	"github.com/goreleaser/goreleaser/pkg/config"
	"github.com/goreleaser/goreleaser/pkg/context"
	"github.com/stretchr/testify/require"
)

func TestRelease(t *testing.T) {
	_, back := setup(t)
	defer back()
	var cmd = newReleaseCmd()
	cmd.cmd.SetArgs([]string{"--snapshot", "--timeout=1m", "--parallelism=2", "--deprecated"})
	require.NoError(t, cmd.cmd.Execute())
}

func TestReleaseInvalidConfig(t *testing.T) {
	_, back := setup(t)
	defer back()
	createFile(t, "goreleaser.yml", "foo: bar")
	var cmd = newReleaseCmd()
	cmd.cmd.SetArgs([]string{"--snapshot", "--timeout=1m", "--parallelism=2", "--deprecated"})
	require.EqualError(t, cmd.cmd.Execute(), "yaml: unmarshal errors:\n  line 1: field foo not found in type config.Project")
}

func TestReleaseBrokenProject(t *testing.T) {
	_, back := setup(t)
	defer back()
	createFile(t, "main.go", "not a valid go file")
	var cmd = newReleaseCmd()
	cmd.cmd.SetArgs([]string{"--snapshot", "--timeout=1m", "--parallelism=2"})
	require.EqualError(t, cmd.cmd.Execute(), "failed to parse dir: .: main.go:1:1: expected 'package', found not")
}

func TestReleaseFlags(t *testing.T) {
	var setup = func(opts releaseOpts) *context.Context {
		return setupReleaseContext(context.New(config.Project{}), opts)
	}

	t.Run("snapshot", func(t *testing.T) {
		var ctx = setup(releaseOpts{
			snapshot: true,
		})
		require.True(t, ctx.Snapshot)
		require.True(t, ctx.SkipPublish)
		require.True(t, ctx.SkipPublish)
	})

	t.Run("skips", func(t *testing.T) {
		var ctx = setup(releaseOpts{
			skipPublish:  true,
			skipSign:     true,
			skipValidate: true,
		})
		require.True(t, ctx.SkipSign)
		require.True(t, ctx.SkipPublish)
		require.True(t, ctx.SkipPublish)
	})

	t.Run("parallelism", func(t *testing.T) {
		require.Equal(t, 1, setup(releaseOpts{
			parallelism: 1,
		}).Parallelism)
	})

	t.Run("notes", func(t *testing.T) {
		var notes = "foo.md"
		var header = "header.md"
		var footer = "footer.md"
		var ctx = setup(releaseOpts{
			releaseNotes:  notes,
			releaseHeader: header,
			releaseFooter: footer,
		})
		require.Equal(t, notes, ctx.ReleaseNotes)
		require.Equal(t, header, ctx.ReleaseHeader)
		require.Equal(t, footer, ctx.ReleaseFooter)
	})

	t.Run("rm dist", func(t *testing.T) {
		require.True(t, setup(releaseOpts{
			rmDist: true,
		}).RmDist)
	})
}
