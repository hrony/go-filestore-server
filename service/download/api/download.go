package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-filestore-server/common"
	"go-filestore-server/config"
	"go-filestore-server/database/ceph"
	"go-filestore-server/database/oss"
	dbcli "go-filestore-server/service/dbproxy/client"
	"log"
	"net/http"
	"strings"
)

// DownloadURLHandler : 生成文件的下载地址
func DownloadURLHandler(c *gin.Context) {
	filehash := c.Request.FormValue("filehash")
	dbResp, err := dbcli.GetFileMeta(filehash)
	fmt.Println(filehash, dbResp, err)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusServerError,
			"msg":  "server error",
		})
		return
	}

	tblFile := dbcli.ToTableFile(dbResp.Data)

	// TODO: 判断文件存在OSS，还是Ceph，还是在本地
	if strings.HasPrefix(tblFile.FileAddr.String, config.DefaultConfig.TempLocalRootDir) ||
		strings.HasPrefix(tblFile.FileAddr.String, config.DefaultConfig.CephRootDir) {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")
		tmpURL := fmt.Sprintf("http://%s/file/download?filehash=%s&username=%s&token=%s",
			c.Request.Host, filehash, username, token)
		c.Data(http.StatusOK, "application/octet-stream", []byte(tmpURL))
	} else if strings.HasPrefix(tblFile.FileAddr.String, config.DefaultConfig.OSSRootDir) {
		// oss下载url
		signedURL := oss.DownloadURL(tblFile.FileAddr.String)
		log.Println(tblFile.FileAddr.String)
		c.Data(http.StatusOK, "application/octet-stream", []byte(signedURL))
	}
}

// DownloadHandler : 文件下载接口
func DownloadHandler(c *gin.Context) {
	fsha1 := c.Request.FormValue("filehash")
	username := c.Request.FormValue("username")
	// TODO: 处理异常情况
	fResp, ferr := dbcli.GetFileMeta(fsha1)
	ufResp, uferr := dbcli.QueryUserFileMeta(username, fsha1)
	fmt.Println(fsha1, username, fResp, ferr)
	fmt.Println(ufResp, uferr)
	if ferr != nil || uferr != nil || !fResp.Suc || !ufResp.Suc {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusServerError,
			"msg":  "server error",
		})
		return
	}
	uniqFile := dbcli.ToTableFile(fResp.Data)
	userFile := dbcli.ToTableUserFile(ufResp.Data)

	if strings.HasPrefix(uniqFile.FileAddr.String, config.DefaultConfig.TempLocalRootDir) {
		// 本地文件， 直接下载
		c.FileAttachment(uniqFile.FileAddr.String, userFile.FileName)
	} else if strings.HasPrefix(uniqFile.FileAddr.String, config.DefaultConfig.CephRootDir) {
		// ceph中的文件，通过ceph api先下载
		bucket := ceph.GetCephBucket("userfile")
		data, _ := bucket.Get(uniqFile.FileAddr.String)
		//	c.Header("content-type", "application/octect-stream")
		c.Header("content-disposition", "attachment; filename=\""+userFile.FileName+"\"")
		c.Data(http.StatusOK, "application/octect-stream", data)
	}
}
