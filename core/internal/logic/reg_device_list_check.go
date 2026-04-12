package logic

import (
	"strings"
	"techsupport/core/internal/models"
	"techsupport/core/internal/constants"
	"time"
)

func CalculateScoreForFirstDevice(userFirstDevice, dbFirstDevice string, weights models.Weights) float64 {

	if userFirstDevice == "" || dbFirstDevice == "" {
		return constants.NoMatch
	}

	if strings.EqualFold(userFirstDevice, dbFirstDevice) {
		return weights.FirstDevice * constants.IdealMatch
	}
	return constants.NoMatch
}

func CalculateScoreForDevices(userDevices []string, dbRecord models.DBRecord, weights models.Weights) float64 {
 if len(userDevices) == 0 || len(dbRecord.Devices) == 0 {
  return constants.NoMatch
 }

 deviceMap := make(map[string]struct{}, len(dbRecord.Devices))
 for _, d := range dbRecord.Devices {
  deviceMap[strings.ToLower(d)] = struct{}{}
 }

 matches := 0.0
 for _, u := range userDevices {
  if _, exists := deviceMap[strings.ToLower(u)]; exists {
   matches++
   delete(deviceMap, strings.ToLower(u))
  }
 }

 if matches == 0.0 {
  return constants.NoMatch
 }

 matchRatio := float64(matches) / float64(len(dbRecord.Devices))
 penalty := GetDevicePenalty(dbRecord.RegDate, len(dbRecord.Devices))

 if matchRatio > constants.IdealMatch {
  matchRatio = constants.IdealMatch
 }

 return weights.Devices * matchRatio* penalty
}


func GetDevicePenalty(regDate time.Time, devicesCount int) float64 {
    yearsActive := time.Since(regDate).Hours() / constants.OneYearInHours
    if yearsActive < 1 {
        yearsActive = 1
    }

    allowedDevices := max(int(yearsActive * 4), 10)

    if devicesCount > allowedDevices {
        return constants.MostlyMatch
    }

    return constants.IdealMatch
}