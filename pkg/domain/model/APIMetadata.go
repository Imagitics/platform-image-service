package model

type APIMetadata struct {
	TenantID string            `json:"tenant_id"`
	APIName  string            `json:"api_name"`
	Params   map[string]string `json:"parameters"`
}

func (instance *APIMetadata) getParams(key string) string {
	return instance.Params[key]
}
