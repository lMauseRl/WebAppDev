package usecase

import "github.com/lud0m4n/WebAppDev/internal/http/repository"

type UseCase struct {
	Repository *repository.Repository
}

func NewUseCase(repo *repository.Repository) *UseCase {
	return &UseCase{
		Repository: repo,
	}
}
