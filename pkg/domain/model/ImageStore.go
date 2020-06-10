package model

type ImageStoreData struct {
	TenantId           string
	Searchterm         string
	SearchTermAlias    string
	StoreType          string
	ImageCount         int
	StoreUrlByImageURL map[string]string
}

type S3UploadRequest struct {
	Bucket   string `json:"bucket"`
	TenantId string `json:"tenant_id"`
	Directory   string `json:"directory"`
	Region   string `json:"region"`
}
