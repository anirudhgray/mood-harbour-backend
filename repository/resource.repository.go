package repository

import (
	"github.com/anirudhgray/mood-harbour-backend/infra/database"
	"github.com/anirudhgray/mood-harbour-backend/models"
	"gorm.io/gorm"
)

type ResourceRepository struct {
	db *gorm.DB
}

func NewResourceRepository() *ResourceRepository {
	return &ResourceRepository{database.DB}
}

// ResourceRepositoryInterface is the interface for the ResourceRepository.
type ResourceRepositoryInterface interface {
	CreateResource(resource *models.Resource) error
	GetResourceByID(resourceID uint) (models.Resource, error)
	GetResources() ([]models.Resource, error)
	GetResourcesByUserID(userID uint) ([]models.Resource, error)
	GetAdminResources() ([]models.Resource, error)
	DeleteResource(resourceID uint) error
	UpdateResource(resource *models.Resource) (models.Resource, error)
	AddReview(resourceID uint, review *models.Review) error
	GetReviewsByResourceID(resourceID uint) ([]models.Review, error)
	DeleteReview(reviewID uint) error
	UpdateReview(review *models.Review) error
	GetReviewByID(reviewID uint) (models.Review, error)
	GetReviewsByUserID(userID uint) ([]models.Review, error)
}

// CreateResource creates a new resource in the database.
func (rr *ResourceRepository) CreateResource(resource *models.Resource) error {
	return rr.db.Create(resource).Error
}

// GetResourceByID gets a resource by its ID.
func (rr *ResourceRepository) GetResourceByID(resourceID uint) (models.Resource, error) {
	var resource models.Resource
	err := rr.db.Where("id = ?", resourceID).First(&resource).Error
	return resource, err
}

// GetResources gets all the resources in the database.
func (rr *ResourceRepository) GetResources() ([]models.Resource, error) {
	var resources []models.Resource
	err := rr.db.Find(&resources).Error
	return resources, err
}

// GetResourcesByUserID gets all the resources for a specific user.
func (rr *ResourceRepository) GetResourcesByUserID(userID uint) ([]models.Resource, error) {
	var resources []models.Resource
	err := rr.db.Where("user_id = ?", userID).Find(&resources).Error
	return resources, err
}

// GetAdminResources gets all the resources for an admin.
func (rr *ResourceRepository) GetAdminResources() ([]models.Resource, error) {
	var resources []models.Resource
	err := rr.db.Where("admin_post = ?", true).Find(&resources).Error
	return resources, err
}

// DeleteResource deletes a resource by its ID.
func (rr *ResourceRepository) DeleteResource(resourceID uint) error {
	return rr.db.Where("id = ?", resourceID).Delete(&models.Resource{}).Error
}

// UpdateResource updates a resource in the database.
func (rr *ResourceRepository) UpdateResource(resource *models.Resource) (models.Resource, error) {
	err := rr.db.Save(resource).Error
	return *resource, err
}

// AddReview adds a review to a resource.
func (rr *ResourceRepository) AddReview(resourceID uint, review *models.Review) error {
	return rr.db.Create(review).Error
}

// GetReviewsByResourceID gets all the reviews for a specific resource.
func (rr *ResourceRepository) GetReviewsByResourceID(resourceID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := rr.db.Where("resource_id = ?", resourceID).Find(&reviews).Error
	return reviews, err
}

// DeleteReview deletes a review by its ID.
func (rr *ResourceRepository) DeleteReview(reviewID uint) error {
	return rr.db.Where("id = ?", reviewID).Delete(&models.Review{}).Error
}

// UpdateReview updates a review in the database.
func (rr *ResourceRepository) UpdateReview(review *models.Review) error {
	return rr.db.Save(review).Error
}

// GetReviewByID gets a review by its ID.
func (rr *ResourceRepository) GetReviewByID(reviewID uint) (models.Review, error) {
	var review models.Review
	err := rr.db.Where("id = ?", reviewID).First(&review).Error
	return review, err
}

// GetReviewsByUserID gets all the reviews for a specific user.
func (rr *ResourceRepository) GetReviewsByUserID(userID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := rr.db.Where("user_id = ?", userID).Find(&reviews).Error
	return reviews, err
}
