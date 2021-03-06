package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pegasus-cloud/iam_client/iam"
	"github.com/pegasus-cloud/iam_client/protos"
	"github.com/pegasus-cloud/iam_client/utility"
)

func createUser(c *gin.Context) {
	userInfo := &user{}
	if err := c.ShouldBindWith(userInfo, binding.JSON); err != nil {
		utility.ResponseWithType(c, http.StatusBadRequest, &utility.ErrResponse{
			Message: err.Error(),
		})
		return
	}
	if userInfo.Force {
		if output, err := iam.GetUser(userInfo.UserID); err != nil {
			utility.ResponseWithType(c, http.StatusInternalServerError, &utility.ErrResponse{
				Message: databaseErrMsg,
			})
			return
		} else if output.ID != "" {
			utility.ResponseWithType(c, http.StatusBadRequest, &utility.ErrResponse{
				Message: userExistErrMsg,
			})
			return
		}
	} else {
		for true {
			userInfo.UserID = utility.GetRandNumeric(16)
			if output, err := iam.GetUser(userInfo.UserID); err != nil {
				utility.ResponseWithType(c, http.StatusInternalServerError, &utility.ErrResponse{
					Message: databaseErrMsg,
				})
				return
			} else if output.ID == "" {
				break
			}
		}
	}
	if err := iam.CreateUser(&protos.UserInfo{
		ID:           userInfo.UserID,
		DisplayName:  userInfo.DisplayName,
		PasswordHash: userInfo.Password,
		Description:  userInfo.Description,
		Extra:        userInfo.Extra,
	}); err != nil {
		utility.ResponseWithType(c, http.StatusInternalServerError, &utility.ErrResponse{
			Message: err.Error(),
		})
		return
	}
	getUserOutput, err := iam.GetUser(userInfo.UserID)
	if err != nil {
		utility.ResponseWithType(c, http.StatusInternalServerError, &utility.ErrResponse{
			Message: err.Error(),
		})
		return
	}
	utility.ResponseWithType(c, http.StatusCreated, user{
		UserID:      getUserOutput.ID,
		DisplayName: getUserOutput.DisplayName,
		Description: getUserOutput.Description,
		Extra:       getUserOutput.Extra,
		CreatedAt:   getUserOutput.CreatedAt,
		UpdatedAt:   getUserOutput.UpdatedAt,
	})
}
