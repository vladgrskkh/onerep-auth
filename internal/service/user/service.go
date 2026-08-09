package user

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

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

const birthDateFormat = "2006-01-02"

// UpdateProfile applies the non-nil fields of the command to the stored user.
// Field values are validated here, not at the handler boundary.
func (s *UserService) UpdateProfile(ctx context.Context, cmd UpdateProfileCommand) (userdomain.UserProfile, error) {
	u, err := s.users.FindByID(ctx, cmd.UserID)
	if err != nil {
		return userdomain.UserProfile{}, err
	}

	if cmd.DisplayName != nil {
		name := strings.TrimSpace(*cmd.DisplayName)
		if name == "" || utf8.RuneCountInString(name) > authdomain.MaxDisplayNameLength {
			return userdomain.UserProfile{}, authdomain.ErrInvalidDisplayName
		}
		u.DisplayName = name
	}
	if cmd.Gender != nil {
		gender := authdomain.Gender(*cmd.Gender)
		if !gender.IsValid() {
			return userdomain.UserProfile{}, authdomain.ErrInvalidGender
		}
		u.Gender = gender
	}
	if cmd.BirthDate != nil {
		birthDate, parseErr := time.Parse(birthDateFormat, *cmd.BirthDate)
		if parseErr != nil {
			return userdomain.UserProfile{}, authdomain.ErrInvalidBirthDate
		}
		u.BirthDate = &birthDate
	}
	if cmd.AvatarURL != nil {
		u.AvatarURL = cmd.AvatarURL
	}

	u.UpdatedAt = time.Now()
	updated, err := s.users.Update(ctx, u)
	if err != nil {
		return userdomain.UserProfile{}, err
	}

	return userdomain.ToProfile(updated, cmd.UserID), nil
}
