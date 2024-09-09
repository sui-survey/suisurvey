package model

type BlobUploadResponse struct {
	NewlyCreated struct {
		BlobObject struct {
			BlobID string `json:"blob_id"`
		} `json:"blob_object"`
	} `json:"newly_created"`
}
