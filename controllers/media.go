package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sp3ctr4/database"
)

func HandleUpload(c *gin.Context) {

	/*
		HandleUpload performs the following actions:

		1. Retrieve and validate the authenticated user's information.

		2. Sanitize the file name, ensuring it's valid and secure.
		3. Validate the file type and size, ensuring it's a supported video format.
		4. Save the file to a structured directory under the user's ID.
		5. Insert a record into the `media` table with the file path, user ID, and type (public/private).
		6. Return a response with the file's URL/path or an error message.
		7. Implement thorough error handling and logging.
		8. Consider additional security measures like rate limiting, virus scanning, and CSRF protection.
	*/
	var user database.User
	contextUser, ok := c.Get("user")
	if !ok {
		c.JSON(404, gin.H{
			"data": gin.H{
				"message": "user not found",
			},
		})
		c.Abort()
		return
	}
	user = contextUser.(database.User)
	fmt.Println(user, "user")

	file, _ := c.FormFile("file")

	sanitizedFilename := filepath.Base(file.Filename)[:255]
	extension := filepath.Ext(file.Filename)
	videoExtensions := []string{
		".mp4",
		".avi",
		".mkv",
		".mov",
		".wmv",
		".flv",
		".webm",
		".m4v",
		".3gp",
		".mpg",
		".mpeg",
		".ogv",
		".vob",
	}

	fmt.Println(extension, "extension", videoExtensions, sanitizedFilename, file.Size)
	var isValidFileFormat bool
	isValidFileFormat = false
	for _, ext := range videoExtensions {
		if ext == extension {
			isValidFileFormat = true
			break

		}

	}
	if !isValidFileFormat || file.Size > 100_00_00_00 {
		c.JSON(http.StatusForbidden, gin.H{
			"data": "file format forbidden or file too large",
		})
		c.Abort()
		return
	}

	// get os home direcotory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": "failed to initalize storage",
		})
		c.Abort()
		return
	}
	userId := fmt.Sprint(user.ID)
	newUUID := uuid.New().String()
	destinationPath := filepath.Join(homeDir, "Desktop", "uploads", userId, newUUID+"-"+sanitizedFilename)
	fmt.Println("destination", destinationPath)

	// Upload the file to specific dst.
	c.SaveUploadedFile(file, destinationPath)

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
