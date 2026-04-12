package logic

import (
	"fmt"
	"strings"
	"techsupport/core/internal/constants"
	"techsupport/core/internal/models"
	"techsupport/core/pkg"
	"time"
)

var _ pkg.ScoreCalculator = (*DevicesCalculator)(nil)

type DevicesCalculator struct{}

func (c DevicesCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	res := checkAllDevices(user.UserClaim.Devices, db.Devices, db.RegDate)

	return models.CalcResult{
		Name:    "Device History Match",
		Code:    "devices_list",
		Value:   res.Value,
		Weight:  weights.Devices,
		Result:  res.Value * weights.Devices,
		Status:  res.Status,
		Comment: res.Comment,
	}
}

func checkAllDevices(userDevices []string, dbDevices []string, regDate time.Time) rawCheckResult {
	if len(userDevices) == 0 || len(dbDevices) == 0 {
		return rawCheckResult{Value: constants.NoMatch, Status: "no_data", Comment: "Device list is empty"}
	}

	deviceMap := make(map[string]struct{}, len(dbDevices))
	for _, d := range dbDevices {
		deviceMap[strings.ToLower(d)] = struct{}{}
	}

	matches := 0.0
	for _, u := range userDevices {
		lowU := strings.ToLower(u)
		if _, exists := deviceMap[lowU]; exists {
			matches++
			delete(deviceMap, lowU)
		}
	}

	if matches == 0 {
		return rawCheckResult{Value: constants.NoMatch, Status: "no_match", Comment: "No common devices found"}
	}

	matchRatio := matches / float64(len(dbDevices))
	if matchRatio > constants.IdealMatch {
		matchRatio = constants.IdealMatch
	}

	penalty := getDevicePenalty(regDate, len(dbDevices))
	finalValue := matchRatio * penalty

	comment := fmt.Sprintf("Matches: %.0f/%d. Ratio: %.2f. Penalty: %.2f", 
		matches, len(dbDevices), matchRatio, penalty)

	status := "partial"
	if finalValue >= constants.IdealMatch {
		status = "match"
	}

	return rawCheckResult{
		Value:   finalValue,
		Status:  status,
		Comment: comment,
	}
}

func getDevicePenalty(regDate time.Time, devicesCount int) float64 {
	yearsActive := time.Since(regDate).Hours() / constants.OneYearInHours
	if yearsActive < 1 {
		yearsActive = 1
	}

	allowedDevices := max(int(yearsActive*4), 10)

	if devicesCount > allowedDevices {
		return constants.MostlyMatch
	}

	return constants.IdealMatch
}