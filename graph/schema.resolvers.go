package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"

	"github.com/seb-schulz/onegate/graph/model"
	"github.com/spf13/viper"
)

// CreateCredentialOptions is the resolver for the createCredentialOptions field.
func (r *queryResolver) CreateCredentialOptions(ctx context.Context) (*model.CreateCredentialOptions, error) {
	return &model.CreateCredentialOptions{
		Challenge: mustRandomEncodedBytes(32),
		Rp: model.RelyingParty{
			Name: viper.GetString("rp.name"),
			ID:   viper.GetString("rp.id"),
		},
		PubKeyCredParams: []*model.PubKeyCredParam{
			{Alg: -7, Type: "public-key"},
			{Alg: -257, Type: "public-key"},
		},
		UserID: mustRandomEncodedBytes(16),
	}, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
