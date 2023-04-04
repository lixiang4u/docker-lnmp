package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lixiang4u/docker-lnmp/util"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
)

type FileController struct {
}

func (x FileController) List(ctx *gin.Context) {
	var root = ctx.Query("path")

	if root == "" {
		root = util.AppDirectory()
	}
	var files []gin.H
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Join(root, d.Name()) != path {
			return nil
		}

		fi, err := d.Info()
		if err != nil {
			return err
		}

		log.Println("", d.Name(), ", ", path)
		files = append(files, gin.H{
			"path":   path,
			"name":   d.Name(),
			"is_dir": d.IsDir(),
			"mode":   fi.Mode(),
			"time":   fi.ModTime().Unix(),
			"size":   fi.Size(),
		})
		return nil
	})

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  nil,
		"data": files,
		"root": root,
	})
}
