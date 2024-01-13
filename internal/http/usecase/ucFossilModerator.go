package usecase

import (
	"errors"
	"strings"

	"github.com/lud0m4n/WebAppDev/internal/model"
)

func (uc *UseCase) GetFossilForModerator(searchSpecies, startFormationDate, endFormationDate, fossilStatus string, moderatorID uint) ([]model.FossilRequest, error) {
	searchSpecies = strings.ToLower(searchSpecies + "%")
	fossilStatus = strings.ToLower(fossilStatus + "%")

	if moderatorID <= 0 {
		return nil, errors.New("недопустимый ИД модератора")
	}

	deliveries, err := uc.Repository.GetFossilForModerator(searchSpecies, startFormationDate, endFormationDate, fossilStatus, moderatorID)
	if err != nil {
		return nil, err
	}

	return deliveries, nil
}

func (uc *UseCase) GetFossilByIDForModerator(fossilID int, moderatorID uint) (model.FossilGetResponse, error) {
	if fossilID <= 0 {
		return model.FossilGetResponse{}, errors.New("недопустимый ИД останка")
	}
	if moderatorID <= 0 {
		return model.FossilGetResponse{}, errors.New("недопустимый ИД модератора")
	}

	deliveries, err := uc.Repository.GetFossilByIDForModerator(fossilID, moderatorID)
	if err != nil {
		return model.FossilGetResponse{}, err
	}

	return deliveries, nil
}

func (uc *UseCase) UpdateFossilForModerator(fossilID int, moderatorID uint, species model.FossilUpdateSpeciesRequest) error {
	if fossilID <= 0 {
		return errors.New("недопустимый ИД останка")
	}
	if moderatorID <= 0 {
		return errors.New("недопустимый ИД модератора")
	}

	err := uc.Repository.UpdateFossilForModerator(fossilID, moderatorID, &species)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) UpdateFossilStatusForModerator(fossilID int, moderatorID uint, fossilStatus model.FossilUpdateStatusRequest) error {
	if fossilID <= 0 {
		return errors.New("недопустимый ИД останка")
	}
	if moderatorID <= 0 {
		return errors.New("недопустимый ИД модератора")
	}
	if fossilStatus.Status != model.FOSSIL_STATUS_COMPLETED && fossilStatus.Status != model.FOSSIL_STATUS_REJECTED {
		return errors.New("текущий статус останка уже завершен или отклонен")
	}

	err := uc.Repository.UpdateFossilStatusForModerator(fossilID, moderatorID, &fossilStatus)
	if err != nil {
		return err
	}

	return nil
}
