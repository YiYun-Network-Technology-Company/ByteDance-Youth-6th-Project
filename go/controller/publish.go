package controller

import (
	"SkyLine/entity"
	"SkyLine/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
)

type VideoListResponse struct {
	entity.Response
	VideoList []entity.DouyinVideo `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")

	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, entity.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, entity.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	fileName := filepath.Base(data.Filename)
	user := usersLoginInfo[token]
	finalName := fmt.Sprintf("%d_%s", user.Id, fileName)
	saveFile := filepath.Join("./public/", finalName)
	// 使用os.Stat函数检查文件是否存在
	//_, err := os.Stat(saveFile)
	//if err == nil {
	//	fmt.Println("文件存在")
	//} else if os.IsNotExist(err) {
	//	fmt.Println("文件不存在")
	//} else {
	//	fmt.Println("发生了其他错误：", err)
	//}
	fileContent, err := ioutil.ReadFile(saveFile)
	if err != nil {
		fmt.Println("读取文件内容时发生错误：", err)
		return
	}
	// 调用上传函数
	err = service.UploadFile(fileName, fileContent, 2)
	if err != nil {
		fmt.Println("上传文件时发生错误：", err)
	} else {
		fmt.Println("文件上传成功！")
	}
}

// @Summary  获取某一用户的所发布的搜游视频
// @Description  这个接口，会根据用户id去查询该用户发布的所有的视频
// @Tags         视频相关接口
// @Param        userid  query  int64  ture  "用户的userid"
// @Param        Token  query  string  ture  "该参数只有在用户登录状态下进行设置"
// @Router       /publish/list [get]
func PublishList(c *gin.Context) {
	publishListRequest := &entity.PublishListRequest{c.GetInt64("user_id"), c.GetString("token")}
	video, err := service.SelectVideoListByUserId(publishListRequest)
	if err != nil {
		fmt.Printf("视频获取出错:%v\n", err)
		c.JSON(http.StatusInternalServerError, entity.FeedResponse{
			Response:  entity.Response{StatusCode: 1, StatusMsg: "获取视频列表出现错误"},
			VideoList: nil,
			NextTime:  time.Now().Unix(),
		})
	}
	//将获取的video输出，方便测试
	//fmt.Printf("%#v", video)
	douyinVideos := make([]entity.DouyinVideo, len(video))
	for i := range video {
		author := entity.Author{
			Avatar:          video[i].UserDetail.Avatar,
			BackgroundImage: video[i].UserDetail.BackgroundImage,
			FavoriteCount:   video[i].UserDetail.FavoriteCount,
			FollowCount:     video[i].UserDetail.FollowCount,
			FollowerCount:   video[i].UserDetail.FollowerCount,
			ID:              video[i].UserDetail.ID,
			IsFollow:        video[i].UserDetail.IsFollow,
			Name:            video[i].UserDetail.Name,
			Signature:       video[i].UserDetail.Signature,
			TotalFavorited:  video[i].UserDetail.TotalFavorited,
			WorkCount:       video[i].UserDetail.WorkCount,
		}
		douyinVideo := &entity.DouyinVideo{
			Author:        author,
			CommentCount:  video[i].CommentCount,
			CoverURL:      video[i].CoverUrl,
			FavoriteCount: video[i].FavoriteCount,
			ID:            video[i].VideoId,
			IsFavorite:    video[i].IsFollow,
			PlayURL:       video[i].PlayUrl,
			Title:         video[i].Title,
		}
		douyinVideos[i] = *douyinVideo
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: entity.Response{
			StatusCode: 0,
			StatusMsg:  "查询视频列表成功",
		},
		VideoList: douyinVideos,
	})
}
