package core

// Kudo represents a oos kudo.
type Kudo struct {
	UserID      string `json:"user_id" bson:"userId"`
	RepoID      string `json:"id" bson:"repoId"`
	RepoName    string `json:"full_name" bson:"repoName"`
	RepoURL     string `json:"html_url" bson:"repoUrl"`
	Language    string `json:"language" bson:"language"`
	Description string `json:"description" bson:"description"`
	Notes       string `json:"notes" bson:"notes"`
}
