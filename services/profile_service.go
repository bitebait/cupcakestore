package services

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
	"strings"
)

type ProfileService interface {
	Create(profile *models.Profile) error
	FindByUserId(id uint) (models.Profile, error)
	Update(profile *models.Profile) error
}

type profileService struct {
	profileRepository repositories.ProfileRepository
}

func NewProfileService(profileRepository repositories.ProfileRepository) ProfileService {
	return &profileService{
		profileRepository: profileRepository,
	}
}

func (s *profileService) Create(profile *models.Profile) error {
	s.normalizeProfile(profile)
	return s.profileRepository.Create(profile)
}

func (s *profileService) FindByUserId(id uint) (models.Profile, error) {
	return s.profileRepository.FindByUserId(id)
}

func (s *profileService) Update(profile *models.Profile) error {
	s.normalizeProfile(profile)
	return s.profileRepository.Update(profile)
}

func (s *profileService) normalizeProfile(profile *models.Profile) {
	profile.FirstName = strings.Title(profile.FirstName)
	profile.LastName = strings.Title(profile.LastName)
	profile.State = strings.Title(profile.State)
	profile.City = strings.Title(profile.City)
	profile.Address = strings.Title(profile.Address)
}
