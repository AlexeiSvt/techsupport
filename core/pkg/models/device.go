package models

// Device represents a unique hardware entity within the graph ecosystem.
// It acts as a central pivot point (hub) that allows the engine to link 
// multiple sessions, transactions, and user claims to the same physical hardware.
type Device struct {
    // DeviceID is the unique hardware fingerprint (UUID, IMEI, or secure hash).
    // In Neo4j, this is the primary lookup key for the Device node.
    DeviceID string `json:"device_id" cypher:"id"`

    // Model is the human-readable marketing name of the hardware (e.g., "iPhone 15 Pro").
    // Used for consistency checks against UserClaim.FirstDeviceID.
    Model string `json:"model" cypher:"model"`
}