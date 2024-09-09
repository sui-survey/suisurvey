package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"kevinsheeran/walrus-backend/model"
	"kevinsheeran/walrus-backend/result"
	"net/http"
	"strings"
)

// CreateForm
// @Summary CreateForm API
// @Tags CreateForm
// @Produce json
// @Param data body model.CreateFormDto true "data"
// @Success 200 {object} result.Result
// @router /api/v1/create-form [post]
func CreateForm(c *gin.Context) {
	// Bind json parameters
	var dto model.CreateFormDto
	if err := c.BindJSON(&dto); err != nil {
		result.Failed(c, http.StatusBadRequest, "Invalid input")
	}
	// TODO: Check form duplication

	// Call walrus publisher API
	response, err := callWalrusPublisher(&dto)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Parse the response to extract blobId
	var walrusResponse model.WalrusResponse
	if err := json.Unmarshal([]byte(response), &walrusResponse); err != nil {
		result.Failed(c, http.StatusInternalServerError, "Error parsing response")
		return
	}

	blobId := walrusResponse.NewlyCreated.BlobObject.BlobId
	fmt.Printf("blobId: %s\n", blobId)

	for _, item := range dto.ItemList {
		fmt.Printf("item name: %s\n", item.Name)
		createSurvey(&dto, item.Name, blobId)
	}

	result.Success(c, response)

}

// GetForm
// @Summary GetForm API
// @Tags GetForm
// @Produce json
// @Param blobId path string true "Blob ID"
// @Success 200 {object} result.Result
// @Failure 400 {string} string "Invalid blob ID"
// @router /api/v1/form/{blobId} [get]
func GetForm(c *gin.Context) {

	// Bind json parameters
	blobId := c.Param("blobId")
	if blobId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Blob ID is required"})
		return
	}

	// Call aggregator API
	response, err := callWalrusAggregator(blobId)
	if err != nil {
		result.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(c, response)
}

func callWalrusPublisher(data *model.CreateFormDto) (string, error) {
	url := "https://publisher-devnet.walrus.space/v1/store"

	// Convert the DTO data to JSON format
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error marshalling data to JSON: %v", err)
	}

	req, err := http.NewRequest("PUT", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, error := client.Do(req)
	if error != nil {
		return "", fmt.Errorf("error executing request: %v", error)
	}
	defer response.Body.Close()

	// Read the response body
	body, error := ioutil.ReadAll(response.Body)
	if error != nil {
		return "", fmt.Errorf("error reading response: %v", error)
	}

	return string(body), nil
}

func callWalrusAggregator(blobId string) (string, error) {

	url := fmt.Sprintf("https://aggregator-devnet.walrus.space/v1/%s", blobId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("err creating request: %v", err)
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("err executing request: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("err reading response: %v", err)
	}

	return string(body), nil
}
