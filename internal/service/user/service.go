package user

import (
	"context"
	"time"

	"github.com/google/uuid"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	userdomain "github.com/vladgrskkh/onerep-auth/internal/domain/user"
)

type UserProfileFinder interface {
	FindByID(ctx context.Context, id uuid.UUID) (authdomain.User, error)
}

type UserProfileUpdater interface {
	Update(ctx context.Context, user authdomain.User) (authdomain.User, error)
}

type UserService struct {
	finder  UserProfileFinder
	updater UserProfileUpdater
}

func NewUserService(finder UserProfileFinder, updater UserProfileUpdater) *UserService {
	return &UserService{finder: finder, updater: updater}
}

func (s *UserService) GetProfile(ctx context.Context, userID, requesterID uuid.UUID) (authdomain.UserProfile, error) {
	u, err := s.finder.FindByID(ctx, userID)
	if err != nil {
		return authdomain.UserProfile{}, err
	}
	return u.ToProfile(requesterID), nil
}

func (s *UserService) UpdateProfile(
	ctx context.Context,
	userID uuid.UUID,
	input userdomain.UpdateProfileInput,
) (authdomain.UserProfile, error) {
	u, err := s.finder.FindByID(ctx, userID)
	if err != nil {
		return authdomain.UserProfile{}, err
	}

	if input.DisplayName != nil {
		u.DisplayName = *input.DisplayName
	}
	if input.Gender != nil {
		u.Gender = *input.Gender
	}
	if input.BirthDate != nil {
		u.BirthDate = input.BirthDate
	}
	if input.AvatarURL != nil {
		u.AvatarURL = input.AvatarURL
	}

	u.UpdatedAt = time.Now()
	updated, err := s.updater.Update(ctx, u)
	if err != nil {
		return authdomain.UserProfile{}, err
	}

	return updated.ToProfile(userID), nil
}
