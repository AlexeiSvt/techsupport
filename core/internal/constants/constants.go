package constants

const (
	Insolvent_Weight_RegDate          = 12.5 // 12.5 is the weight for registration date for Insolvent users
	Insolvent_Weight_RegCountry       = 5.0  // 5.0 is the weight for registration country for Insolvent users
	Insolvent_Weight_RegCity          = 12.5 // 12.5 is the weight for registration city for Insolvent users
	Insolvent_Weight_FirstEmail       = 11.5 // 11.5 is the weight for first email for Insolvent users
	Insolvent_Weight_Phone            = 16.0 // 16.0 is the weight for phone for Insolvent users
	Insolvent_Weight_FirstDevice      = 17.5 // 17.5 is the weight for first device for Insolvent users
	Insolvent_Weight_Devices          = 25.0 // 25.0 is the weight for devices for Insolvent users
	Insolvent_Weight_FirstTransaction = 0.0  // 0.0 is the weight for first transaction for Insolvent users, as they don't have transactions
)

const (
	Solvent_Weight_RegDate          = 7.5  // 7.5 is the weight for registration date for Solvent users
	Solvent_Weight_RegCountry       = 5.0  // 5.0 is the weight for registration country for Solvent users
	Solvent_Weight_RegCity          = 12.5 // 12.5 is the weight for registration city for Solvent users
	Solvent_Weight_FirstEmail       = 11.5 // 11.5 is the weight for first email for Solvent users
	Solvent_Weight_Phone            = 16.0 // 16.0 is the weight for phone for Solvent users
	Solvent_Weight_FirstDevice      = 12.5 // 12.5 is the weight for first device for Solvent users
	Solvent_Weight_Devices          = 15.0 // 15.0 is the weight for devices for Solvent users
	Solvent_Weight_FirstTransaction = 20.0 // 20.0 is the weight for first transaction for Solvent users, as they have transactions
)

const (
	IdealMatch   = 1.0 // 1.0 is the score for an ideal match
	MostlyMatch  = 0.7 // 0.7 is the threshold for a mostly match
	PartialMatch = 0.5 // 0.5 is the score for a partial match
	NoMatch      = 0.0 // 0.0 is the score for no match
)

const (
	AvgAmountOfHoursInMonth = 24 * 30.44 // Average number of hours in a month
	ToleranceHours          = 24.0 * 5   // 5 days tolerance in hours
	OneYearInHours          = 24.0 * 365 // One year in hours
	MinIntervalHours = 2.0 // 2.0 is the minimum interval in hours
	IdealMatchofMonths   = 2 // Ideal match if the registration dates are within 2 months
	OneYearofMonths = 12 // One year is 12 months
	PartialMatchofMonths = 4 // Partial match if the registration dates are within 4 months
)

const (
	LocMatch    = 0.3  // 30% of the score is based on location
	DeviceMatch = 0.25 // 25% of the score is based on device information
	IPMatch     = 0.15 // 15% of the score is based on IP address
)

const (
	MaxTransactionScore = 20.0 // 20.0 is the maximum score for a transaction
	CountryScore        = 6.0  // 6.0 is the score for a country match
	CityScore           = 4.0  // 4.0 is the score for a city match
	DeviceScore         = 7.0  // 7.0 is the score for a device match
	IPScore             = 3.0  // 3.0 is the score for an IP match
	MinScoreForPartialMatch = 10.0 // 10.0 is the minimum score for a partial match for transactions
)

const (
	BruteforcePenalty   = 20.0 // 20.0 is the penalty for bruteforce attempts
	NewDevicesThreshold = 3.0  // 3.0 is the threshold for new devices
)

const (
	SuddenMultiplier       = 5.0 // 5.0 is the multiplier for sudden high donations
	FirstDonationThreshold = 5.0 // 5.0 is the threshold for the first donation for F2P users
)

const MinLen = 7 // 7 is the minimum length for phone numbers to be considered for partial matching

const (
	FullPenalty = 100
	ForVPN = 70
	ForProxy = 70
	ForHosting = 40
	ForDatacenter = 40
)

const ApiBaseURL = "https://api.ipapi.is"