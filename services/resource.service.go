package services

import (
	"errors"

	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/anirudhgray/mood-harbour-backend/repository"
)

type ResourceService struct {
	resourceRepo repository.ResourceRepositoryInterface
	userRepo     repository.UserRepositoryInterface
}

func NewResourceService(resourceRepo repository.ResourceRepositoryInterface, userRepo repository.UserRepositoryInterface) *ResourceService {
	return &ResourceService{resourceRepo, userRepo}
}

type ResourceServiceInterface interface {
	CreateResourceEntry(userID uint, title, content, url string, external, adminPost bool) (models.ResourceResponse, error)
	GetAllResources() ([]models.ResourceResponse, error)
	GetResourceByID(resourceID uint) (models.ResourceResponse, error)
	DeleteResource(userID, resourceID uint) error
	UpdateResource(resourceID uint, userID uint, title, content, url string, external, adminPost bool) (models.ResourceResponse, error)
	GetAdminResources() ([]models.ResourceResponse, error)
	AddReview(resourceID, userID uint, content string, rating models.Rating) error
}

// CreateResourceEntry creates a new resource entry in the database.
func (rs *ResourceService) CreateResourceEntry(userID uint, title, content, url string, external, adminPost bool) (models.ResourceResponse, error) {
	resource := models.Resource{
		CreatedBy: userID,
		Title:     title,
		Content:   content,
		URL:       url,
		External:  external,
		AdminPost: adminPost,
	}

	err := rs.resourceRepo.CreateResource(&resource)
	if err != nil {
		return models.ResourceResponse{}, err
	}

	reviews, err := rs.resourceRepo.GetReviewsByResourceID(resource.ID)
	if err != nil {
		return models.ResourceResponse{}, err
	}

	return models.ResourceResponse{
		ID:       resource.ID,
		Resource: resource,
		Reviews:  reviews,
	}, nil
}

// GetAllResources gets all the resources in the database.
func (rs *ResourceService) GetAllResources() ([]models.ResourceResponse, error) {
	resources, err := rs.resourceRepo.GetResources()
	if err != nil {
		return nil, err
	}

	var resourceResponses []models.ResourceResponse
	for _, resource := range resources {
		reviews, err := rs.resourceRepo.GetReviewsByResourceID(resource.ID)
		if err != nil {
			return nil, err
		}

		resourceResponses = append(resourceResponses, models.ResourceResponse{
			ID:       resource.ID,
			Resource: resource,
			Reviews:  reviews,
		})
	}

	return resourceResponses, nil
}

// GetResourceByID gets a resource by its ID.
func (rs *ResourceService) GetResourceByID(resourceID uint) (models.ResourceResponse, error) {
	resource, err := rs.resourceRepo.GetResourceByID(resourceID)
	if err != nil {
		return models.ResourceResponse{}, err
	}

	reviews, err := rs.resourceRepo.GetReviewsByResourceID(resource.ID)
	if err != nil {
		return models.ResourceResponse{}, err
	}

	return models.ResourceResponse{
		ID:       resource.ID,
		Resource: resource,
		Reviews:  reviews,
	}, nil
}

// DeleteResource deletes a resource by its ID.
func (rs *ResourceService) DeleteResource(userID, resourceID uint) error {
	// Check if the user is the owner of the resource or is an admin
	resource, err := rs.resourceRepo.GetResourceByID(resourceID)
	if err != nil {
		return err
	}

	user, err := rs.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if resource.CreatedBy != user.ID && !user.Admin {
		// return error unauthorized
		return errors.New("unauthorized")
	}

	return rs.resourceRepo.DeleteResource(resourceID)
}

// UpdateResource updates a resource in the database.
func (rs *ResourceService) UpdateResource(resourceID uint, userID uint, title, content, url string, external, adminPost bool) (models.ResourceResponse, error) {
	resource := models.Resource{
		CreatedBy: userID,
		Title:     title,
		Content:   content,
		URL:       url,
		External:  external,
		AdminPost: adminPost,
	}

	updatedResource, err := rs.resourceRepo.UpdateResource(&resource)
	if err != nil {
		return models.ResourceResponse{}, err
	}

	reviews, err := rs.resourceRepo.GetReviewsByResourceID(resourceID)
	if err != nil {
		return models.ResourceResponse{}, err
	}

	return models.ResourceResponse{
		ID:       updatedResource.ID,
		Resource: updatedResource,
		Reviews:  reviews,
	}, nil
}

// GetAdminResources gets all the resources for an admin.
func (rs *ResourceService) GetAdminResources() ([]models.ResourceResponse, error) {
	resources, err := rs.resourceRepo.GetAdminResources()
	if err != nil {
		return nil, err
	}

	var resourceResponses []models.ResourceResponse
	for _, resource := range resources {
		reviews, err := rs.resourceRepo.GetReviewsByResourceID(resource.ID)
		if err != nil {
			return nil, err
		}

		resourceResponses = append(resourceResponses, models.ResourceResponse{
			ID:       resource.ID,
			Resource: resource,
			Reviews:  reviews,
		})
	}

	return resourceResponses, nil
}

// AddReview adds a review to a resource.
func (rs *ResourceService) AddReview(resourceID, userID uint, content string, rating models.Rating) error {
	review := models.Review{
		ResourceID: resourceID,
		UserID:     userID,
		Content:    content,
		Rating:     rating,
	}

	return rs.resourceRepo.AddReview(resourceID, &review)
}
