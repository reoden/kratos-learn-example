package data

import (
	"context"
	"fmt"

	"crypto/sha512"
	"time"
	"user/internal/biz"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// 定义数据表结构体
type User struct {
    ID          int64      `gorm:"primarykey"`
    Mobile      string     `gorm:"index:idx_mobile;unique;type:varchar(11) comment '手机号码，用户唯一标识';not null"`
    Password    string     `gorm:"type:varchar(100);not null "` // 用户密码的保存需要注意是否加密
    NickName    string     `gorm:"type:varchar(25) comment '用户昵称'"`
    Birthday    *time.Time `gorm:"type:datetime comment '出生日日期'"`
    Gender      string     `gorm:"column:gender;default:male;type:varchar(16) comment 'female:女,male:男'"`
    Role        int        `gorm:"column:role;default:1;type:int comment '1:普通用户，2:管理员'"`
    CreatedAt   time.Time  `gorm:"column:add_time"`
    UpdatedAt   time.Time  `gorm:"column:update_time"`
    DeletedAt   gorm.DeletedAt
    IsDeletedAt bool
}

type userRepo struct {
	uc *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		uc: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) CreateUser(ctx context.Context, u *biz.User) (*biz.User, error) {
	var user User

	result := r.uc.db.Where(&biz.User{Mobile: u.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	user.Mobile = u.Mobile
	user.NickName = u.NickName
	user.Password = encrypt(u.Password)

	res := r.uc.db.Create(&user)

	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, res.Error.Error())
	}

	return &biz.User{
		ID: user.ID,
		NickName: user.NickName,
		Password: user.Password,
		Mobile: user.Mobile,
		Gender: user.Gender,
		Role: user.Role,
	}, nil
}

func encrypt(psd string) string {
	options := &password.Options{
		SaltLen: 16,
		Iterations: 10000,
		KeyLen: 32,
		HashFunction: sha512.New,
	}

	salt, encodedPwd := password.Encode(psd, options)

	return fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
}
