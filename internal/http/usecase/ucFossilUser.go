package usecase

import (
	"errors"
	"strings"

	"github.com/lud0m4n/WebAppDev/internal/model"
)

func (uc *UseCase) GetFossilForUser(searchSpecies, startFormationDate, endFormationDate, fossilStatus string, userID uint) ([]model.FossilRequest, error) {
	searchSpecies = strings.ToUpper(searchSpecies + "%")
	fossilStatus = strings.ToLower(fossilStatus + "%")

	if userID <= 0 {
		return nil, errors.New("недопустимый ИД пользователя")
	}

	fossil, err := uc.Repository.GetFossilForUser(searchSpecies, startFormationDate, endFormationDate, fossilStatus, userID)
	if err != nil {
		return nil, err
	}

	return fossil, nil
}

func (uc *UseCase) GetFossilByIDForUser(fossilID int, userID uint) (model.FossilGetResponse, error) {
	if fossilID <= 0 {
		return model.FossilGetResponse{}, errors.New("недопустимый ИД доставки")
	}
	if userID <= 0 {
		return model.FossilGetResponse{}, errors.New("недопустимый ИД пользователя")
	}

	fossil, err := uc.Repository.GetFossilByIDForUser(fossilID, userID)
	if err != nil {
		return model.FossilGetResponse{}, err
	}

	return fossil, nil
}

func (uc *UseCase) DeleteFossilForUser(fossilID int, userID uint) error {
	if fossilID <= 0 {
		return errors.New("недопустимый ИД доставки")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}

	err := uc.Repository.DeleteFossilForUser(fossilID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) UpdateFossilForUser(fossilID int, userID uint, species model.FossilUpdateSpeciesRequest) error {
	if fossilID <= 0 {
		return errors.New("недопустимый ИД доставки")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}
	if len(species.Species) != 6 {
		return errors.New("недопустимый номер рейса")
	}
	//////////////
	err := uc.Repository.UpdateFossilForUser(fossilID, userID, &species)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) UpdateFossilStatusForUser(fossilID int, userID uint) error {
	if fossilID <= 0 {
		return errors.New("недопустимый ИД доставки")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}

	err := uc.Repository.UpdateFossilStatusForUser(fossilID, userID)
	if err != nil {
		return err
	}

	return nil
}
