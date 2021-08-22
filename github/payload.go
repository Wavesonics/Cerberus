package github

type GithubWebhookPayload struct {
	Ref string `json:"ref"`
}
