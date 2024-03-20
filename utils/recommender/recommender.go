package recommender

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/anirudhgray/mood-harbour-backend/infra/database"
	"github.com/anirudhgray/mood-harbour-backend/models"
)

// GetUserLikesDislikes retrieves the resource IDs that a user has liked and disliked based on their review ratings.
// It queries the database for reviews associated with the given user ID and categorizes the resources into liked and disliked lists.
func GetUserLikesDislikes(userID uint) (liked []uint, notLiked []uint, err error) {
	var reviews []models.Review
	result := database.DB.Where("user_id = ?", userID).Find(&reviews)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	var likedResources []uint
	var notLikedResources []uint

	for _, review := range reviews {
		if review.Rating >= 4 {
			likedResources = append(likedResources, review.ResourceID)
		} else {
			notLikedResources = append(notLikedResources, review.ResourceID)
		}
	}

	return likedResources, notLikedResources, nil
}

// GetUsersWhoReviewedResource retrieves the user IDs of users who have reviewed a specific resource.
// It queries the database for reviews associated with the given resource ID and returns the user IDs of users who have reviewed the resource.
func GetUsersWhoReviewedResource(resourceID uint) (userslist []uint, err error) {
	var reviews []models.Review
	result := database.DB.Where("resource_id = ?", resourceID).Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}

	var users []uint

	for _, review := range reviews {
		users = append(users, review.UserID)
	}

	return users, nil
}

// GetUsersWithSimilarInteractions retrieves user IDs who have interacted with resources that are similar to the given user's liked and disliked resources.
// It queries the database for users who have reviewed resources that match the liked and disliked resources of the given user.
func GetUsersWithSimilarInteractions(likedResources []uint, dislikedResources []uint, userID uint) (userslist []uint, err error) {
	var users []uint

	for _, resourceID := range likedResources {
		reviewedUsers, err := GetUsersWhoReviewedResource(resourceID)
		if err != nil {
			return nil, err
		}
		users = append(users, reviewedUsers...)
	}

	for _, resourceID := range dislikedResources {
		reviewedUsers, err := GetUsersWhoReviewedResource(resourceID)
		if err != nil {
			return nil, err
		}
		users = append(users, reviewedUsers...)
	}

	// remove duplicates and remove the current user's ID
	currentUserID := userID
	uniqueUsers := make(map[uint]bool)
	var resultUsers []uint

	for _, userID := range users {
		if userID != currentUserID && !uniqueUsers[userID] {
			uniqueUsers[userID] = true
			resultUsers = append(resultUsers, userID)
		}
	}

	return resultUsers, nil
}

func formatFloat(num float64, prc int) string {
	var (
		zero, dot = "0", "."

		str = fmt.Sprintf("%."+strconv.Itoa(prc)+"f", num)
	)

	return strings.TrimRight(strings.TrimRight(str, zero), dot)
}

// CalculateUserSimilarity of two users via modified Jaccard Coefficient.
//
//	S(U1, U2) = (|L1 intersec L2| + |D1 intersect D2| - |L1 intersect D2| - |L2 intersect D1|) / |L1 union L2 union D1 union D2|
func CalculateUserSimilarity(currentUserLiked, currentUserDisliked, otherUserLiked, otherUserDisliked []uint) (similarityCoefficient float64) {
	L1, L2, D1, D2 := currentUserLiked, otherUserLiked, currentUserDisliked, otherUserDisliked

	L1IntersectL2Size := intersectionSize(currentUserLiked, otherUserLiked)
	D1IntersectD2Size := intersectionSize(currentUserDisliked, otherUserDisliked)
	L1IntersectD2Size := intersectionSize(currentUserLiked, otherUserDisliked)
	L2IntersectD1Size := intersectionSize(otherUserLiked, currentUserDisliked)

	numerator := float64(L1IntersectL2Size + D1IntersectD2Size - L1IntersectD2Size - L2IntersectD1Size)
	denominator := float64(unionSize4(L1, L2, D1, D2))

	if numerator == 0 {
		return 0
	}

	similarity := numerator / denominator

	res, _ := strconv.ParseFloat(formatFloat(similarity, 3), 64)

	return res
}

func intersectionSize(set1, set2 []uint) int {
	count := 0
	set2Map := make(map[uint]bool)
	for _, item := range set2 {
		set2Map[item] = true
	}
	for _, item := range set1 {
		if set2Map[item] {
			count++
		}
	}
	return count
}

func unionSize4(arr1, arr2, arr3, arr4 []uint) int {
	unionSet := make(map[uint]bool)

	for _, item := range arr1 {
		unionSet[item] = true
	}
	for _, item := range arr2 {
		unionSet[item] = true
	}
	for _, item := range arr3 {
		unionSet[item] = true
	}
	for _, item := range arr4 {
		unionSet[item] = true
	}

	return len(unionSet)
}

