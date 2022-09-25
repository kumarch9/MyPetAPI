package signout

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignOutAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		// claims := &tokens.Claims{}
		//key := tokens.Jwt_Key

		_, errInReadCookie_user := c.Cookie("user_token")
		if errInReadCookie_user != nil {
			_, errInReadCookie := c.Cookie("adm_token")
			if errInReadCookie != nil {
				c.AbortWithStatus(http.StatusNotFound)
				c.Redirect(http.StatusTemporaryRedirect, "/home")
				return
			}
			c.SetSameSite(http.SameSiteLaxMode)
			c.SetCookie("adm_token", "", -1, "", "", true, true)
			c.Redirect(http.StatusTemporaryRedirect, "/home")
			return
		}
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("user_token", "", -1, "", "", true, true)
		c.Redirect(http.StatusTemporaryRedirect, "/home")

	}

}
