package search

import (
	"context"

	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"github.com/kavkaco/Kavka-Core/log"
	"github.com/kavkaco/Kavka-Core/utils/vali"
)

const DefaultSearchLimit = 20

type SearchService struct {
	searchRepository repository.SearchRepository
	validator        *vali.Vali
	logger           *log.SubLogger
}

func NewSearchService(logger *log.SubLogger, searchRepository repository.SearchRepository) *SearchService {
	return &SearchService{searchRepository, vali.Validator(), logger}
}

func (s *SearchService) Search(ctx context.Context, input string, limit int) (*model.SearchResultDTO, *vali.ValiErr) {
	errs := s.validator.Validate(searchValidation{input})
	if len(errs) > 0 {
		return nil, &vali.ValiErr{ValidationErrors: errs}
	}

	if limit <= 0 || limit > DefaultSearchLimit {
		limit = DefaultSearchLimit
	}

	result, err := s.searchRepository.Search(ctx, input, limit)
	if err != nil {
		return nil, &vali.ValiErr{Error: err}
	}

	return result, nil
}
