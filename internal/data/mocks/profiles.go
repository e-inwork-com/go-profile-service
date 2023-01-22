package mocks

import (
	"time"

	"github.com/e-inwork-com/go-profile-service/internal/data"
	"github.com/google/uuid"
)

type ProfileModel struct{}

func (m ProfileModel) Insert(profile *data.Profile) error {
	profile.ID = MockFirstUUID()
	profile.CreatedAt = time.Now()
	profile.Version = 1

	return nil
}

func (m ProfileModel) GetByID(id uuid.UUID) (*data.Profile, error) {
	profileID := MockFirstUUID()

	if id == profileID {
		var profile = &data.Profile{
			ID:             profileID,
			CreatedAt:      time.Now(),
			ProfileUser:    MockFirstUUID(),
			ProfileName:    "John Doe",
			ProfilePicture: MockFirstUUID().String() + ".jpg",
			Version:        1,
		}
		return profile, nil
	}

	return nil, data.ErrRecordNotFound
}

func (m ProfileModel) GetByProfileUser(profileUser uuid.UUID) (*data.Profile, error) {
	profileUserID := MockFirstUUID()

	if profileUser == profileUserID {
		var profile = &data.Profile{
			ID:             MockFirstUUID(),
			CreatedAt:      time.Now(),
			ProfileUser:    profileUserID,
			ProfileName:    "John Doe",
			ProfilePicture: MockFirstUUID().String() + ".jpg",
			Version:        1,
		}
		return profile, nil
	}

	return nil, data.ErrRecordNotFound
}

func (m ProfileModel) Update(profile *data.Profile) error {
	profile.Version += 1

	return nil
}
