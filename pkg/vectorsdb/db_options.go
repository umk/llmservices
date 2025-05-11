package vectorsdb

type dbConfig struct {
	repackPercent int
}

// Option configures a Database.
type Option func(*dbConfig)

// WithRepackPercent sets the percentage of deleted items that triggers a repack.
func WithRepackPercent(percentage int) Option {
	return func(cfg *dbConfig) {
		cfg.repackPercent = percentage
	}
}
