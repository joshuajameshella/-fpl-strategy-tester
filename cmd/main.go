package main

import "fpl-strategy-tester/internal"

func main() {
	resolver := internal.NewResolver()
	resolver.ResolveDatabase()
}
