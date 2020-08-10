package models

import "App_Conf/dao"

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Age      int    `json:"age"`
	Mobile   string `json:"mobile"`
	Sex      string `json:"sex"`
	Address  string `json:"address"`
}

//增加一个用户
func CreateAUser(user *User) (err error) {
	err = dao.DB.Create(&user).Error
	return
}

func GetUserList() (userList *[]User, err error) {
	var list []User
	if err = dao.DB.Find(&list).Error; err != nil {
		return nil, err
	}
	return &list, nil
}

func GetAUser(id string) (user *User, err error) {
	user = new(User)
	dao.DB.Find("id=?", id).First(user)
	return
}

func UpdateAUser(user *User) (err error) {
	err = dao.DB.Save(user).Error
	return
}

func DeleteAUser(id string) (err error) {
	err = dao.DB.Where("id=?", id).Delete(&User{}).Error
	return
}
