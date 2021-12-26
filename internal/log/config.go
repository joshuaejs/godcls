// config centralizes the log's configuration.

package log

// Config centralizes the configuration for the log
type Config struct {
	Segment struct {
		MaxStoreBytes uint64
		MaxIndexBytes uint64
		InitialOffset uint64
	}
}
