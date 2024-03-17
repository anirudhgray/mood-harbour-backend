package services

import (
	"github.com/anirudhgray/mood-harbour-backend/models"
	"github.com/anirudhgray/mood-harbour-backend/repository"
)

type ResourceService struct {
	resourceRepo repository.ResourceRepositoryInterface
}

func NewResourceService(resourceRepo repository.ResourceRepositoryInterface) *ResourceService {
	return &ResourceService{resourceRepo}
}

type ResourceServiceInterface interface {
	CreateResourceEntry(userID uint, title, content, url string, external, adminPost bool) (models.ResourceResponse, error)
	GetAllResources() ([]models.ResourceResponse, error)
	GetResourceByID(resourceID uint) (models.ResourceResponse, error)
	DeleteResource(resourceID uint) error
	UpdateResource(resourceID uint, userID uint, title, content, url string, external, adminPost bool) (models.ResourceResponse, error)
	GetAdminResources() ([]models.ResourceResponse, error)
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
func (rs *ResourceService) DeleteResource(resourceID uint) error {
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
