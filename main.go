package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const PROJECT_ROOT = "./projects"

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")
	// r.Use(cors.Default())

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.POST("/compile", compileLatex)
	r.GET("/download/:filename", downloadPDF)
	r.GET("/files", getFiles)
	r.POST("/files", createFile)
	r.GET("/files/:filename", getFileContent)
	r.PUT("/files/:filename", saveFile)

	r.Run(":8080")
}

func getFiles(c *gin.Context) {
	files, err := listFiles(PROJECT_ROOT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list files"})
		return
	}
	c.JSON(http.StatusOK, files)
}

func listFiles(dir string) ([]gin.H, error) {
	var files []gin.H
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(PROJECT_ROOT, path)
		files = append(files, gin.H{
			"name": info.Name(),
			"path": relPath,
			"type": getFileType(info),
		})
		return nil
	})
	return files, err
}

func getFileType(info os.FileInfo) string {
	if info.IsDir() {
		return "directory"
	}
	return "file"
}

func createFile(c *gin.Context) {
	var input struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Path string `json:"path"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fullPath := filepath.Join(PROJECT_ROOT, input.Path, input.Name)

	if input.Type == "directory" {
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
			return
		}
	} else {
		if _, err := os.Create(fullPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "File or directory created successfully"})
}

func getFileContent(c *gin.Context) {
	filename := c.Param("filename")
	content, err := ioutil.ReadFile(filepath.Join(PROJECT_ROOT, filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": string(content)})
}

func saveFile(c *gin.Context) {
	filename := c.Param("filename")
	var input struct {
		Content string `json:"content"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := ioutil.WriteFile(filepath.Join(PROJECT_ROOT, filename), []byte(input.Content), 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File saved successfully"})
}

func compileLatex(c *gin.Context) {
    var input struct {
        LaTeX string `json:"latex"`
    }

    if err := c.BindJSON(&input); err != nil {
        log.Printf("Error binding JSON: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }


    // Create a temporary directory
    tmpDir, err := ioutil.TempDir("", "latex-")
    if err != nil {
        log.Printf("Failed to create temporary directory: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary directory"})
        return
    }
    defer os.RemoveAll(tmpDir)

    tmpFile := filepath.Join(tmpDir, "input.tex")
    if err := ioutil.WriteFile(tmpFile, []byte(input.LaTeX), 0644); err != nil {
        log.Printf("Failed to write to temporary file: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write to temporary file"})
        return
    }

    cmd := exec.Command("pdflatex", "-interaction=nonstopmode", "-aux-directory="+tmpDir, "-output-directory="+tmpDir, tmpFile)
    env := os.Environ()
    env = append(env, "MIKTEX_ENABLE_INSTALLER=1")
    cmd.Env = env
    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Printf("LaTeX compilation failed: %v\nOutput: %s", err, string(output))
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("LaTeX compilation failed: %s", string(output))})
        return
    }

    pdfPath := filepath.Join(tmpDir, "input.pdf")
    pdfContent, err := ioutil.ReadFile(pdfPath)
    if err != nil {
        log.Printf("Failed to read generated PDF: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read generated PDF"})
        return
    }

    pdfBase64 := base64.StdEncoding.EncodeToString(pdfContent)

    downloadFilename := fmt.Sprintf("latex_output_%d.pdf", time.Now().UnixNano())
    
    downloadDir := "./downloads"
    os.MkdirAll(downloadDir, os.ModePerm)
    downloadPath := filepath.Join(downloadDir, downloadFilename)
    if err := os.Rename(pdfPath, downloadPath); err != nil {
        log.Printf("Failed to prepare PDF for download: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare PDF for download"})
        return
    }

    result := gin.H{
        "pdf": pdfBase64,
        "downloadUrl": "/download/" + downloadFilename,
    }
    log.Printf("Compilation result: %+v", result)
    c.JSON(http.StatusOK, result)
}

func downloadPDF(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join("./downloads", filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "PDF not found"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/pdf")
	c.File(filePath)

	// Delete the file after download
	go func() {
		time.Sleep(5 * time.Second)
		os.Remove(filePath)
	}()
}