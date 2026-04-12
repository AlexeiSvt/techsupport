package logic

import (
	"strings"
	"techsupport/core/internal/constants"
	"techsupport/core/internal/models"
	"techsupport/core/pkg"
)

var _ pkg.ScoreCalculator = (*FirstDeviceCalculator)(nil)

type FirstDeviceCalculator struct{}

func (c FirstDeviceCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	res := checkFirstDevice(user.UserClaim.FirstDevice, db.FirstDevice)

	return models.CalcResult{
		Name:    "First Device Match",
		Code:    "first_device",
		Value:   res.Value,
		Weight:  weights.FirstDevice,
		Result:  res.Value * weights.FirstDevice,
		Status:  res.Status,
		Comment: res.Comment,
	}
}

func checkFirstDevice(userDev, dbDev string) rawCheckResult {
	if userDev == "" || dbDev == "" {
		return rawCheckResult{Value: constants.NoMatch, Status: "no_data", Comment: "Missing first device data"}
	}

	if strings.EqualFold(userDev, dbDev) {
		return rawCheckResult{Value: constants.IdealMatch, Status: "match", Comment: "First device matches exactly"}
	}

	return rawCheckResult{Value: constants.NoMatch, Status: "no_match", Comment: "First device mismatch"}
}
