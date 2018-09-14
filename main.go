package main

import (
	//"bytes"
	"log"
	"strings"
	//"context"
	//"fmt"
	//"github.com/docker/docker/api/types"
	//"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	//"io/ioutil"
	//"encoding/json"
	//"github.com/json-iterator/go"
	"net/http"
	"os/exec"
)

func main() {
	log.SetOutput(gin.DefaultWriter)
	router := gin.Default()
	router.GET("/healthz", health_check)
	router.GET("/container/:name", get_container)
	router.Run(":8080")
}

func health_check(c *gin.Context) {
	contCmd := exec.Command("docker", "version")

	_, err := contCmd.Output()
	if err != nil {
		log.Println("Command Error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, "ok")
}

func get_container(c *gin.Context) {
	contName := c.Param("name")

	contCmd := exec.Command(
		"docker",
		"ps",
		"-f",
		"name="+contName,
		"--format",
		"'{{.ID}}-{{.Names}}'",
	)

	dataOut, err := contCmd.Output()

	if err != nil {
		log.Println("Command Error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if len(dataOut) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Container name not found"})
	} else {
		conts := strings.Split(string(dataOut), "\n")

		var msg string
		for _, cont := range conts {
			cont = strings.Trim(cont, "'")
			log.Println("Cont: ", cont)
			if len(cont) > 0 && strings.HasSuffix(cont, contName) {
				log.Println("Found: ", cont)
				msg = cont
			} else {
				log.Println("NotFound")
			}
		}
		log.Println("Message: ", msg)
		if len(msg) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Container name not found"})
		} else {
			c.JSON(http.StatusOK, gin.H{"data": msg})
		}

	}
}
