package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/vladgrskkh/onerep-auth/internal/domain"
)

type userProfileFinder interface {
	FindByID(ctx context.Context, id uuid.UUID) (domain.User, error)
}

type userProfileUpdater interface {
	Update(ctx context.Context, user domain.User) (domain.User, error)
}

type UserService struct {
	finder  userProfileFinder
	updater userProfileUpdater
}

func NewUserService(finder userProfileFinder, updater userProfileUpdater) *UserService {
	return &UserService{finder: finder, updater: updater}
}

func (s *UserService) GetProfile(ctx context.Context, userID, requesterID uuid.UUID) (domain.UserProfile, error) {
	user, err := s.finder.FindByID(ctx, userID)
	if err != nil {
		return domain.UserProfile{}, err
	}
	return user.ToProfile(requesterID), nil
}

type UpdateProfileInput struct {
	DisplayName *string        `json:"display_name,omitempty"`
	Gender      *domain.Gender `json:"gender,omitempty"`
	BirthDate   *time.Time     `json:"birth_date,omitempty"`
	AvatarURL   *string        `json:"avatar_url,omitempty"`
}

func (s *UserService) UpdateProfile(
	ctx context.Context,
	userID uuid.UUID,
	input UpdateProfileInput,
) (domain.UserProfile, error) {
	user, err := s.finder.FindByID(ctx, userID)
	if err != nil {
		return domain.UserProfile{}, err
	}

	if input.DisplayName != nil {
		user.DisplayName = *input.DisplayName
	}
	if input.Gender != nil {
		user.Gender = *input.Gender
	}
	if input.BirthDate != nil {
		user.BirthDate = input.BirthDate
	}
	if input.AvatarURL != nil {
		user.AvatarURL = input.AvatarURL
	}

	user.UpdatedAt = time.Now()
	updated, err := s.updater.Update(ctx, user)
	if err != nil {
		return domain.UserProfile{}, err
	}

	return updated.ToProfile(userID), nil
}
