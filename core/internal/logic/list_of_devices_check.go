package logic

import (
    "fmt"
    "strings"
    "techsupport/core/internal/constants"
    "techsupport/core/internal/errors"
    "techsupport/core/pkg/models"
    "techsupport/core/pkg"
    logPkg "techsupport/log/pkg"
    "time"
)

var _ pkg.ScoreCalculator = (*DevicesCalculator)(nil)

type DevicesCalculator struct {
    Log logPkg.Logger
}

func (c *DevicesCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
    if c.Log != nil {
        c.Log.Debugw("calculating all devices match", 
            "user_devices_count", len(user.UserClaim.Devices), 
            "db_devices_count", len(db.Devices),
        )
    }

    res := c.checkAllDevices(user.UserClaim.Devices, db.Devices, db.RegDate)

    calcRes := models.CalcResult{
        Name:    "Device History Match",
        Code:    "devices_list",
        Value:   res.Value,
        Weight:  weights.Devices,
        Result:  res.Value * weights.Devices,
        Status:  res.Status,
        Comment: res.Comment,
    }

    if c.Log != nil {
        c.Log.Infow("device history calculation finished", 
            "status", calcRes.Status, 
            "result", calcRes.Result,
        )
    }

    return calcRes
}

func (c *DevicesCalculator) checkAllDevices(userDevices []string, dbDevices []string, regDate time.Time) rawCheckResult {
    if len(userDevices) == 0 || len(dbDevices) == 0 {
        if c.Log != nil {
            c.Log.Warnw(errors.ErrEmptyDeviceList.Error(), 
                "u_count", len(userDevices), 
                "db_count", len(dbDevices),
            )
        }
        return rawCheckResult{
            Value:   constants.NoMatch, 
            Status:  errors.StatusNoData, 
            Comment: "Device list is empty",
        }
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
        return rawCheckResult{
            Value:   constants.NoMatch, 
            Status:  errors.StatusNoMatch, 
            Comment: "No common devices found",
        }
    }

    matchRatio := matches / float64(len(dbDevices))
    if matchRatio > constants.IdealMatch {
        matchRatio = constants.IdealMatch
    }

    penalty := c.getDevicePenalty(regDate, len(dbDevices))
    finalValue := matchRatio * penalty

    comment := fmt.Sprintf("Matches: %.0f/%d. Ratio: %.2f. Penalty: %.2f", 
        matches, len(dbDevices), matchRatio, penalty)

    status := errors.StatusPartial
    if finalValue >= constants.IdealMatch {
        status = errors.StatusMatch
    }

    if c.Log != nil {
        c.Log.Debugw("device match logic details", 
            "matches", matches, 
            "ratio", matchRatio, 
            "penalty", penalty,
        )
    }

    return rawCheckResult{
        Value:   finalValue,
        Status:  status,
        Comment: comment,
    }
}

func (c *DevicesCalculator) getDevicePenalty(regDate time.Time, devicesCount int) float64 {
    durationSinceReg := time.Since(regDate)
    
    if durationSinceReg < 0 {
        if c.Log != nil {
            c.Log.Errorw(errors.ErrFutureRegDate.Error(), "reg_date", regDate)
        }
        return constants.IdealMatch
    }

    yearsActive := durationSinceReg.Hours() / constants.OneYearInHours
    if yearsActive < 1 {
        yearsActive = 1
    }

    allowedDevices := 10
    calculatedAllowed := int(yearsActive * 4)
    if calculatedAllowed > allowedDevices {
        allowedDevices = calculatedAllowed
    }

    if devicesCount > allowedDevices {
        if c.Log != nil {
            c.Log.Debugw("applying device count penalty", 
                "count", devicesCount, 
                "allowed", allowedDevices,
            )
        }
        return constants.MostlyMatch
    }

    return constants.IdealMatch
}