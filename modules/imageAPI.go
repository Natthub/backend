package modules

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
)

func getBase64(path string) string {
	f, _ := os.Open(path)
	// Read entire JPG into byte slice.
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)
	// Print encoded data to console.
	// ... The base64 image can be used as a data URI in a browser.
	fmt.Println("ENCODED: " + encoded)
	return encoded
}

func GetProductImage(c *gin.Context) {
	filename := c.Param("filename")
	path := "./image_product/" + filename
	c.File(path)
}

func GetReceiptImage(c *gin.Context) {
	filename := c.Param("filename")
	path := "./image_receipt/" + filename
	c.File(path)
}
