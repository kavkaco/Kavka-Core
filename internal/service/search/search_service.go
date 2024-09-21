package search

import (
	"context"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"github.com/kavkaco/Kavka-Core/log"
	"github.com/kavkaco/Kavka-Core/utils/vali"
)

type SearchService struct {
	searchRepository repository.SearchRepository
	validator        *vali.Vali
}

func NewSearchService(logger *log.SubLogger, searchRepository repository.SearchRepository) *SearchService {
	return &SearchService{searchRepository, vali.Validator()}
}

func (s *SearchService) Search(ctx context.Context, input string) (*model.SearchResultDTO, *vali.Varror) {
	varrors := s.validator.Validate(searchValidation{input})
	if len(varrors) > 0 {
		return nil, &vali.Varror{ValidationErrors: varrors}
	}

	result, err := s.searchRepository.Search(ctx, input)
	if err != nil {
		return nil, &vali.Varror{Error: err}
	}

	return result, nil
}
