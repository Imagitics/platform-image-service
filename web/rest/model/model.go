package model

type S3UploadRequest struct {
	Bucket   string `json:"bucket"`
	TenantId string `json:"tenant_id"`
	File     []byte `json:"entity`
	Region   string `json:"region`
}

type ImageURL struct {
	Url  string `json:"Url"`
	Name string `json:"Identifier"`
}

type SearchAPIImageResponse struct {
	TotalResults int        `json:"totalResults"`
	SearchUrls   []ImageURL `json:"searchResults"`
}

type SearchAPIImageRequest struct {
	TenantID    string `json:"tenantId"`
	SearchTerm  string `json:"searchTerm"`
	IncludeFace bool   `json:"includeFace"`
	SearchAlias string `json:"searchAlias"`
}