// CalculateSimilaritiesWithOtherUsers returns a mapping of similar users (returned by GetUsersWithSimilarInteractions) to their similarity coefficient with current user.
func CalculateSimilaritiesWithOtherUsers(currentUserID uint, similarUsers []uint, currentUserLiked, currentUserDisliked []uint) (sims map[uint]float64, err error) {
	userSimilarities := make(map[uint]float64)

	for _, otherUserID := range similarUsers {
		otherUserLiked, otherUserDisliked, err := GetUserLikesDislikes(otherUserID)
		if err != nil {
			return nil, err
		}

		similarity := CalculateUserSimilarity(currentUserLiked, currentUserDisliked, otherUserLiked, otherUserDisliked)
		userSimilarities[otherUserID] = similarity
	}

	return userSimilarities, nil
}

// GetUnreviewedResources fetches resources which the current user has not yet reviewed (will select recommendations from among these)
func GetUnreviewedResources(currentUserID uint) ([]uint, error) {
	userLikedResources, userDislikedResources, err := GetUserLikesDislikes(currentUserID)
	if err != nil {
		return nil, err
	}

	userReviewedResourceIDs := append(userLikedResources, userDislikedResources...)

	var unreviewedResources []uint

	var allResourceIDs []uint
	result := database.DB.Model(&models.Resource{}).Pluck("id", &allResourceIDs)
	if result.Error != nil {
		return nil, result.Error
	}

	// Filter out the resources that the user has reviewed
	for _, resourceID := range allResourceIDs {
		if !contains(userReviewedResourceIDs, resourceID) {
			unreviewedResources = append(unreviewedResources, resourceID)
		}
	}

	return unreviewedResources, nil
}

func contains(arr []uint, item uint) bool {
	for _, val := range arr {
		if val == item {
			return true
		}
	}
	return false
}

// CalculateRecommendationProbabilities returns mapping on user's unreviewed resources to the "probability" of them liking it (i.e., order of recommendation).
//
// ZL = sum of similarity coefficients of other similar users who have liked a particular resource.
//
// ML = number of such users as above
//
// ZD = sum of similarity coefficients of other similar users who have disliked a particular resource.
//
// MD = number of such users as above
//
//	P(U, B) = (ZL - ZD) / (|ML| + |MD|)
//
// Produces a value between -1 and 1
func CalculateRecommendationProbabilities(currentUserID uint, unreviewedResources []uint, similarUsers map[uint]float64) map[uint]float64 {
	recommendationProbabilities := make(map[uint]float64)

	for _, resourceID := range unreviewedResources {
		ZL := 0.0
		ZD := 0.0
		ML := 0
		MD := 0

		likedByUsers, dislikedByUsers, err := GetLikersDislikersForResource(resourceID)
		if err != nil {
			continue
		}

		for _, userID := range likedByUsers {
			if similarity, ok := similarUsers[userID]; ok {
				ZL += similarity
				ML++
			}
		}

		for _, userID := range dislikedByUsers {
			if similarity, ok := similarUsers[userID]; ok {
				ZD += similarity
				MD++
			}
		}

		if ML+MD == 0 {
			continue
		}

		probability := (ZL - ZD) / float64(ML+MD)
		recommendationProbabilities[resourceID] = probability
	}

	return recommendationProbabilities
}

func GetLikersDislikersForResource(resourceID uint) ([]uint, []uint, error) {
	var likedByUsers []uint
	var dislikedByUsers []uint

	usersWhoLiked, err := GetUsersWhoLikedResource(resourceID)
	if err != nil {
		return nil, nil, err
	}
	likedByUsers = append(likedByUsers, usersWhoLiked...)

	// Fetch users who disliked the resource
	usersWhoDisliked, err := GetUsersWhoDislikedResource(resourceID)
	if err != nil {
		return nil, nil, err
	}
	dislikedByUsers = append(dislikedByUsers, usersWhoDisliked...)

	return likedByUsers, dislikedByUsers, nil
}

func GetUsersWhoLikedResource(resourceID uint) ([]uint, error) {
	var likedByUsers []uint

	var reviews []models.Review
	result := database.DB.Where("resource_id = ? AND rating >= 4", resourceID).Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, review := range reviews {
		likedByUsers = append(likedByUsers, review.UserID)
	}

	return likedByUsers, nil
}

func GetUsersWhoDislikedResource(resourceID uint) ([]uint, error) {
	var dislikedByUsers []uint

	var reviews []models.Review
	result := database.DB.Where("resource_id = ? AND rating < 4", resourceID).Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, review := range reviews {
		dislikedByUsers = append(dislikedByUsers, review.UserID)
	}

	return dislikedByUsers, nil
}

type Recommendation struct {
	Resource    models.Resource
	Probability float64
}

func GetRecommendedResourcesSortedAndPaginated(recommendationProbabilities map[uint]float64, page int, perPage int) []Recommendation {
	var recommendations []Recommendation

	for resourceID, probability := range recommendationProbabilities {
		var resource models.Resource
		database.DB.Model(&models.Resource{}).First(&resource, resourceID)
		recommendations = append(recommendations, Recommendation{
			Resource:    resource,
			Probability: probability,
		})
	}

	// Sort recommendations by probability in descending order
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Probability > recommendations[j].Probability
	})

	// Paginate the recommendations
	startIdx := (page - 1) * perPage
	endIdx := int(math.Min(float64(startIdx+perPage), float64(len(recommendations))))

	if startIdx < len(recommendations) {
		return recommendations[startIdx:endIdx]
	}
	return []Recommendation{}
}
