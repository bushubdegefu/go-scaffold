package temps

import (
	"os"
	"text/template"
)

func KeyCloakFrame() {
	// ####################################################
	//  rabbit template
	kc_tmpl, err := template.New("RenderData").Parse(keycloakMiddleware)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("nosqlconn", os.ModePerm)
	if err != nil {
		panic(err)
	}

	nosqlconn_file, err := os.Create("uitils/keycloak.go")
	if err != nil {
		panic(err)
	}
	defer nosqlconn_file.Close()

	err = rpc_tmpl.Execute(nosqlconn_file, RenderData)
	if err != nil {
		panic(err)
	}
}

var keycloakMiddleware = `
package keycloakmiddleware

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v4"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"

	"github.com/coreos/go-oidc/v3/oidc"
)

type Claims struct {
	jwt.RegisteredClaims
	Username   string {{.BackTick}}json:"preferred_username"{{.BackTick}}
	Name       string {{.BackTick}}json:"name"{{.BackTick}}
	GivenName  string {{.BackTick}}json:"given_name"{{.BackTick}}
	FamilyName string {{.BackTick}}json:"family_name"{{.BackTick}}
}

func middlewreuseageexample() {
	validator, err := NewKeycloakJWTValidator("http://localhost:8080/realms/myrealm", "myclient")
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	profile := app.Group("profile", keyauth.New(keyauth.Config{
		Validator: validator,
	}))
	profile.Get("name", func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(*Claims)
		return c.SendString(claims.Username)
	})

	app.Listen(":3000")
}

func NewKeycloakJWTValidator(issuerUrl, clientId string) (func(*fiber.Ctx, string) (bool, error), error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, issuerUrl)
	if err != nil {
		return nil, err
	}
	verifier := provider.Verifier(&oidc.Config{
		ClientID: clientId,
	})
	return func(c *fiber.Ctx, key string) (bool, error) {
		var ctx = c.UserContext()
		_, err := verifier.Verify(ctx, key)
		if err != nil {
			return false, err
		}
		token, _ := jwt.ParseWithClaims(key, &Claims{},
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v",
						token.Header["alg"])
				}
				return key, nil
			})
		c.Locals("claims", token.Claims)
		return true, nil
	}, nil
}
`
