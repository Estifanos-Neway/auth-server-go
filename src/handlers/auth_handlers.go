package handlers

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/estifanos-neway/auth-server-with-go/src/types"
	"github.com/gin-gonic/gin"
	"github.com/hasura/go-graphql-client"
)

type myRoundTripper struct {
	rt http.RoundTripper
}

func (mrt myRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("x-hasura-admin-secret", os.Getenv("X_HASURA_ADMIN_SECRETE"))
	return mrt.rt.RoundTrip(req)
}

var httpClient *http.Client = &http.Client{
	Transport: myRoundTripper{rt: http.DefaultTransport},
}

var gqClient = graphql.NewClient("https://funny-starfish-21.hasura.app/v1/graphql", httpClient)

func notImplemented(c *gin.Context) {
	c.IndentedJSON(http.StatusNotImplemented, gin.H{"message": "Not_Implemented_Yet"})
}

type usersByEmailQuery struct {
	UsersByEmail types.User `graphql:"users_by_email(args:{user_email:$user_email})"`
}

func getUserByEmail(email string) (types.User, error) {
	query := usersByEmailQuery{}
	variables := map[string]interface{}{
		"user_email": email,
	}

	if err := gqClient.Query(context.Background(), &query, variables); err != nil {
		return query.UsersByEmail, err
	}

	return query.UsersByEmail, nil
}

type insertUsersMutation struct {
	InsertUsers struct {
		Returning []struct {
			Id               string `json:"id"`
			Email            string `json:"email"`
			Full_Name        string `json:"fullName"`
			Avatar_Image_Url string `json:"avatarImageUrl"`
			Member_Since     string `json:"memberSince"`
		}
	} `graphql:"insert_users(objects:{email:$email,full_name:$full_name,password_hash:$password_hash })"`
}

func insertUser(signUpInput types.SignUpInput) (insertUsersMutation, error) {
	mutation := insertUsersMutation{}
	variables := map[string]interface{}{
		"email":         signUpInput.Email,
		"full_name":     signUpInput.FullName,
		"password_hash": signUpInput.PasswordHash,
	}

	if err := gqClient.Mutate(context.Background(), &mutation, variables); err != nil {
		return mutation, err
	}

	return mutation, nil
}

func SignUpUser(c *gin.Context) {
	var signUpInput types.SignUpInput
	if err := c.BindJSON(&signUpInput); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid_Sign_Up_Input"})
		return
	} else if err := signUpInput.IsValidSignUpInput(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	existingUser, err := getUserByEmail(signUpInput.Email)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal_Error"})
		return
	}
	if existingUser.Email == signUpInput.Email {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Email_Already_Exist"})
		return
	}
	newUser, err := insertUser(signUpInput)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal_Error"})
		return
	}
	c.IndentedJSON(http.StatusOK, newUser.InsertUsers.Returning[0])
}

func SignInUser(c *gin.Context) {
	var signInInput types.SignInInput
	if err := c.BindJSON(&signInInput); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid_Sign_In_Input"})
		return
	} else if err := signInInput.IsValidSignInInput(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user, err := getUserByEmail(signInInput.Email)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal_Error"})
		return
	}
	if user.Email != signInInput.Email || user.Password_Hash != signInInput.PasswordHash {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User_Not_Found"})
		return
	}
	user.Password_Hash = ""
	authTokens := types.AuthTokens{
		AccessToken:  "",
		RefreshToken: "",
	}
	userLogin := types.UserLogin{
		User:   user,
		Tokens: authTokens,
	}
	c.IndentedJSON(http.StatusOK, userLogin)
}
