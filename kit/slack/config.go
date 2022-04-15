package slack

// Config for slack application
type Config struct {
	AppName        string `json:"app_name"`
	AppToken       string `json:"app_token"`
	WorkspaceToken string `json:"ws_token"`
}
