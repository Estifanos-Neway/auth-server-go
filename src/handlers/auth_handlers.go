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
	UsersByEmail struct {
		Id               string `json:"id"`
		Email            string `json:"email"`
		Full_Name        string `json:"fullName"`
		Avatar_Image_Url string `json:"avatarImageUrl"`
		Member_Since     string `json:"memberSince"`
		Password_Hash    string `json:"passwordHash"`
	} `graphql:"users_by_email(args:{user_email:$user_email})"`
}

func getUserByEmail(email string) (usersByEmailQuery, error) {
	query := usersByEmailQuery{}
	variables := map[string]interface{}{
		"user_email": email,
	}

	if err := gqClient.Query(context.Background(), &query, variables); err != nil {
		return query, err
	}
	return query, nil

}

func SignUpUser(c *gin.Context) {
	var signUpInput types.SignUpInput
	if err := c.BindJSON(&signUpInput); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid_Sign_In_Input"})
		return
	}
	existingUser, err := getUserByEmail(signUpInput.Email)
	if err != nil {
		log.Println(err)
		return
	}
	if existingUser.UsersByEmail.Email == signUpInput.Email {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Email_Already_Exist"})
		return
	}
	c.IndentedJSON(http.StatusOK, existingUser)
}

func SignInUser(c *gin.Context) {
	notImplemented(c)
}
