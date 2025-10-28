package prana

// ActionCreator is a convenience type for creating actions.
type ActionCreator struct {
	actionType string
}

// NewActionCreator creates a new action creator.
func NewActionCreator(actionType string) *ActionCreator {
	return &ActionCreator{
		actionType: actionType,
	}
}

// Create creates an action with the given payload.
func (ac *ActionCreator) Create(payload any) Action {
	return Action{
		Type:    ac.actionType,
		Payload: payload,
	}
}

// Type returns the action type.
func (ac *ActionCreator) Type() string {
	return ac.actionType
}
