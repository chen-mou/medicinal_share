package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	redis2 "github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"medicinal_share/main/entity"
	"medicinal_share/tool"
	"medicinal_share/tool/db/elasticsearch"
	"medicinal_share/tool/db/mysql"
	"medicinal_share/tool/db/redis"
	"medicinal_share/tool/encrypt/md5"
	"strconv"
	"time"
)

const (
	userKeyPrefix     = "USER:"
	userLockKeyPrefix = "USER:LOCK:"
)

func GetByName(username string) *entity.User {
	res := &entity.User{}
	err := mysql.GetConnect().Where("username = ?", username).Take(res).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}
	res.Username = username
	return res
}

func getByIdFromDB(userId int64) *entity.User {
	var user entity.User
	err := mysql.GetConnect().Model(&entity.User{}).
		Joins("UserInfo").Preload("Role").
		Where("user.id = ?", userId).
		Take(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}
	return &user
}

func getByIdFromCache(userId int64) *entity.User {
	user := &entity.User{}
	err := redis.Get(userKeyPrefix+strconv.FormatInt(userId, 10), user)
	if err != nil {
		if errors.Is(err, redis2.Nil) {
			return nil
		}
		panic(err)
	}
	return user
}

func GetById(userId int64) *entity.User {
	id := strconv.FormatInt(userId, 10)
	key := userKeyPrefix + id
	lock := userLockKeyPrefix + id
	c := redis.NewCache(lock, key)
	return c.Get(&entity.User{}, func() any {
		return getByIdFromDB(userId)
	}).(*entity.User)
}

func Create(name, password string, tx *gorm.DB) *entity.User {
	id, err := tool.GetId("user")
	if err != nil {
		panic(err)
	}
	user := entity.User{
		Id:       id,
		Username: name,
		Password: password,
	}
	err = tx.Create(&user).Error
	if err != nil {
		panic(err)
	}
	return &user
}

func CreateData(userId int64, tx *gorm.DB) *entity.UserData {
	data := entity.UserData{
		UserId:   userId,
		Nickname: "新用户" + md5.Hash(time.Now().String())[25:],
		Avatar:   1,
	}
	err := tx.Create(&data).Error
	if err != nil {
		panic(err)
	}
	return &data
}

func GetDataByUserId(userId int64) *entity.UserData {
	data := entity.UserData{}
	err := mysql.GetConnect().Where("user_id = ?", userId).
		Joins("AvatarFile").
		Preload("AvatarFile.File").Take(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(err)
	}
	return &data
}

func UpdateData(data *entity.UserData, tx *gorm.DB) *entity.UserData {
	err := tx.Model(&entity.UserData{}).Where("user_id = ?", data.UserId).Updates(*data).Error
	if err != nil {
		panic(err)
	}
	if data.InfoId != nil {
		err = tx.Model(&entity.DoctorInfo{}).
			Where("user_id = ?", data.UserId).
			Update("info_id", data.InfoId).Error
		if err != nil {
			panic(err)
		}
	}
	return data
}

func UpdatePassword(userId int64, password string, tx *gorm.DB) {
	err := tx.Model(&entity.User{}).Where("id = ?", userId).Update("password", password).Error
	if err != nil {
		panic(err)
	}
}

type DoctorStatus uint

const (
	Online  = iota //在线
	Busy           //忙碌
	Offline        //下线
)

// UpdateDoctorStatus TODO:更新医生当前的状态
func UpdateDoctorStatus(userId int64, status DoctorStatus) {}

// GetBestDoctor 获取最佳匹配的医生
func GetBestDoctor(tags []int64, long float64, latit float64) int64 {
	tag := make([]int64, 0)
	db := mysql.GetConnect()
	err := db.Model(&entity.TagRelation{}).
		Select("relation_id").
		Where("relation_type = 'tag'").
		Joins("Tag", db.Where("Tag.id in (?) and Tag.type = 'Symptom'", tags)).Scan(&tags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	queryBody := map[string]any{
		"size":    1,
		"_source": "doctor_id",
		"query": map[string]any{"bool": map[string]any{
			"should": map[string]any{
				"terms": map[string][]int64{"tags": tag},
			},
			"filter": map[string]any{
				"term": map[string]DoctorStatus{
					"status": Online,
				},
			},
		}},
		"sort": []map[string]any{
			{
				"_geo_distance": map[string]any{
					"location": map[string]any{
						"lat": latit,
						"lon": long,
					},
					"order":         "asc",
					"distance_type": "plane",
				},
			},
			{"_score": map[string]string{
				"order": "asc",
			}},
		},
	}
	byt, _ := json.Marshal(queryBody)
	res := map[string]any{}
	elasticsearch.Get(&res,
		elasticsearch.GetClient().Search.WithBody(bytes.NewBuffer(byt)),
		elasticsearch.GetClient().Search.WithIndex("doctor_tag"),
	)
	fmt.Println(res)
	return 0
}

func CreateWorker(tx *gorm.DB, hospitalId, userId int64) {
	err := tx.Create(&entity.Worker{
		HospitalId: hospitalId,
		UserId:     userId,
	}).Error
	if err != nil {
		panic(err)
	}
}

// GetBestDoctorTest TODO: 测试方法
func GetBestDoctorTest() int64 {
	return 1234
}
