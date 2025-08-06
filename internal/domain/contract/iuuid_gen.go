package contract

// UUIDGenerator defines the interface for generating UUIDs.
type IUUIDGenerator interface {
	NewUUID() string
}
