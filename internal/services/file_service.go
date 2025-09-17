package services

import (
	"fmt"
	"go-backend/internal/models"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FileService handles file upload, validation, and management
type FileService struct {
	db           *gorm.DB
	uploadPath   string
	maxFileSize  int64
	allowedTypes map[string]bool
	staticURL    string
	auditService *AuditService
}

// FileUploadConfig contains file upload configuration
type FileUploadConfig struct {
	UploadPath   string
	MaxFileSize  int64 // in bytes
	AllowedTypes []string
	StaticURL    string
}

// NewFileService creates a new file service instance
func NewFileService(db *gorm.DB, config FileUploadConfig, auditService *AuditService) *FileService {
	// Convert allowed types slice to map for faster lookup
	allowedTypesMap := make(map[string]bool)
	for _, fileType := range config.AllowedTypes {
		allowedTypesMap[strings.ToLower(fileType)] = true
	}

	// Ensure upload directory exists
	if err := os.MkdirAll(config.UploadPath, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}

	return &FileService{
		db:           db,
		uploadPath:   config.UploadPath,
		maxFileSize:  config.MaxFileSize,
		allowedTypes: allowedTypesMap,
		staticURL:    config.StaticURL,
		auditService: auditService,
	}
}

// FileValidationError represents file validation errors
type FileValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e FileValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// UploadResult contains the result of a file upload
type UploadResult struct {
	FileUpload *models.FileUpload `json:"file_upload"`
	URL        string             `json:"url"`
}

// ValidateFile validates a file before upload
func (s *FileService) ValidateFile(fileHeader *multipart.FileHeader) error {
	// Check file size
	if fileHeader.Size > s.maxFileSize {
		return FileValidationError{
			Field:   "file_size",
			Message: fmt.Sprintf("File size exceeds maximum allowed size of %d bytes", s.maxFileSize),
		}
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		return FileValidationError{
			Field:   "file_extension",
			Message: "File must have an extension",
		}
	}

	// Remove the dot from extension
	ext = ext[1:]
	if !s.allowedTypes[ext] {
		allowedList := make([]string, 0, len(s.allowedTypes))
		for fileType := range s.allowedTypes {
			allowedList = append(allowedList, fileType)
		}
		return FileValidationError{
			Field:   "file_type",
			Message: fmt.Sprintf("File type '%s' is not allowed. Allowed types: %s", ext, strings.Join(allowedList, ", ")),
		}
	}

	return nil
}

// UploadFile uploads a file and stores its metadata
func (s *FileService) UploadFile(fileHeader *multipart.FileHeader, userID uint, category string) (*UploadResult, error) {
	// Validate the file
	if err := s.ValidateFile(fileHeader); err != nil {
		return nil, err
	}

	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	// Generate unique filename
	ext := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Create full file path
	filePath := filepath.Join(s.uploadPath, fileName)

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy the uploaded file to destination
	_, err = io.Copy(dst, file)
	if err != nil {
		// Clean up the created file if copy fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Get file info
	fileInfo, err := dst.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Create file upload record
	fileUpload := &models.FileUpload{
		UserID:       userID,
		OriginalName: fileHeader.Filename,
		FileName:     fileName,
		FilePath:     filePath,
		FileSize:     fileInfo.Size(),
		MimeType:     fileHeader.Header.Get("Content-Type"),
		Category:     category,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save to database
	if err := s.db.Create(fileUpload).Error; err != nil {
		// Clean up the uploaded file if database save fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}

	// Generate public URL
	url := fmt.Sprintf("%s/%s", strings.TrimRight(s.staticURL, "/"), fileName)

	// Log the upload in audit trail
	if s.auditService != nil {
		auditData := AuditEventData{
			EntityType: "file_upload",
			EntityID:   fmt.Sprintf("%d", fileUpload.ID),
			NewValues: map[string]interface{}{
				"original_name": fileUpload.OriginalName,
				"file_name":     fileUpload.FileName,
				"file_size":     fileUpload.FileSize,
				"mime_type":     fileUpload.MimeType,
				"category":      fileUpload.Category,
			},
		}
		s.auditService.LogEvent(userID, ActionFileUpload, auditData)
	}

	return &UploadResult{
		FileUpload: fileUpload,
		URL:        url,
	}, nil
}

// GetFile retrieves file metadata by ID
func (s *FileService) GetFile(fileID uint) (*models.FileUpload, error) {
	var fileUpload models.FileUpload
	err := s.db.Preload("User").First(&fileUpload, fileID).Error
	return &fileUpload, err
}

// GetUserFiles retrieves all files uploaded by a specific user
func (s *FileService) GetUserFiles(userID uint, category string, limit, offset int) ([]models.FileUpload, error) {
	var files []models.FileUpload
	query := s.db.Where("user_id = ?", userID)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&files).Error

	return files, err
}

// GetFilesByCategory retrieves files by category
func (s *FileService) GetFilesByCategory(category string, limit, offset int) ([]models.FileUpload, error) {
	var files []models.FileUpload
	err := s.db.Where("category = ?", category).
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&files).Error

	return files, err
}

// DeleteFile deletes a file and its metadata
func (s *FileService) DeleteFile(fileID, userID uint) error {
	// Get the file record
	var fileUpload models.FileUpload
	if err := s.db.First(&fileUpload, fileID).Error; err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Check if user owns the file or is admin
	if fileUpload.UserID != userID {
		var user models.User
		if err := s.db.First(&user, userID).Error; err != nil {
			return fmt.Errorf("unauthorized")
		}
		if user.Role != models.RoleAdmin {
			return fmt.Errorf("unauthorized to delete this file")
		}
	}

	// Delete the physical file
	if err := os.Remove(fileUpload.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete physical file: %w", err)
	}

	// Delete the database record
	if err := s.db.Delete(&fileUpload).Error; err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	// Log the deletion in audit trail
	if s.auditService != nil {
		auditData := AuditEventData{
			EntityType: "file_upload",
			EntityID:   fmt.Sprintf("%d", fileID),
			OldValues: map[string]interface{}{
				"original_name": fileUpload.OriginalName,
				"file_name":     fileUpload.FileName,
				"file_size":     fileUpload.FileSize,
				"category":      fileUpload.Category,
			},
		}
		s.auditService.LogEvent(userID, ActionFileDownload, auditData)
	}

	return nil
}

// GetFileContent serves file content for download
func (s *FileService) GetFileContent(fileID, userID uint) (*models.FileUpload, *os.File, error) {
	// Get the file record
	var fileUpload models.FileUpload
	if err := s.db.First(&fileUpload, fileID).Error; err != nil {
		return nil, nil, fmt.Errorf("file not found: %w", err)
	}

	// Check if user owns the file or is admin (for private files)
	if fileUpload.UserID != userID {
		var user models.User
		if err := s.db.First(&user, userID).Error; err != nil {
			return nil, nil, fmt.Errorf("unauthorized")
		}
		if user.Role != models.RoleAdmin && user.Role != models.RoleModerator {
			return nil, nil, fmt.Errorf("unauthorized to access this file")
		}
	}

	// Open the file
	file, err := os.Open(fileUpload.FilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}

	// Log the download in audit trail
	if s.auditService != nil {
		auditData := AuditEventData{
			EntityType: "file_upload",
			EntityID:   fmt.Sprintf("%d", fileID),
			NewValues: map[string]interface{}{
				"action":        "download",
				"original_name": fileUpload.OriginalName,
			},
		}
		s.auditService.LogEvent(userID, ActionFileDownload, auditData)
	}

	return &fileUpload, file, nil
}

// UpdateFileMetadata updates file metadata (category, etc.)
func (s *FileService) UpdateFileMetadata(fileID, userID uint, updates map[string]interface{}) error {
	// Get the file record
	var fileUpload models.FileUpload
	if err := s.db.First(&fileUpload, fileID).Error; err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Check if user owns the file or is admin
	if fileUpload.UserID != userID {
		var user models.User
		if err := s.db.First(&user, userID).Error; err != nil {
			return fmt.Errorf("unauthorized")
		}
		if user.Role != models.RoleAdmin {
			return fmt.Errorf("unauthorized to update this file")
		}
	}

	// Add update timestamp
	updates["updated_at"] = time.Now()

	// Update the record
	if err := s.db.Model(&fileUpload).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update file metadata: %w", err)
	}

	return nil
}

// CleanupOrphanedFiles removes files that exist on disk but not in database
func (s *FileService) CleanupOrphanedFiles() error {
	// Get all files from upload directory
	files, err := filepath.Glob(filepath.Join(s.uploadPath, "*"))
	if err != nil {
		return fmt.Errorf("failed to read upload directory: %w", err)
	}

	// Get all filenames from database
	var dbFiles []models.FileUpload
	if err := s.db.Select("file_name").Find(&dbFiles).Error; err != nil {
		return fmt.Errorf("failed to query database files: %w", err)
	}

	// Create a map of database filenames for quick lookup
	dbFileMap := make(map[string]bool)
	for _, file := range dbFiles {
		dbFileMap[file.FileName] = true
	}

	// Remove orphaned files
	orphanedCount := 0
	for _, filePath := range files {
		fileName := filepath.Base(filePath)
		if !dbFileMap[fileName] {
			if err := os.Remove(filePath); err == nil {
				orphanedCount++
			}
		}
	}

	return nil
}

// GetFileStats returns file upload statistics
func (s *FileService) GetFileStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total files
	var totalFiles int64
	if err := s.db.Model(&models.FileUpload{}).Count(&totalFiles).Error; err != nil {
		return nil, err
	}
	stats["total_files"] = totalFiles

	// Total storage used
	var totalSize int64
	if err := s.db.Model(&models.FileUpload{}).Select("COALESCE(SUM(file_size), 0)").Scan(&totalSize).Error; err != nil {
		return nil, err
	}
	stats["total_storage_bytes"] = totalSize
	stats["total_storage_mb"] = float64(totalSize) / (1024 * 1024)

	// Files by category
	var categoryCounts []struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}
	if err := s.db.Model(&models.FileUpload{}).
		Select("category, COUNT(*) as count").
		Group("category").
		Find(&categoryCounts).Error; err != nil {
		return nil, err
	}
	stats["files_by_category"] = categoryCounts

	// Recent uploads (last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	var recentCount int64
	if err := s.db.Model(&models.FileUpload{}).
		Where("created_at > ?", thirtyDaysAgo).
		Count(&recentCount).Error; err != nil {
		return nil, err
	}
	stats["recent_uploads"] = recentCount

	return stats, nil
}

// GetUploadedFileURL returns the public URL for a file
func (s *FileService) GetUploadedFileURL(fileName string) string {
	return fmt.Sprintf("%s/%s", strings.TrimRight(s.staticURL, "/"), fileName)
}
