package auth

import (
	"context"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kinde-oss/kinde-go/oauth2/authorization_code"
)

/*

to use it u have to pass through an Auth Provider:

	privateGroup := router.Group("/private")
	err = auth.UseKindeAuth(
		privateGroup,
		os.Getenv("KINDE_DOMAIN"),
		os.Getenv("KINDE_CLIENT_ID"),
		os.Getenv("KINDE_CLIENT_SECRET"),
		os.Getenv("KINDE_REDIRECT_URI"),
	)

*/

func UseKindeAuth(router *gin.RouterGroup, kindeDomain, clientID, clientSecret, baseRedirectURL string, options ...authorization_code.Option) error {
	basePath := router.BasePath()
	if basePath == "/" {
		basePath = ""
	}
	redirectURI := fmt.Sprintf("%s%s%s", baseRedirectURL, basePath, "/kinde/callback")

	router.Use(func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		sessionStorage := &SessionStorage{session: session}

		reqOptions := append(options,
			authorization_code.WithSessionHooks(sessionStorage),
			authorization_code.WithTokenValidation(true),
		)
		kindeClient, err := authorization_code.NewAuthorizationCodeFlow(kindeDomain,
			clientID,
			clientSecret,
			redirectURI,
			reqOptions...,
		)

		if err != nil {
			fmt.Printf("Error creating Kinde client: %v", err)
			ctx.String(500, "Error creating Kinde client")
			ctx.Abort()
			return
		}

		ctx.Set("kinde_client", kindeClient)
	})

	router.GET("/kinde/callback", func(ctx *gin.Context) {
		if client, ok := ctx.Get("kinde_client"); ok {
			if kindeClient, ok := client.(authorization_code.IAuthorizationCodeFlow); ok {
				err := kindeClient.ExchangeCode(context.Background(), ctx.Query("code"), ctx.Query("state"))
				if err != nil {
					ctx.AbortWithError(500, err)
					return
				}
				ctx.Redirect(302, "/")
				return
			}
		}
		ctx.AbortWithError(500, fmt.Errorf("kinde client not found"))
	})

	router.Use(func(ctx *gin.Context) {

		if client, ok := ctx.Get("kinde_client"); ok {
			if kindeClient, ok := client.(authorization_code.IAuthorizationCodeFlow); ok {

				if isAuthenticated, _ := kindeClient.IsAuthenticated(context.Background()); !isAuthenticated {
					// Check for invitation_code query parameter
					invitationCode := ctx.Query("invitation_code")
					authURL := kindeClient.GetAuthURLWithInvitation(invitationCode)
					ctx.Redirect(302, authURL)
					ctx.Abort()
				}
				return
			}
		}

		ctx.AbortWithError(401, fmt.Errorf("unauthorized"))

	})

	return nil
}
