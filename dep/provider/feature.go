package provider

import (
	"github.com/cross-team/clublink-backend/app/usecase/authorizer"
	"github.com/cross-team/clublink-backend/app/usecase/feature"
	"github.com/cross-team/clublink-backend/app/usecase/repository"
	"github.com/short-d/app/fw/env"
)

// NewFeatureDecisionMakerFactorySwitch creates FeatureDecisionFactory based on
// server environment.
func NewFeatureDecisionMakerFactorySwitch(
	deployment env.Deployment,
	toggleRepo repository.FeatureToggle,
	authorizer authorizer.Authorizer,
) feature.DecisionMakerFactory {
	if deployment.IsDevelopment() {
		return feature.NewStaticDecisionMakerFactory(authorizer)
	}
	return feature.NewDynamicDecisionMakerFactory(toggleRepo, authorizer)
}
