package utils

// CalculatePopularity computes the popularity score for a blog.
func CalculatePopularity(views, likes, dislikes, comments int) float64 {
	const (
		viewWeight    = 1.0
		likeWeight    = 3.0
		dislikeWeight = -2.0
		commentWeight = 2.0
	)
	return float64(views)*viewWeight + float64(likes)*likeWeight + float64(dislikes)*dislikeWeight + float64(comments)*commentWeight
}
