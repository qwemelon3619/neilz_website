package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"neilz.space/web/models"
	"neilz.space/web/utils"
)

func RegisterService(ID string, password string) error {
	// 일반적으로 cost는 10을 많이 사용한다. O(2^cost)가 걸리는 연산이라, cost를 너무 크게 잡지는 말자.
	// cost를 적게 쓰면 가볍지만, 보안에 취약하여 브루트 포스 공격에 위험이 있고,
	// cost가 크면, 무겁지만 보안이 강력하다.
	// 패스워드 bcrypt 암호화(cost 10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// db에 저장..
	err = models.SaveUserEmailAndPassword(ID, string(hashedPassword))
	return err
}

func LoginService(userID, userPassword, hostInfo string) (string, string, error) {

	// find user by user_id

	usr, err := models.FindUserByUserID(userID)
	if err != nil {
		return "", "", errors.New("input error")
	}

	// compare password
	if !utils.IsPasswordCorrect(usr.UserPassword, userPassword) {
		return "", "", errors.New("incorrect password")
	}

	// generate Access Token
	accessToken, err := utils.GenerateAccessToken(usr.UserUUID)
	if err != nil {
		return "", "", err
	}

	// generate Refresh Token
	refreshToken, err := utils.GenerateRefreshToken(usr.UserUUID)
	if err != nil {
		return "", "", err
	}

	// save refreshToken to DB
	err = models.SaveRefreshToken(usr.UserUUID, refreshToken, hostInfo)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func LogoutService(hostInfo string) error {
	// delete user refreshToken
	err := models.RemoveRefreshTokenDB(hostInfo)
	if err != nil {
		return err
	}
	return nil
}
