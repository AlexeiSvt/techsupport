package tests

import (
	"techsupport/core/internal/scoring/logic"

	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCalculateScoreForFirstDevice_AllCases(t *testing.T) {
	cases := []struct {
		name      string
		user      string
		db        string
		isDonator bool
		expected  float64
	}{

		{"F2P Identical", "iPhone 12", "iPhone 12", false, 17.5},
		{"F2P Different", "iPhone 12", "Samsung Galaxy", false, 0.0},
		{"F2P EmptyBoth", "", "", false, 0.0},
		{"F2P EmptyUser", "", "iPhone 12", false, 0.0},
		{"F2P EmptyDB", "iPhone 12", "", false, 0.0},
		{"F2P CaseInsensitive", "iphone 12", "IPHONE 12", false, 17.5},
		{"F2P DifferentModels", "iPhone 12", "iPhone 13", false, 0.0},
		{"F2P WithSpaces", "iPhone 12", "iPhone 12 ", false, 0.0},
		{"F2P LeadingSpace", " iPhone 12", "iPhone 12", false, 0.0},
		{"F2P Pixel", "Google Pixel 6", "Google Pixel 6", false, 17.5},
		{"F2P PixelDifferent", "Google Pixel 6", "Google Pixel 7", false, 0.0},
		{"F2P iPad", "iPad Air", "iPad Air", false, 17.5},
		{"F2P LongDeviceName", "Samsung Galaxy S21 Ultra Premium Edition", "Samsung Galaxy S21 Ultra Premium Edition", false, 17.5},
		{"F2P WithNumbers", "OnePlus 9 Pro", "OnePlus 9 Pro", false, 17.5},
		{"F2P DifferentNumbers", "OnePlus 9 Pro", "OnePlus 8 Pro", false, 0.0},
		{"F2P Generic", "Mobile Device", "Mobile Device", false, 17.5},
		{"F2P Unicode", "Смартфон", "Смартфон", false, 17.5},
		{"F2P SingleChar", "A", "A", false, 17.5},

		{"P2W Identical", "iPhone 12", "iPhone 12", true, 12.5},
		{"P2W Different", "iPhone 12", "Samsung Galaxy", true, 0.0},
		{"P2W EmptyBoth", "", "", true, 0.0},
		{"P2W EmptyUser", "", "iPhone 12", true, 0.0},
		{"P2W EmptyDB", "iPhone 12", "", true, 0.0},
		{"P2W CaseInsensitive", "iphone 12", "IPHONE 12", true, 12.5},
		{"P2W DifferentModels", "iPhone 12", "iPhone 13", true, 0.0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := logic.CalculateScoreForFirstDevice(c.user, c.db, logic.GetWeights(c.isDonator))
			assert.Equal(t, c.expected, result, "Test failed: %s", c.name)
		})
	}
}
