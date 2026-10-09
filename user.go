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
	res, err := s.client.Request(context.Background(), http.MethodGet, ProfilePath, nil)
	if err != nil {
		return nil, err
	}
	profile := res.GetProfileSettingsView()
	if profile == nil {
		return nil, fmt.Errorf("no profile in response")
	}
	return profile, nil
}

// Subscription: Get the subscription of the account tied to the device secret.
//
// A device without a subscription has the "basic" plan.
func (s *UserService) Subscription() (*Subscription, error) {
	res, err := s.client.Request(context.Background(), http.MethodGet, SubscriptionPath, nil)
	if err != nil {
		return nil, err
	}
	subscription := res.GetSubscriptionView().GetUserSubscription()
	if subscription == nil {
		return nil, fmt.Errorf("no subscription in response")
	}
	return subscription, nil
}
