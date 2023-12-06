package main

import (
	"flag"
	"log"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"

	_ "embed"

	"github.com/gin-gonic/gin"
	"github.com/huyinghuan/simple-swagger-server/static"
)

func findJSONFiles(dirPath string) (map[string]string, error) {
	filesMap := make(map[string]string)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil // Skip directories
		}

		if filepath.Ext(path) == ".json" {
			// Get file name without extension
			fileName := strings.TrimSuffix(info.Name(), ".json")

			// Convert absolute path to relative path
			relPath, err := filepath.Rel(dirPath, path)
			if err != nil {
				return err
			}

			filesMap[fileName] = relPath
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filesMap, nil
}

func authMiddleware(auth string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/login.html" {
			ctx.Next()
			return
		}
		if auth == "" {
			ctx.Next()
			return
		}
		token := ctx.Query("token")
		if token == "" {
			token, _ = ctx.Cookie("token")
		}
		if token != auth {
			//	ctx.JSON(401, gin.H{"error": "认证失败"})
			ctx.Redirect(302, "/login.html")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

//go:embed login.html
var html string

func NewApp(docsDir string, auth string) *gin.Engine {
	router := gin.New()
	authFn := authMiddleware(auth)
	router.Use(gin.Recovery())
	router.POST("/api/docs", authFn, func(ctx *gin.Context) {
		// 获取目标文件夹下的所有.json文件
		files, err := findJSONFiles(docsDir)
		if err != nil {
			ctx.JSON(501, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, files)
	})
	defaultDocsDir := docsDir // 默认文档目录 ./docs/xxx docs/xxx
	router.POST("/api/upload", authFn, func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(501, gin.H{"error": err.Error()})
			return
		}
		fileName := file.Filename
		filePath := path.Join(defaultDocsDir, fileName)
		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			ctx.JSON(501, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "上传成功"})
	})
	router.POST("/api/login", func(ctx *gin.Context) {
		var data map[string]string
		ctx.BindJSON(&data)
		token := data["token"]

		if token != auth {
			ctx.JSON(401, gin.H{"error": "认证失败"})
			return
		}
		ctx.SetCookie("token", token, 3600*24*30, "/", "", false, true)
		ctx.JSON(200, gin.H{"message": "登录成功", "success": true})
	})
	router.DELETE("/api/delete", authFn, func(ctx *gin.Context) {
		data := map[string]string{}
		ctx.BindJSON(&data)
		url := data["url"]
		url = strings.Replace(url, "..", "", -1)
		url = strings.Replace(url, "/docs/", "", 1)
		filePath := path.Join(defaultDocsDir, url)
		if err := os.Remove(filePath); err != nil {
			ctx.JSON(501, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"message": "删除成功"})
	})
	router.GET("/*filepath", authFn, func(ctx *gin.Context) {

		file := ctx.Param("filepath")
		file = strings.Replace(file, "..", "", -1)

		if file == "/login.html" {
			ctx.Data(200, "text/html", []byte(html))
			return
		}

		if file == "/" {
			file = "/index.html"
		}

		if strings.HasPrefix(file, "/docs/") {
			file = strings.Replace(file, "/docs/", "", 1)
			realFilePath := path.Join(defaultDocsDir, file)
			ctx.File(realFilePath)
			return
		}
		ctype := mime.TypeByExtension(filepath.Ext(file))
		body, err := static.UI.ReadFile("ui" + file)
		ctx.Header("Content-Type", ctype)
		if err != nil {
			ctx.JSON(404, gin.H{"error": err.Error()})
			return
		}
		ctx.Data(200, ctype, body)
	})
	return router
}
func main() {
	var port string
	var docsDir string
	var auth string
	flag.StringVar(&port, "port", "8888", "端口号")
	flag.StringVar(&docsDir, "docs", "docs", "文档目录")
	flag.StringVar(&auth, "auth", "", "简单的认证,不填不认证")
	flag.Parse()

	if port == "" {
		log.Fatal("端口不能为空,启动时需添加参数,如: --port 8888 ")
	}
	app := NewApp(docsDir, auth)
	log.Println("listen on port:", port)
	app.Run(":" + port)
}
