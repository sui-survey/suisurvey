package model

type Storage struct {
	ID          string `json:"id"`
	StartEpoch  int    `json:"startEpoch"`
	EndEpoch    int    `json:"endEpoch"`
	StorageSize int    `json:"storageSize"`
}

type BlobObject struct {
	ID              string  `json:"id"`
	StoredEpoch     int     `json:"storedEpoch"`
	BlobId          string  `json:"blobId"`
	Size            int     `json:"size"`
	ErasureCodeType string  `json:"erasureCodeType"`
	CertifiedEpoch  int     `json:"certifiedEpoch"`
	Storage         Storage `json:"storage"`
}

type NewlyCreated struct {
	BlobObject  BlobObject `json:"blobObject"`
	EncodedSize int        `json:"encodedSize"`
	Cost        int        `json:"cost"`
}

type WalrusResponse struct {
	NewlyCreated NewlyCreated `json:"newlyCreated"`
}
