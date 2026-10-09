package mangoplus

import (
	"context"
	"fmt"
	"net/http"
)

const (
	ProfilePath      = "profile"
	SubscriptionPath = "subscription"
)

type UserService service

// Profile: Get the profile of the account tied to the device secret.
func (s *UserService) Profile() (*ProfileSettingsView, error) {
	return s.ProfileContext(context.Background())
}

// ProfileContext: Same as Profile, with a context.
func (s *UserService) ProfileContext(ctx context.Context) (*ProfileSettingsView, error) {
	res, err := s.client.Request(ctx, http.MethodGet, ProfilePath, nil)
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}
	profile := res.GetProfileSettingsView()
	if profile == nil {
		return nil, fmt.Errorf("get profile: no profile in response")
	}
	return profile, nil
}

// Subscription: Get the subscription of the account tied to the device secret.
//
// A device without a subscription has the "basic" plan.
func (s *UserService) Subscription() (*Subscription, error) {
	return s.SubscriptionContext(context.Background())
}

// SubscriptionContext: Same as Subscription, with a context.
func (s *UserService) SubscriptionContext(ctx context.Context) (*Subscription, error) {
	res, err := s.client.Request(ctx, http.MethodGet, SubscriptionPath, nil)
	if err != nil {
		return nil, fmt.Errorf("get subscription: %w", err)
	}
	subscription := res.GetSubscriptionView().GetUserSubscription()
	if subscription == nil {
		return nil, fmt.Errorf("get subscription: no subscription in response")
	}
	return subscription, nil
}
