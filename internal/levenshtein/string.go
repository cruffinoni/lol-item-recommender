package levenshtein

func StringSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 100.0
	}

	len1 := len(s1)
	len2 := len(s2)
	if len1 == 0 || len2 == 0 {
		return 0.0
	}

	maxLength := len1
	if len2 > maxLength {
		maxLength = len2
	}

	distance := levenshteinDistance(s1, s2)
	similarity := 100.0 - (float64(distance) / float64(maxLength) * 100.0)
	return similarity
}

func levenshteinDistance(s1, s2 string) int {
	len1 := len(s1)
	len2 := len(s2)

	// Créer la matrice de distances
	matrix := make([][]int, len1+1)
	for i := range matrix {
		matrix[i] = make([]int, len2+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	// Calculer les distances
	for j := 1; j <= len2; j++ {
		for i := 1; i <= len1; i++ {
			if s1[i-1] == s2[j-1] {
				matrix[i][j] = matrix[i-1][j-1]
			} else {
				min := matrix[i-1][j]
				if matrix[i][j-1] < min {
					min = matrix[i][j-1]
				}
				if matrix[i-1][j-1] < min {
					min = matrix[i-1][j-1]
				}
				matrix[i][j] = min + 1
			}
		}
	}

	// Retourner la distance entre les deux chaînes
	return matrix[len1][len2]
}
