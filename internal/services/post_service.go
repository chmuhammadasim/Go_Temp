package services

import (
	"errors"
	"go-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

// PostService provides post-specific business logic using the generic CRUD service
type PostService struct {
	*CRUDService[models.Post]
	db           *gorm.DB
	auditService *AuditService
}

// NewPostService creates a new post service instance
func NewPostService(db *gorm.DB, auditService *AuditService) *PostService {
	return &PostService{
		CRUDService:  NewCRUDService[models.Post](db),
		db:           db,
		auditService: auditService,
	}
}

// CreatePost creates a new post with audit logging
func (s *PostService) CreatePost(userID uint, title, content string) (*models.Post, error) {
	post := &models.Post{
		UserID:    userID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.Create(post); err != nil {
		return nil, err
	}

	// Log the creation in audit trail
	if s.auditService != nil {
		auditData := AuditEventData{
			EntityType: "post",
			EntityID:   string(rune(post.ID)),
			NewValues: map[string]interface{}{
				"title":   post.Title,
				"content": post.Content,
				"user_id": post.UserID,
			},
		}
		s.auditService.LogEvent(userID, ActionCreate, auditData)
	}

	return post, nil
}

// UpdatePost updates a post with authorization check
func (s *PostService) UpdatePost(postID, userID uint, title, content *string) (*models.Post, error) {
	// Get the existing post
	existingPost, err := s.GetByID(postID, "User")
	if err != nil {
		return nil, err
	}

	// Check if user owns the post or is admin
	if existingPost.UserID != userID {
		var user models.User
		if err := s.db.First(&user, userID).Error; err != nil {
			return nil, errors.New("unauthorized")
		}
		if user.Role != models.RoleAdmin && user.Role != models.RoleModerator {
			return nil, errors.New("unauthorized to edit this post")
		}
	}

	// Store old values for audit
	oldValues := map[string]interface{}{
		"title":   existingPost.Title,
		"content": existingPost.Content,
	}

	// Prepare updates
	updates := make(map[string]interface{})
	newValues := make(map[string]interface{})

	if title != nil {
		updates["title"] = *title
		newValues["title"] = *title
	}

	if content != nil {
		updates["content"] = *content
		newValues["content"] = *content
	}

	// Add update timestamp
	updates["updated_at"] = time.Now()

	// Perform the update
	if len(updates) > 1 { // More than just updated_at
		if err := s.Update(postID, updates); err != nil {
			return nil, err
		}

		// Log the update in audit trail
		if s.auditService != nil {
			auditData := AuditEventData{
				EntityType: "post",
				EntityID:   string(rune(postID)),
				OldValues:  oldValues,
				NewValues:  newValues,
			}
			s.auditService.LogEvent(userID, ActionUpdate, auditData)
		}
	}

	// Get and return the updated post
	return s.GetByID(postID, "User")
}

// DeletePost deletes a post with authorization check
func (s *PostService) DeletePost(postID, userID uint) error {
	// Get the existing post
	existingPost, err := s.GetByID(postID)
	if err != nil {
		return err
	}

	// Check if user owns the post or is admin
	if existingPost.UserID != userID {
		var user models.User
		if err := s.db.First(&user, userID).Error; err != nil {
			return errors.New("unauthorized")
		}
		if user.Role != models.RoleAdmin && user.Role != models.RoleModerator {
			return errors.New("unauthorized to delete this post")
		}
	}

	// Delete the post
	if err := s.Delete(postID); err != nil {
		return err
	}

	// Log the deletion in audit trail
	if s.auditService != nil {
		auditData := AuditEventData{
			EntityType: "post",
			EntityID:   string(rune(postID)),
			OldValues: map[string]interface{}{
				"title":   existingPost.Title,
				"content": existingPost.Content,
			},
		}
		s.auditService.LogEvent(userID, ActionDelete, auditData)
	}

	return nil
}

// GetUserPosts gets all posts by a specific user
func (s *PostService) GetUserPosts(userID uint, options QueryOptions) (*PaginatedResult[models.Post], error) {
	conditions := map[string]interface{}{
		"user_id": userID,
	}
	
	// Add User preload to options if not already present
	found := false
	for _, preload := range options.Preload {
		if preload == "User" {
			found = true
			break
		}
	}
	if !found {
		options.Preload = append(options.Preload, "User")
	}
	
	return s.FindMany(conditions, options)
}

// GetPublishedPosts gets all published posts
func (s *PostService) GetPublishedPosts(options QueryOptions) (*PaginatedResult[models.Post], error) {
	// Add User preload to options if not already present
	found := false
	for _, preload := range options.Preload {
		if preload == "User" {
			found = true
			break
		}
	}
	if !found {
		options.Preload = append(options.Preload, "User")
	}
	
	return s.GetAll(options)
}

// SearchPosts searches posts by title and content
func (s *PostService) SearchPosts(query string, options QueryOptions) (*PaginatedResult[models.Post], error) {
	options.Search = query
	
	// Add User preload to options if not already present
	found := false
	for _, preload := range options.Preload {
		if preload == "User" {
			found = true
			break
		}
	}
	if !found {
		options.Preload = append(options.Preload, "User")
	}
	
	return s.GetAll(options)
}

// GetPostStats returns post statistics
func (s *PostService) GetPostStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total posts
	totalPosts, err := s.Count(map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	stats["total_posts"] = totalPosts

	// Recent posts (last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	var recentCount int64
	err = s.db.Model(&models.Post{}).
		Where("created_at > ?", thirtyDaysAgo).
		Count(&recentCount).Error
	if err != nil {
		return nil, err
	}
	stats["recent_posts"] = recentCount

	// Posts by user count
	var userPostCounts []struct {
		UserID uint  `json:"user_id"`
		Count  int64 `json:"count"`
	}
	err = s.db.Model(&models.Post{}).
		Select("user_id, COUNT(*) as count").
		Group("user_id").
		Order("count DESC").
		Limit(10).
		Find(&userPostCounts).Error
	if err != nil {
		return nil, err
	}
	stats["top_authors"] = userPostCounts

	return stats, nil
}

// GetPostsByDateRange gets posts within a specific date range
func (s *PostService) GetPostsByDateRange(startDate, endDate time.Time, options QueryOptions) (*PaginatedResult[models.Post], error) {
	// Add date range filter
	if options.Filter.Filters == nil {
		options.Filter.Filters = make(map[string]interface{})
	}
	
	options.Filter.Filters["created_at"] = map[string]interface{}{
		"from": startDate,
		"to":   endDate,
	}
	
	// Add User preload to options if not already present
	found := false
	for _, preload := range options.Preload {
		if preload == "User" {
			found = true
			break
		}
	}
	if !found {
		options.Preload = append(options.Preload, "User")
	}
	
	return s.GetAll(options)
}

// BulkDeletePosts deletes multiple posts (admin only)
func (s *PostService) BulkDeletePosts(postIDs []uint, userID uint) error {
	// Check if user is admin
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("unauthorized")
	}
	if user.Role != models.RoleAdmin {
		return errors.New("only admins can perform bulk operations")
	}

	// Convert to []interface{} for the generic method
	ids := make([]interface{}, len(postIDs))
	for i, id := range postIDs {
		ids[i] = id
	}

	// Delete the posts
	if err := s.DeleteBatch(ids); err != nil {
		return err
	}

	// Log the bulk deletion in audit trail
	if s.auditService != nil {
		auditData := AuditEventData{
			EntityType: "post",
			NewValues: map[string]interface{}{
				"bulk_deleted_ids": postIDs,
				"count":            len(postIDs),
			},
		}
		s.auditService.LogEvent(userID, ActionDelete, auditData)
	}

	return nil
}