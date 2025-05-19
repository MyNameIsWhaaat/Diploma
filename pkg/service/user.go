package service

import "github.com/MyNameIsWhaaat/algo-learning/pkg/domain"

func (s *Service) GetUserProfileData(userID int) (domain.UserProfileData, error) {
	return s.repo.GetUserProfileData(userID)
}

func (s *Service) MarkUserNotNew(userID int) error {
	return s.repo.UpdateUserIsNewFlag(userID, false)
}
