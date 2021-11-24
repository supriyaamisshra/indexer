package fetcher

import (
	"net/http"
)

type Fetcher interface {
	// fetch following / follower data
	FetchConnections(address string) ([]ConnectionEntry, error)
	// fetch user identity data
	FetchIdentity(address string) (IdentityEntryList, error)
}

type fetcher struct {
	httpClient *http.Client
}

var _ Fetcher = &fetcher{}

func NewFetcher() *fetcher {
	return &fetcher{
		httpClient: httpClient(),
	}
}
