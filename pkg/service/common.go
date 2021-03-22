package service

type hotComment struct {
	CommentatorName string
	CommentContent  string
}

type Comment struct {
	Data Data `json:"data"`
}

type Data struct {
	Html string `json:"html"`
}

func sliceToString(slice []hotComment) string {
	var result string
	for _, comment := range slice {
		result += comment.CommentatorName + "\n" + comment.CommentContent + "\n\n"
	}

	return result
}
