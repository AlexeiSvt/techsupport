package logic

import (
	"strings"
	"techsupport/core/internal/constants"
	"techsupport/core/internal/errors"
	"techsupport/core/internal/models"
	"techsupport/core/pkg"
	logPkg "techsupport/log/pkg"
)

var _ pkg.ScoreCalculator = (*FirstDeviceCalculator)(nil)

type FirstDeviceCalculator struct {
	Log logPkg.Logger
}

func (c *FirstDeviceCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	userDev := user.UserClaim.FirstDevice
	dbDev := db.FirstDevice

	if c.Log != nil {
		c.Log.Debugw("checking first device match", "user_device", userDev, "db_device", dbDev)
	}

	res := c.checkFirstDevice(userDev, dbDev)

	calcRes := models.CalcResult{
		Name:    "First Device Match",
		Code:    "first_device",
		Value:   res.Value,
		Weight:  weights.FirstDevice,
		Result:  res.Value * weights.FirstDevice,
		Status:  res.Status,
		Comment: res.Comment,
	}

	if c.Log != nil {
		c.Log.Infow("device score calculated",
			"status", calcRes.Status,
			"result", calcRes.Result,
		)
	}

	return calcRes
}

func (c *FirstDeviceCalculator) checkFirstDevice(userDev, dbDev string) rawCheckResult {
	u := strings.TrimSpace(userDev)
	d := strings.TrimSpace(dbDev)

	if u == "" || d == "" {
		if c.Log != nil {
			c.Log.Warnw(errors.ErrEmptyDeviceData.Error(), "u_dev_empty", u == "", "db_dev_empty", d == "")
		}

		return rawCheckResult{
			Value:   constants.NoMatch,
			Status:  errors.StatusNoData,
			Comment: "Missing first device data",
		}
	}

	if strings.EqualFold(u, d) {
		return rawCheckResult{
			Value:   constants.IdealMatch,
			Status:  errors.StatusMatch,
			Comment: "First device matches exactly",
		}
	}

	if c.Log != nil {
		c.Log.Debugw(errors.ErrDeviceNoMatch.Error(), "u_dev", u, "db_dev", d)
	}

	return rawCheckResult{
		Value:   constants.NoMatch,
		Status:  errors.StatusNoMatch,
		Comment: "First device mismatch",
	}
}