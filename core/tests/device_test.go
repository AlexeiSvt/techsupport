package tests

import (
	"techsupport/core/internal/logic"
	"techsupport/core/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirstDeviceCalculator_AllCases(t *testing.T) {
	calc := logic.FirstDeviceCalculator{}

	cases := []struct {
		name      string
		userDev   string
		dbDev     string
		isDonator bool
		expected  float64
	}{
		{"F2P Identical", "iPhone 12", "iPhone 12", false, 17.5},
		{"F2P Different", "iPhone 12", "Samsung Galaxy", false, 0.0},
		{"F2P EmptyBoth", "", "", false, 0.0},
		{"F2P CaseInsensitive", "iphone 12", "IPHONE 12", false, 17.5},
		{"F2P Pixel", "Google Pixel 6", "Google Pixel 6", false, 17.5},
		{"F2P LongDeviceName", "Samsung Galaxy S21 Ultra Premium Edition", "Samsung Galaxy S21 Ultra Premium Edition", false, 17.5},

		{"P2W Identical", "iPhone 12", "iPhone 12", true, 12.5},
		{"P2W Different", "iPhone 12", "Samsung Galaxy", true, 0.0},
		{"P2W EmptyBoth", "", "", true, 0.0},
		{"P2W CaseInsensitive", "iphone 12", "IPHONE 12", true, 12.5},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			weights := logic.GetWeights(c.isDonator)

			user := models.UserData{
				UserClaim: models.UserClaim{FirstDevice: c.userDev},
			}
			db := models.DBRecord{
				FirstDevice: c.dbDev,
			}

			result := calc.Calculate(user, db, weights)

			assert.InDelta(t, c.expected, result.Result, 0.001, "Test failed: %s", c.name)

			if c.expected > 0 {
				assert.Equal(t, "match", result.Status)
			} else if c.userDev == "" || c.dbDev == "" {
				assert.Equal(t, "no_data", result.Status)
			} else {
				assert.Equal(t, "no_match", result.Status)
			}
		})
	}
}