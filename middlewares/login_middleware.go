package middlewares

import "github.com/gin-gonic/gin"

func CheckLoginned(c *gin.Context) {
	cookie, err := c.Cookie("access-token")
	if err == nil && cookie != "" {
		// 쿠키가 존재하면 로그인 상태로 간주하고 context에 저장
		c.Set("isLoggedIn", true)
	} else {
		// 쿠키가 없으면 로그아웃 상태로 간주
		c.Set("isLoggedIn", false)
	}

	c.Next() // 다음 핸들러 호출
}
