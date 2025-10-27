package slsk

func EvaluateFile(file File) int {
	score := 0

	if file.SampleRate < 44000 {
		score -= 1
	} else if file.SampleRate > 44400 {
		score += 1
	}

	if file.BitDepth == 16 {
		score += 1
	} else if file.BitDepth == 24 {
		score += 2
	}

	if file.Extension == "flac" {
		score += 1
	} else if file.Extension == "mp3" {
		score -= 1
	}

	return score
}

func EvaluateFileList(files []File) float64 {
	scoreList := []int{}

	for _, f := range files {
		scoreList = append(scoreList, EvaluateFile(f))
	}

	sum := 0

	for _, s := range scoreList {
		sum += s
	}

	return float64(float64(sum) / float64(len(scoreList)))
}

// returns host name too
func GetBestFileList(releases []SearchResponses) ([]File, string) {
	bestRelease := []File{}
	var bestScore float64 = -99
	bestReleaseHost := ""

	for _, r := range releases {
		score := EvaluateFileList(r.Files)

		if score > float64(bestScore) {
			bestRelease = r.Files
			bestScore = score
			bestReleaseHost = r.Username
		}
	}

	return bestRelease, bestReleaseHost
}
