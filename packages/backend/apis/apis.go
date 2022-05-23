package apis

import (
	db2 "backend/db"
	"backend/models"
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Msg  string      `json:"msg" example:"ok" default:"ok"`
	Code int         `json:"code" example:"0" format:"int" default:"0"`
	Data interface{} `json:"data"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &Response{
		Msg:  "ok",
		Code: 0,
		Data: data,
	})
}

func ErrorResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, &Response{
		Msg:  msg,
		Code: 4999,
		Data: nil,
	})
}

type GetLoginNonceResponseData struct {
	Nonce string `json:"nonce"`
}

type MetamaskLoginParameter struct {
	PublicAddress string `json:"publicAddress"`
	Signature     string `json:"signature"`
	Nonce         string `json:"nonce"`
}

type MetamaskLoginResponseData struct {
	AccessToken string `json:"accessToken"`
}

type GetProfileResponseData struct {
	Uid           uint   `json:"uid"`
	PublicAddress string `json:"publicAddress"`
	Username      string `json:"username"`
}

type ModifyUsernameParameter struct {
	Username string `json:"username"`
}

func GetLoginNonceAPI(c *gin.Context) {
	publicAddress, _ := c.GetQuery("publicAddress")

	if publicAddress == "" || !utils.IsAddressValid(publicAddress) {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	db, _ := db2.GetDB()
	var user models.User
	db.Where("public_address = ?", publicAddress).Limit(1).Find(&user)
	// user not in database
	if user.ID > 0 {
		c.JSON(http.StatusOK, GetLoginNonceResponseData{Nonce: user.Nonce})
		return
	}

	c.JSON(http.StatusOK, GetLoginNonceResponseData{Nonce: utils.RandStringRunes(6)})
}

func MetamaskLoginAPI(c *gin.Context) {
	params := MetamaskLoginParameter{}
	err := c.BindJSON(&params)

	if err != nil {
		c.JSON(http.StatusBadRequest, "Request should have signature and publicAddress")
		return
	}

	isValid := utils.VerifySig(params.PublicAddress, params.Signature, []byte("I am signing my one-time nonce: "+params.Nonce))
	fmt.Printf("is valid: %v", isValid)
	if !isValid {
		c.JSON(http.StatusUnauthorized, "login failed")
		return
	}

	var user models.User
	db, _ := db2.GetDB()

	db.Where("public_address = ?", params.PublicAddress).Limit(1).Find(&user)

	if user.ID > 0 {
		if user.Nonce != params.Nonce {
			c.JSON(http.StatusUnauthorized, "login failed")
			return
		}

		user.Nonce = utils.RandStringRunes(6)
		db.Save(&user)

		c.JSON(http.StatusOK, MetamaskLoginResponseData{
			AccessToken: utils.GenerateToken(utils.CustomClaims{
				Uid:           user.ID,
				PublicAddress: user.PublicAddress,
			}),
		})
		return
	}

	newUser := models.User{PublicAddress: params.PublicAddress, Nonce: utils.RandStringRunes(8)}
	err = db.Create(&newUser).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, "login failed")
		return
	}

	c.JSON(http.StatusOK, MetamaskLoginResponseData{
		AccessToken: utils.GenerateToken(utils.CustomClaims{
			Uid:           newUser.ID,
			PublicAddress: user.PublicAddress,
		}),
	})
}

func GetProfileAPI(c *gin.Context) {
	publicAddress := c.GetString("publicAddress")

	db, _ := db2.GetDB()
	var user models.User
	db.Where("public_address = ?", publicAddress).Limit(1).Find(&user)

	SuccessResponse(c, GetProfileResponseData{
		Uid:           user.ID,
		Username:      user.Username,
		PublicAddress: user.PublicAddress,
	})
}

func ModifyUsernameAPI(c *gin.Context) {
	publicAddress := c.GetString("publicAddress")
	params := ModifyUsernameParameter{}

	err := c.BindJSON(&params)
	if err != nil || params.Username == "" {
		c.JSON(http.StatusBadRequest, "bad request")
		return
	}

	db, _ := db2.GetDB()
	var user models.User
	db.Where("public_address = ?", publicAddress).Limit(1).Find(&user)

	user.Username = params.Username
	db.Save(&user)

	SuccessResponse(c, GetProfileResponseData{
		Uid:           user.ID,
		Username:      user.Username,
		PublicAddress: user.PublicAddress,
	})
}
