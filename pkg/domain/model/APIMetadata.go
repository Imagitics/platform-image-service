package model

type APIMetadata struct {
	TenantID string
	APIName  string
	Params   map[string]string
}

func (instance *APIMetadata) getParams(key string) string {
	return instance.Params[key]
}
