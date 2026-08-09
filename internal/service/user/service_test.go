package user_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	authdomain "github.com/vladgrskkh/onerep-auth/internal/domain/auth"
	"github.com/vladgrskkh/onerep-auth/internal/service/user"
	usermocks "github.com/vladgrskkh/onerep-auth/internal/service/user/mocks"
)

type ServiceTestSuite struct {
	suite.Suite

	svc  *user.UserService
	repo *usermocks.MockUserProfileRepository
}

func (s *ServiceTestSuite) SetupTest() {
	s.repo = usermocks.NewMockUserProfileRepository(s.T())
	s.svc = user.NewUserService(s.repo)
}

func (s *ServiceTestSuite) expectFound(found authdomain.User) {
	s.repo.EXPECT().FindByID(mock.Anything, found.ID).Return(found, nil)
}

func (s *ServiceTestSuite) expectUpdate() *authdomain.User {
	var stored authdomain.User
	s.repo.EXPECT().
		Update(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, u authdomain.User) (authdomain.User, error) {
			stored = u
			return u, nil
		})
	return &stored
}

func (s *ServiceTestSuite) TestUpdateProfile_AllFields() {
	userID := uuid.Must(uuid.NewV7())
	s.expectFound(authdomain.User{
		ID:          userID,
		Email:       "a@b.com",
		DisplayName: "Old Name",
		AvatarURL:   new("old.png"),
		Gender:      authdomain.GenderFemale,
	})
	stored := s.expectUpdate()

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
	s.Equal("New Name", stored.DisplayName)
	s.Equal(authdomain.GenderMale, stored.Gender)
	s.Equal("2000-01-02", stored.BirthDate.Format("2006-01-02"))
	s.Equal("new.png", *stored.AvatarURL)
	s.False(stored.UpdatedAt.IsZero())
	s.Equal("New Name", *profile.DisplayName)
}

func (s *ServiceTestSuite) TestUpdateProfile_InvalidBirthDate() {
	userID := uuid.Must(uuid.NewV7())
	s.expectFound(authdomain.User{ID: userID, DisplayName: "Alice"})

	birthDate := "01/02/2000"
	_, err := s.svc.UpdateProfile(context.Background(), user.UpdateProfileCommand{
		UserID:    userID,
		BirthDate: &birthDate,
	})
	s.Require().ErrorIs(err, authdomain.ErrInvalidBirthDate)
	s.repo.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything)
}

func (s *ServiceTestSuite) TestUpdateProfile_InvalidGender() {
	userID := uuid.Must(uuid.NewV7())
	s.expectFound(authdomain.User{ID: userID, DisplayName: "Alice"})

	gender := "attack"
	_, err := s.svc.UpdateProfile(context.Background(), user.UpdateProfileCommand{
		UserID: userID,
		Gender: &gender,
	})
	s.Require().ErrorIs(err, authdomain.ErrInvalidGender)
	s.repo.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything)
}

func (s *ServiceTestSuite) TestUpdateProfile_DisplayNameTooLong() {
	userID := uuid.Must(uuid.NewV7())
	s.expectFound(authdomain.User{ID: userID, DisplayName: "Alice"})

	name := strings.Repeat("a", authdomain.MaxDisplayNameLength+1)
	_, err := s.svc.UpdateProfile(context.Background(), user.UpdateProfileCommand{
		UserID:      userID,
		DisplayName: &name,
	})
	s.Require().ErrorIs(err, authdomain.ErrInvalidDisplayName)
	s.repo.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything)
}

func (s *ServiceTestSuite) TestUpdateProfile_EmptyDisplayName() {
	userID := uuid.Must(uuid.NewV7())
	s.expectFound(authdomain.User{ID: userID, DisplayName: "Alice"})

	name := "   "
	_, err := s.svc.UpdateProfile(context.Background(), user.UpdateProfileCommand{
		UserID:      userID,
		DisplayName: &name,
	})
	s.Require().ErrorIs(err, authdomain.ErrInvalidDisplayName)
	s.repo.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything)
}

func (s *ServiceTestSuite) TestUpdateProfile_PartialUpdateKeepsOtherFields() {
	userID := uuid.Must(uuid.NewV7())
	birthDate := time.Date(1990, 5, 1, 0, 0, 0, 0, time.UTC)
	avatar := "avatar.png"
	s.expectFound(authdomain.User{
		ID:          userID,
		Email:       "a@b.com",
		DisplayName: "Alice",
		AvatarURL:   &avatar,
		Gender:      authdomain.GenderFemale,
		BirthDate:   &birthDate,
	})
	stored := s.expectUpdate()

	gender := "other"
	_, err := s.svc.UpdateProfile(context.Background(), user.UpdateProfileCommand{
		UserID: userID,
		Gender: &gender,
	})
	s.Require().NoError(err)
	s.Equal("Alice", stored.DisplayName)
	s.Equal(authdomain.GenderOther, stored.Gender)
	s.Equal(&birthDate, stored.BirthDate)
	s.Equal(&avatar, stored.AvatarURL)
}

func (s *ServiceTestSuite) TestUpdateProfile_TrimsDisplayName() {
	userID := uuid.Must(uuid.NewV7())
	s.expectFound(authdomain.User{ID: userID, DisplayName: "Alice"})
	stored := s.expectUpdate()

	name := "  New Name  "
	_, err := s.svc.UpdateProfile(context.Background(), user.UpdateProfileCommand{
		UserID:      userID,
		DisplayName: &name,
	})
	s.Require().NoError(err)
	s.Equal("New Name", stored.DisplayName)
}

func TestServiceSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(ServiceTestSuite))
}
