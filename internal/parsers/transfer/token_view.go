package transfer

import (
	"github.com/baking-bad/bcdhub/internal/elastic"
	"github.com/baking-bad/bcdhub/internal/models"
	"github.com/baking-bad/bcdhub/internal/models/tzip"
	"github.com/jinzhu/gorm"
)

// TokenKey -
type TokenKey struct {
	Address    string
	Network    string
	Entrypoint string
}

// TokenEvents -
type TokenEvents map[TokenKey]tzip.EventImplementation

// NewTokenViews -
func NewTokenViews(es elastic.IElastic) (TokenEvents, error) {
	tokens, err := es.GetTZIPWithViews()
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	views := make(TokenEvents)
	for _, token := range tokens {
		if len(token.Events) == 0 {
			continue
		}

		for _, view := range token.Events {
			for _, implementation := range view.Implementations {
				for _, entrypoint := range implementation.MichelsonParameterView.Entrypoints {
					views[TokenKey{
						Address:    token.Address,
						Network:    token.Network,
						Entrypoint: entrypoint,
					}] = implementation
				}
			}
		}
	}

	return views, nil
}

// Get -
func (events TokenEvents) Get(address, network, entrypoint string) (tzip.EventImplementation, bool) {
	view, ok := events[TokenKey{
		Address:    address,
		Network:    network,
		Entrypoint: entrypoint,
	}]
	return view, ok
}

// GetByOperation -
func (events TokenEvents) GetByOperation(operation models.Operation) (tzip.EventImplementation, bool) {
	event, ok := events[TokenKey{
		Address:    operation.Destination,
		Network:    operation.Network,
		Entrypoint: operation.Entrypoint,
	}]
	return event, ok
}