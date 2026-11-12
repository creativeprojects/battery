package battery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBatteries(t *testing.T) {
	tests := []struct {
		name     string
		testdata string
		found    int
		state    AgnosticState
	}{
		{
			name:     "mac os 26 - no battery installed",
			testdata: "macos26_nobattery.xml",
			found:    0,
		},
		{
			name:     "mac os 27 - charging",
			testdata: "macos27_charging.xml",
			found:    1,
			state:    Charging,
		},
		{
			name:     "mac os 27 - discharging",
			testdata: "macos27_discharging.xml",
			found:    1,
			state:    Discharging,
		},
		{
			name:     "mac os 27 - full",
			testdata: "macos27_full.xml",
			found:    1,
			state:    Full,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := os.ReadFile(filepath.Join("testdata", tt.testdata))
			require.NoError(t, err)

			darwinBatteries, err := parseBatteries(raw)
			require.NoError(t, err)

			batteries := convertBatteries(darwinBatteries)
			require.Len(t, batteries, tt.found)
			if tt.found == 1 {
				assert.Equal(t, tt.state, batteries[0].State.Raw)
			}
		})
	}
}
