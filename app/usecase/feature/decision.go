package feature

import (
	"github.com/cross-team/clublink-backend/app/entity"
	"github.com/cross-team/clublink-backend/app/usecase/instrumentation"
)

// DecisionMaker determines whether a feature should be turned on or off under
// certain conditions.
type DecisionMaker interface {
	IsFeatureEnable(featureID string, user *entity.User) bool
}

// DecisionMakerFactory creates feature decision maker.
type DecisionMakerFactory interface {
	NewDecision(instrumentation instrumentation.Instrumentation) DecisionMaker
}
