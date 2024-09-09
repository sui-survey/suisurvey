package main

import (
	"kevinsheeran/walrus-backend/router"
)

// @title Survey
// @version 1.0
// @description backend survey for walrus
// @securityDefinition.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	//global.Log = core.InitLogger()
	// Read the YAML file
	//data, err := ioutil.ReadFile("config.yaml")
	//if err != nil {
	//	log.Fatalf("Error reading YAML file: %s", err)
	//}
	//
	//// Unmarshal the YAML data into the Config struct
	//err = yaml.Unmarshal(data, &config.Config)
	//if err != nil {
	//	log.Fatalf("Error parsing YAML file: %s", err)
	//}
	router := router.InitRouter()

	router.Run("0.0.0.0:31411")

}
