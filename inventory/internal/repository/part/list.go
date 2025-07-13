package part

import (
	"context"

	repoModel "github.com/space-wanderer/microservices/inventory/internal/repository/model"
)

func (r *repository) ListParts(_ context.Context, filter *repoModel.PartsFilter) ([]*repoModel.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if filter == nil || isEmptyFilter(filter) {
		parts := make([]*repoModel.Part, 0, len(r.parts))
		for _, part := range r.parts {
			parts = append(parts, part)
		}
		return parts, nil
	}

	var filteredParts []*repoModel.Part

	for _, part := range r.parts {
		if matchesFilter(part, filter) {
			filteredParts = append(filteredParts, part)
		}
	}

	return filteredParts, nil
}

// isEmptyFilter проверяет, пуст ли фильтр
func isEmptyFilter(filter *repoModel.PartsFilter) bool {
	return len(filter.Uuids) == 0 &&
		len(filter.Names) == 0 &&
		len(filter.Categories) == 0 &&
		len(filter.ManufacturerCountries) == 0 &&
		len(filter.Tags) == 0
}

// matchesFilter проверяет, соответствует ли деталь всем условиям фильтра
func matchesFilter(part *repoModel.Part, filter *repoModel.PartsFilter) bool {
	return matchesUuidFilter(part, filter.Uuids) &&
		matchesNameFilter(part, filter.Names) &&
		matchesCategoryFilter(part, filter.Categories) &&
		matchesCountryFilter(part, filter.ManufacturerCountries) &&
		matchesTagFilter(part, filter.Tags)
}

// matchesUuidFilter проверяет соответствие UUID
func matchesUuidFilter(part *repoModel.Part, uuids []string) bool {
	if len(uuids) == 0 {
		return true
	}
	for _, uuid := range uuids {
		if part.UUID == uuid {
			return true
		}
	}
	return false
}

// matchesNameFilter проверяет соответствие имени
func matchesNameFilter(part *repoModel.Part, names []string) bool {
	if len(names) == 0 {
		return true
	}
	for _, name := range names {
		if part.Name == name {
			return true
		}
	}
	return false
}

// matchesCategoryFilter проверяет соответствие категории
func matchesCategoryFilter(part *repoModel.Part, categories []repoModel.Category) bool {
	if len(categories) == 0 {
		return true
	}
	for _, category := range categories {
		if part.Category == category {
			return true
		}
	}
	return false
}

// matchesCountryFilter проверяет соответствие страны производителя
func matchesCountryFilter(part *repoModel.Part, countries []string) bool {
	if len(countries) == 0 {
		return true
	}
	if part.Manufacturer == nil {
		return false
	}
	for _, country := range countries {
		if part.Manufacturer.Country == country {
			return true
		}
	}
	return false
}

// matchesTagFilter проверяет соответствие тегов
func matchesTagFilter(part *repoModel.Part, tags []string) bool {
	if len(tags) == 0 {
		return true
	}
	for _, filterTag := range tags {
		for _, partTag := range part.Tags {
			if partTag == filterTag {
				return true
			}
		}
	}
	return false
}
