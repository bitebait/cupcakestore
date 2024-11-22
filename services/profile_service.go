package services

import (
	"errors"
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

	if err := s.profileRepository.Create(profile); err != nil {
		return errors.New("falha ao criar o perfil do usuário, verifique os dados e tente novamente")
	}

	return nil
}

func (s *profileService) FindByUserId(id uint) (models.Profile, error) {
	res, err := s.profileRepository.FindByUserId(id)

	if err != nil {
		err = errors.New("falha ao encontrar o perfil do usuário, tente novamente")
	}

	return res, err
}

func (s *profileService) Update(profile *models.Profile) error {
	s.normalizeProfile(profile)

	if err := s.profileRepository.Update(profile); err != nil {
		return errors.New("falha ao atualizar o perfil do usuário, verifique os dados e tente novamente")
	}

	return nil
}

func (s *profileService) normalizeProfile(profile *models.Profile) {
	profile.FirstName = normalizeString(profile.FirstName)
	profile.LastName = normalizeString(profile.LastName)
	profile.State = normalizeString(profile.State)
	profile.City = normalizeString(profile.City)
	profile.Address = normalizeString(profile.Address)
}

func normalizeString(s string) string {
	return strings.Title(strings.TrimSpace(s))
}
