package services

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
)

// CRUDService provides generic CRUD operations for any model
type CRUDService[T any] struct {
	db        *gorm.DB
	modelType reflect.Type
}

// NewCRUDService creates a new generic CRUD service for the specified model type
func NewCRUDService[T any](db *gorm.DB) *CRUDService[T] {
	var model T
	return &CRUDService[T]{
		db:        db,
		modelType: reflect.TypeOf(model),
	}
}

// PaginationOptions defines pagination parameters
type PaginationOptions struct {
	Page     int `json:"page" form:"page" validate:"min=1"`
	PageSize int `json:"page_size" form:"page_size" validate:"min=1,max=100"`
}

// SortOptions defines sorting parameters
type SortOptions struct {
	Field     string `json:"field" form:"field"`
	Direction string `json:"direction" form:"direction"` // "asc" or "desc"
}

// FilterOptions defines filtering parameters
type FilterOptions struct {
	Filters map[string]interface{} `json:"filters" form:"filters"`
}

// QueryOptions combines all query options
type QueryOptions struct {
	Pagination PaginationOptions  `json:"pagination"`
	Sort       []SortOptions      `json:"sort"`
	Filter     FilterOptions      `json:"filter"`
	Search     string             `json:"search" form:"search"`
	Preload    []string           `json:"preload" form:"preload"`
}

// PaginatedResult represents a paginated result set
type PaginatedResult[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// Create creates a new record
func (s *CRUDService[T]) Create(model *T) error {
	return s.db.Create(model).Error
}

// CreateBatch creates multiple records in a single transaction
func (s *CRUDService[T]) CreateBatch(models []T) error {
	return s.db.CreateInBatches(models, 100).Error
}

// GetByID retrieves a record by its ID
func (s *CRUDService[T]) GetByID(id interface{}, preload ...string) (*T, error) {
	var model T
	query := s.db

	// Add preloads if specified
	for _, rel := range preload {
		query = query.Preload(rel)
	}

	err := query.First(&model, id).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

// GetAll retrieves all records with optional query options
func (s *CRUDService[T]) GetAll(options QueryOptions) (*PaginatedResult[T], error) {
	var models []T
	var total int64

	query := s.db.Model(new(T))

	// Apply filters
	query = s.applyFilters(query, options.Filter)

	// Apply search
	if options.Search != "" {
		query = s.applySearch(query, options.Search)
	}

	// Count total records (before pagination)
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	query = s.applySorting(query, options.Sort)

	// Apply preloads
	for _, rel := range options.Preload {
		query = query.Preload(rel)
	}

	// Apply pagination
	offset := (options.Pagination.Page - 1) * options.Pagination.PageSize
	query = query.Offset(offset).Limit(options.Pagination.PageSize)

	// Execute query
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	// Calculate pagination info
	totalPages := int((total + int64(options.Pagination.PageSize) - 1) / int64(options.Pagination.PageSize))
	hasNext := options.Pagination.Page < totalPages
	hasPrev := options.Pagination.Page > 1

	return &PaginatedResult[T]{
		Data:       models,
		Total:      total,
		Page:       options.Pagination.Page,
		PageSize:   options.Pagination.PageSize,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}, nil
}

// Update updates a record by ID
func (s *CRUDService[T]) Update(id interface{}, updates map[string]interface{}) error {
	return s.db.Model(new(T)).Where("id = ?", id).Updates(updates).Error
}

// UpdateStruct updates a record using a struct
func (s *CRUDService[T]) UpdateStruct(id interface{}, model *T) error {
	return s.db.Model(model).Where("id = ?", id).Updates(model).Error
}

// Delete soft deletes a record by ID
func (s *CRUDService[T]) Delete(id interface{}) error {
	return s.db.Delete(new(T), id).Error
}

// HardDelete permanently deletes a record by ID
func (s *CRUDService[T]) HardDelete(id interface{}) error {
	return s.db.Unscoped().Delete(new(T), id).Error
}

// DeleteBatch soft deletes multiple records by IDs
func (s *CRUDService[T]) DeleteBatch(ids []interface{}) error {
	return s.db.Delete(new(T), ids).Error
}

// Restore restores a soft-deleted record by ID
func (s *CRUDService[T]) Restore(id interface{}) error {
	return s.db.Model(new(T)).Unscoped().Where("id = ?", id).Update("deleted_at", nil).Error
}

// Exists checks if a record exists with the given conditions
func (s *CRUDService[T]) Exists(conditions map[string]interface{}) (bool, error) {
	var count int64
	query := s.db.Model(new(T))
	
	for field, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}
	
	err := query.Count(&count).Error
	return count > 0, err
}

// Count returns the total number of records matching the conditions
func (s *CRUDService[T]) Count(conditions map[string]interface{}) (int64, error) {
	var count int64
	query := s.db.Model(new(T))
	
	for field, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}
	
	err := query.Count(&count).Error
	return count, err
}

