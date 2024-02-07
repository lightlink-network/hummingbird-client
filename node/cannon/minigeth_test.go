package cannon_test

import (
	"hummingbird/node/cannon"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() *cannon.MiniGeth {
	w := cannon.NewMiniGeth(&cannon.MiniGethOpts{
		Logger:  slog.Default().With("module", "witness_test"),
		NodeURL: "https://mainnet.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161",
		BaseDir: "/tmp/test_cannon",
	})

	// clear the base directory
	os.RemoveAll(w.Opts.BaseDir)
	os.MkdirAll(w.Opts.BaseDir, 0755)

	return w
}

func TestCannon_MiniGeth(t *testing.T) {
	minigeth := setup()
	defer os.RemoveAll(minigeth.Opts.BaseDir)

	err := minigeth.Process(19150042 / 2)
	assert.NoError(t, err)
}
