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
		state    AgnosticState
	}{
		{
			name:     "mac os 27 - charging",
			testdata: "macos27_charging.xml",
			state:    Charging,
		},
		{
			name:     "mac os 27 - discharging",
			testdata: "macos27_discharging.xml",
			state:    Discharging,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := os.ReadFile(filepath.Join("testdata", tt.testdata))
			require.NoError(t, err)
			darwinBatteries, err := parseBatteries(raw)
			require.NoError(t, err)
			require.Len(t, darwinBatteries, 1)
			battery := convertBattery(darwinBatteries[0])
			assert.Equal(t, tt.state, battery.State.Raw)
		})
	}
}
