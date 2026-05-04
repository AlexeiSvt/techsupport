// Package logic provides calculation engines for scoring user data.
// It includes thread-safe implementations for comparing various user attributes
// in high-concurrency environments.
package logic

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"

	"techsupport/core/internal/constants"
	"techsupport/core/internal/errors"
	"techsupport/core/pkg"
	"techsupport/core/pkg/models"
	logPkg "techsupport/log/pkg"
)

// Ensure FirstEmailCalculator implements the ScoreCalculator interface at compile time.
var _ pkg.ScoreCalculator = (*FirstEmailCalculator)(nil)

// FirstEmailCalculator handles the comparison and scoring logic for user email addresses.
// It is fully thread-safe and optimized for concurrent read access to its logger.
type FirstEmailCalculator struct {
	mu  sync.RWMutex // Protects the logger instance during runtime updates.
	log logPkg.Logger

	totalCalculations uint64 // Atomic counter for the total number of calculation attempts.
	matchCount        uint64 // Atomic counter for the total number of successful matches.
}

// NewFirstEmailCalculator initializes and returns a new FirstEmailCalculator with the given logger.
func NewFirstEmailCalculator(logger logPkg.Logger) *FirstEmailCalculator {
	return &FirstEmailCalculator{
		log: logger,
	}
}

// Calculate compares user and database email records to produce a scoring result.
// This method is designed to be executed within a goroutine and respects context cancellation.
// It extracts the first email from the UserClaimInfo slice and compares it to the email in the DBRecordInfoNode.
// The method returns a CalcResultInfoNode containing the outcome of the comparison,
// the applied weight, and any relevant comments or metadata.
func (c *FirstEmailCalculator) Calculate(ctx context.Context, user *models.UserNode, db models.DBRecordInfoNode, weights models.Weights) models.CalcResultInfoNode {

	// Increment total calculations atomically to avoid race conditions.
	atomic.AddUint64(&c.totalCalculations, 1)

	// Check for context cancellation before proceeding with the
	// calculation to allow for graceful shutdowns
	//and timeouts.
	select {
	case <-ctx.Done():
		return models.CalcResultInfoNode{
			Status:  "context_cancelled",
			Comment: ctx.Err().Error(),
		}
	default:
	}

	// Capture the logger instance under a Read-Lock to ensure thread-safe access.
	c.mu.RLock()
	logger := c.log
	c.mu.RUnlock()

	// Extract the first email from the UserClaimInfo slice, if available.
	var userEmail string
	var claimID string
	if len(user.UserClaimInfo) > 0 && user.UserClaimInfo[0] != nil {
		userEmail = user.UserClaimInfo[0].FirstEmail
		claimID = user.UserClaimInfo[0].ClaimID
	} else {
		return models.CalcResultInfoNode{
			Name:    "First Email Match",
			Code:    "first_email",
			Status:  "error_no_claims",
			Comment: "User node does not contain any ClaimInfoNode",
		}
	}

	dbEmail := db.FirstEmail

	// Perform the email comparison using a helper method that handles case-insensitivity and empty values.
	res := c.checkEmail(ctx, userEmail, dbEmail)

	// Increment match counter atomically if the comparison was successful.
	if res.Value == constants.IdealMatch {
		atomic.AddUint64(&c.matchCount, 1)
	}

	// Construct the final calculation result, including the applied weight and any relevant comments.
	calcRes := models.CalcResultInfoNode{
		Name:    "First Email Match",
		Code:    "first_email",
		Value:   res.Value,
		Weight:  weights.FirstEmail,
		Result:  res.Value * weights.FirstEmail,
		Status:  res.Status,
		Comment: res.Comment,
	}

	// Log the calculation result with relevant metadata for tracing and debugging.
	if logger != nil {
		logger.Infow("email score calculation finished",
			"status", calcRes.Status,
			"user_id", user.UserID,
			"claim_id", claimID,
		)
	}

	return calcRes
}

// checkEmail is an internal helper that performs case-insensitive string comparison.
// It handles empty input values and returns a raw check result with specific statuses.
func (c *FirstEmailCalculator) checkEmail(ctx context.Context, userEmail, dbEmail string) rawCheckResult {
	// Validate inputs to ensure both emails are present.
	if userEmail == "" || dbEmail == "" {
		// Warning log is omitted here as locking the logger inside a tight logic loop
		// could impact performance; logging is handled in the main Calculate method.
		return rawCheckResult{
			Value:   constants.NoMatch,
			Status:  errors.StatusNoData,
			Comment: "Missing email data for comparison",
		}
	}

	// Perform case-insensitive match using strings.EqualFold.
	if strings.EqualFold(userEmail, dbEmail) {
		return rawCheckResult{
			Value:   constants.IdealMatch,
			Status:  errors.StatusMatch,
			Comment: "Full case-insensitive match",
		}
	}

	return rawCheckResult{
		Value:   constants.NoMatch,
		Status:  errors.StatusNoMatch,
		Comment: "Emails do not match",
	}
}

// GetStats returns current metrics: total calculations performed and total matches found.
// This data is retrieved atomically without blocking the main execution flow.
func (c *FirstEmailCalculator) GetStats() (uint64, uint64) {
	return atomic.LoadUint64(&c.totalCalculations), atomic.LoadUint64(&c.matchCount)
}

// SetLogger allows for safe runtime updates of the logger instance.
// It uses a Write-Lock to ensure exclusive access during the pointer swap.
func (c *FirstEmailCalculator) SetLogger(newLogger logPkg.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.log = newLogger
}
