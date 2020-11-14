package internal

import "fpl-strategy-tester/internal/database"

// Resolver is the entry-point for accessing the football data
type Resolver struct {
	Database *database.Resolver
}

// NewResolver creates and returns an empty Resolver
func NewResolver() *Resolver {
	return &Resolver{}
}

// ResolveFootballDB returns or initiates a new database connection
func (r *Resolver) ResolveDatabase() *database.Resolver {
	if r.Database == nil {
		repo := database.NewResolver()
		repo.ResolveFPLDB()
		repo.ResolveMySQLQueryBuilder()
		r.Database = repo
	}
	return r.Database
}