// FindOne finds a single record matching the conditions
func (s *CRUDService[T]) FindOne(conditions map[string]interface{}, preload ...string) (*T, error) {
	var model T
	query := s.db

	// Add preloads if specified
	for _, rel := range preload {
		query = query.Preload(rel)
	}

	for field, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	err := query.First(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

// FindMany finds multiple records matching the conditions
func (s *CRUDService[T]) FindMany(conditions map[string]interface{}, options QueryOptions) (*PaginatedResult[T], error) {
	// Add conditions to filters
	if options.Filter.Filters == nil {
		options.Filter.Filters = make(map[string]interface{})
	}
	
	for field, value := range conditions {
		options.Filter.Filters[field] = value
	}

	return s.GetAll(options)
}

// BulkUpdate updates multiple records matching the conditions
func (s *CRUDService[T]) BulkUpdate(conditions map[string]interface{}, updates map[string]interface{}) error {
	query := s.db.Model(new(T))
	
	for field, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}
	
	return query.Updates(updates).Error
}

// Transaction executes a function within a database transaction
func (s *CRUDService[T]) Transaction(fn func(*gorm.DB) error) error {
	return s.db.Transaction(fn)
}

// applyFilters applies filtering conditions to the query
func (s *CRUDService[T]) applyFilters(query *gorm.DB, filter FilterOptions) *gorm.DB {
	for field, value := range filter.Filters {
		switch v := value.(type) {
		case string:
			if strings.Contains(field, "_like") {
				actualField := strings.Replace(field, "_like", "", 1)
				query = query.Where(fmt.Sprintf("%s LIKE ?", actualField), "%"+v+"%")
			} else if strings.Contains(field, "_in") {
				actualField := strings.Replace(field, "_in", "", 1)
				query = query.Where(fmt.Sprintf("%s IN (?)", actualField), v)
			} else {
				query = query.Where(fmt.Sprintf("%s = ?", field), v)
			}
		case []interface{}:
			query = query.Where(fmt.Sprintf("%s IN (?)", field), v)
		case map[string]interface{}:
			// Handle range queries
			if from, ok := v["from"]; ok {
				query = query.Where(fmt.Sprintf("%s >= ?", field), from)
			}
			if to, ok := v["to"]; ok {
				query = query.Where(fmt.Sprintf("%s <= ?", field), to)
			}
		default:
			query = query.Where(fmt.Sprintf("%s = ?", field), value)
		}
	}
	return query
}

// applySorting applies sorting to the query
func (s *CRUDService[T]) applySorting(query *gorm.DB, sorts []SortOptions) *gorm.DB {
	for _, sort := range sorts {
		if sort.Field != "" {
			direction := "ASC"
			if strings.ToLower(sort.Direction) == "desc" {
				direction = "DESC"
			}
			query = query.Order(fmt.Sprintf("%s %s", sort.Field, direction))
		}
	}
	
	// Default sorting if no sort specified
	if len(sorts) == 0 {
		query = query.Order("id DESC")
	}
	
	return query
}

// applySearch applies search functionality to searchable fields
func (s *CRUDService[T]) applySearch(query *gorm.DB, search string) *gorm.DB {
	// Get searchable fields based on the model type
	searchableFields := s.getSearchableFields()
	
	if len(searchableFields) > 0 {
		searchQuery := ""
		args := []interface{}{}
		
		for i, field := range searchableFields {
			if i > 0 {
				searchQuery += " OR "
			}
			searchQuery += fmt.Sprintf("%s LIKE ?", field)
			args = append(args, "%"+search+"%")
		}
		
		query = query.Where(searchQuery, args...)
	}
	
	return query
}

// getSearchableFields returns fields that should be searchable
func (s *CRUDService[T]) getSearchableFields() []string {
	// Default searchable fields based on common naming patterns
	searchableFields := []string{}
	
	// Use reflection to find string fields that might be searchable
	modelType := s.modelType
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldType := field.Type
		
		// Check if it's a string field and likely searchable
		if fieldType.Kind() == reflect.String {
			fieldName := field.Name
			gormTag := field.Tag.Get("gorm")
			
			// Get the database column name from gorm tag or use field name
			dbFieldName := strings.ToLower(fieldName)
			if strings.Contains(gormTag, "column:") {
				parts := strings.Split(gormTag, "column:")
				if len(parts) > 1 {
					colName := strings.Split(parts[1], ";")[0]
					dbFieldName = colName
				}
			} else {
				// Convert camelCase to snake_case
				dbFieldName = camelToSnake(fieldName)
			}
			
			// Add common searchable fields
			if strings.Contains(strings.ToLower(fieldName), "name") ||
			   strings.Contains(strings.ToLower(fieldName), "title") ||
			   strings.Contains(strings.ToLower(fieldName), "description") ||
			   strings.Contains(strings.ToLower(fieldName), "email") ||
			   strings.Contains(strings.ToLower(fieldName), "username") {
				searchableFields = append(searchableFields, dbFieldName)
			}
		}
	}
	
	return searchableFields
}

// DefaultQueryOptions returns default query options with sensible defaults
func DefaultQueryOptions() QueryOptions {
	return QueryOptions{
		Pagination: PaginationOptions{
			Page:     1,
			PageSize: 20,
		},
		Sort:    []SortOptions{},
		Filter:  FilterOptions{Filters: make(map[string]interface{})},
		Search:  "",
		Preload: []string{},
	}
}

// ValidateQueryOptions validates and applies defaults to query options
func ValidateQueryOptions(options *QueryOptions) {
	if options.Pagination.Page <= 0 {
		options.Pagination.Page = 1
	}
	if options.Pagination.PageSize <= 0 {
		options.Pagination.PageSize = 20
	}
	if options.Pagination.PageSize > 100 {
		options.Pagination.PageSize = 100
	}
	if options.Filter.Filters == nil {
		options.Filter.Filters = make(map[string]interface{})
	}
}

// camelToSnake converts camelCase to snake_case
func camelToSnake(s string) string {
	result := ""
	for i, char := range s {
		if i > 0 && char >= 'A' && char <= 'Z' {
			result += "_"
		}
		result += strings.ToLower(string(char))
	}
	return result
}