package model

type ImageStoreData struct {
	TenantId           string            `json:"tenant_id"`
	Searchterm         string            `json:"search_term"`
	SearchTermAlias    string            `json:"search_term_alias"`
	StoreType          string            `json:"store_type"`
	ImageCount         int               `json:"image_count"`
	StoreUrlByImageURL map[string]string `json:"bucket"`
}

type S3UploadRequest struct {
	Bucket    string `json:"bucket"`
	TenantId  string `json:"tenant_id"`
	Directory string `json:"directory"`
	FilePath  string `json:"filepath"`
	Region    string `json:"region"`
}
