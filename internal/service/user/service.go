package user

import (
	"context"
	"time"

	"github.com/google/uuid"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	userdomain "github.com/vladgrskkh/onerep-auth/internal/domain/user"
)

type UserProfileRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (authdomain.User, error)
	Update(ctx context.Context, user authdomain.User) (authdomain.User, error)
}

type UserService struct {
	users UserProfileRepository
}

func NewUserService(users UserProfileRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) GetProfile(ctx context.Context, userID, requesterID uuid.UUID) (userdomain.UserProfile, error) {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return userdomain.UserProfile{}, err
	}
	return userdomain.ToProfile(u, requesterID), nil
}

// UpdateProfile applies the non-nil fields of the profile to the stored user.
func (s *UserService) UpdateProfile(
	ctx context.Context,
	profile userdomain.UserProfile,
) (userdomain.UserProfile, error) {
	u, err := s.users.FindByID(ctx, profile.ID)
	if err != nil {
		return userdomain.UserProfile{}, err
	}

	if profile.DisplayName != nil {
		u.DisplayName = *profile.DisplayName
	}
	if profile.Gender != nil {
		u.Gender = *profile.Gender
	}
	if profile.BirthDate != nil {
		u.BirthDate = profile.BirthDate
	}
	if profile.AvatarURL != nil {
		u.AvatarURL = profile.AvatarURL
	}

	u.UpdatedAt = time.Now()
	updated, err := s.users.Update(ctx, u)
	if err != nil {
		return userdomain.UserProfile{}, err
	}

	return userdomain.ToProfile(updated, profile.ID), nil
}
