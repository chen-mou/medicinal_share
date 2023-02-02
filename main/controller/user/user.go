package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/main/middleware"
	"medicinal_share/main/service/file"
	"medicinal_share/main/service/user"
	"medicinal_share/tool"
	"medicinal_share/tool/encrypt/jwtutil"
	"strconv"
	"strings"
)

func Login(ctx *gin.Context) {
	usr := entity.User{}
	ctx.Bind(&usr)
	u := user.Login(usr.Username, usr.Password)
	token, _ := jwtutil.GetToken(map[string]string{
		"id": strconv.FormatInt(u.Id, 10),
	})
	ctx.AbortWithStatusJSON(200, gin.H{
		"code":  0,
		"data":  u,
		"token": token,
	})
}

func Register(ctx *gin.Context) {
	usr := entity.User{}
	ctx.Bind(&usr)
	u := user.Register(usr.Username, usr.Password)
	token, _ := jwtutil.GetToken(map[string]string{
		"id": strconv.FormatInt(u.Id, 10),
	})
	ctx.AbortWithStatusJSON(200, gin.H{
		"code":  0,
		"data":  u,
		"token": token,
	})
}

func GetUserData(ctx *gin.Context) {
	usr := tool.GetNowUser(ctx)
	data := user.GetData(usr.Id)
	ctx.AbortWithStatusJSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}

func UploadAvatar(ctx *gin.Context) {
	f, err := ctx.FormFile("file")
	user := tool.GetNowUser(ctx)
	if err != nil {
		panic(err)
	}
	if f.Size >= 5<<20 {
		panic(middleware.NewCustomErr(middleware.ERROR, "文件过大"))
	}
	suffix := strings.Split(f.Filename, ".")[1]
	if _, ok := file.Suffix[suffix]; !ok {
		panic(middleware.NewCustomErr(middleware.ERROR, "文件类型有误"))
	}
	file.Upload(f, "avatar", user.Id, func(i int64, db *gorm.DB) error {
		return db.Model(&entity.UserData{}).Where("user_id = ?", user.Id).Update("avatar", i).Error
	})
	ctx.AbortWithStatusJSON(200, gin.H{
		"data": 0,
	})
}

func UpdateInfo(ctx *gin.Context) {

}
