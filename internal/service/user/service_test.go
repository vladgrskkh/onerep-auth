package user_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	"github.com/vladgrskkh/onerep-auth/internal/service/user"
)

type fakeUserProfileRepository struct {
	found     authdomain.User
	updated   authdomain.User
	updateErr error
}

func (f *fakeUserProfileRepository) FindByID(_ context.Context, _ uuid.UUID) (authdomain.User, error) {
	return f.found, nil
}

func (f *fakeUserProfileRepository) Update(_ context.Context, u authdomain.User) (authdomain.User, error) {
	f.updated = u
	return f.updated, f.updateErr
}

type ServiceTestSuite struct {
	suite.Suite

	svc  *user.UserService
	repo *fakeUserProfileRepository
}

func (s *ServiceTestSuite) SetupTest() {
	s.repo = &fakeUserProfileRepository{}
	s.svc = user.NewUserService(s.repo)
}

func (s *ServiceTestSuite) TestUpdateProfile_AllFields() {
	userID := uuid.Must(uuid.NewV7())
	s.repo.found = authdomain.User{
		ID:          userID,
		Email:       "a@b.com",
		DisplayName: "Old Name",
		AvatarURL:   new("old.png"),
		Gender:      authdomain.GenderFemale,
	}

	birthDate := "2000-01-02"
	avatar := "new.png"
	cmd := user.UpdateProfileCommand{
		UserID:      userID,
		DisplayName: new("New Name"),
		Gender:      new("male"),
		BirthDate:   &birthDate,
		AvatarURL:   &avatar,
	}

	profile, err := s.svc.UpdateProfile(context.Background(), cmd)
	s.Require().NoError(err)
	s.Equal("New Name", s.repo.updated.DisplayName)
	s.Equal(authdomain.GenderMale, s.repo.updated.Gender)
	s.Equal("2000-01-02", s.repo.updated.BirthDate.Format("2006-01-02"))
	s.Equal("new.png", *s.repo.updated.AvatarURL)
	s.False(s.repo.updated.UpdatedAt.IsZero())
	s.Equal("New Name", *profile.DisplayName)
}

func (s *ServiceTestSuite) TestUpdateProfile_InvalidBirthDate() {
	userID := uuid.Must(uuid.NewV7())
	s.repo.found = authdomain.User{ID: userID, DisplayName: "Alice"}

	birthDate := "01/02/2000"
	_, err := s.svc.UpdateProfile(context.Background(), user.UpdateProfileCommand{
		UserID:    userID,
		BirthDate: &birthDate,
	})
	s.Require().ErrorIs(err, authdomain.ErrInvalidBirthDate)
	s.Equal(authdomain.User{}, s.repo.updated)
}

func (s *ServiceTestSuite) TestUpdateProfile_InvalidGender() {
	userID := uuid.Must(uuid.NewV7())
	s.repo.found = authdomain.User{ID: userID, DisplayName: "Alice"}

	gender := "attack"
	_, err := s.svc.UpdateProfile(context.Background(), user.UpdateProfileCommand{
		UserID: userID,
		Gender: &gender,
	})
	s.Require().ErrorIs(err, authdomain.ErrInvalidGender)
	s.Equal(authdomain.User{}, s.repo.updated)
}

func (s *ServiceTestSuite) TestUpdateProfile_DisplayNameTooLong() {
	userID := uuid.Must(uuid.NewV7())
	s.repo.found = authdomain.User{ID: userID, DisplayName: "Alice"}

	name := strings.Repeat("a", authdomain.MaxDisplayNameLength+1)
	_, err := s.svc.UpdateProfile(context.Background(), user.UpdateProfileCommand{
		UserID:      userID,
		DisplayName: &name,
	})
	s.Require().ErrorIs(err, authdomain.ErrInvalidDisplayName)
	s.Equal(authdomain.User{}, s.repo.updated)
}

func (s *ServiceTestSuite) TestUpdateProfile_EmptyDisplayName() {
	userID := uuid.Must(uuid.NewV7())
	s.repo.found = authdomain.User{ID: userID, DisplayName: "Alice"}

	name := "   "
	_, err := s.svc.UpdateProfile(context.Background(), user.UpdateProfileCommand{
		UserID:      userID,
		DisplayName: &name,
	})
	s.Require().ErrorIs(err, authdomain.ErrInvalidDisplayName)
	s.Equal(authdomain.User{}, s.repo.updated)
}

func (s *ServiceTestSuite) TestUpdateProfile_PartialUpdateKeepsOtherFields() {
	userID := uuid.Must(uuid.NewV7())
	birthDate := time.Date(1990, 5, 1, 0, 0, 0, 0, time.UTC)
	avatar := "avatar.png"
	s.repo.found = authdomain.User{
		ID:          userID,
		Email:       "a@b.com",
		DisplayName: "Alice",
		AvatarURL:   &avatar,
		Gender:      authdomain.GenderFemale,
		BirthDate:   &birthDate,
	}

	gender := "other"
	_, err := s.svc.UpdateProfile(context.Background(), user.UpdateProfileCommand{
		UserID: userID,
		Gender: &gender,
	})
	s.Require().NoError(err)
	s.Equal("Alice", s.repo.updated.DisplayName)
	s.Equal(authdomain.GenderOther, s.repo.updated.Gender)
	s.Equal(&birthDate, s.repo.updated.BirthDate)
	s.Equal(&avatar, s.repo.updated.AvatarURL)
}

func (s *ServiceTestSuite) TestUpdateProfile_TrimsDisplayName() {
	userID := uuid.Must(uuid.NewV7())
	s.repo.found = authdomain.User{ID: userID, DisplayName: "Alice"}

	name := "  New Name  "
	_, err := s.svc.UpdateProfile(context.Background(), user.UpdateProfileCommand{
		UserID:      userID,
		DisplayName: &name,
	})
	s.Require().NoError(err)
	s.Equal("New Name", s.repo.updated.DisplayName)
}

func TestServiceSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(ServiceTestSuite))
}
