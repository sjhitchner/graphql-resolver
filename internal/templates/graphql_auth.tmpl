package graphql

import (
	"context"
	//"encoding/base64"
	//"encoding/json"
	"log"
	//"net"
	"net/http"
	"strings"

	//jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	ggg "github.com/sjhitchner/slack-clone/backend/interfaces/context"
	//gcontext "github.com/OscarYuen/go-graphql-starter/context"
	//"github.com/OscarYuen/go-graphql-starter/model"
	//"github.com/OscarYuen/go-graphql-starter/service"
)

// https://github.com/OscarYuen/go-graphql-starter/blob/master/handler/auth.go
// TODO constants
// TODO context variables, use integers for speed?
const (
	HeaderAuthorization = "Authorization"
	HeaderBearer        = "Bearer"
)

func WrapAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userId, isAuthorized, err := ValidateAuthorizationHeader(ctx, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Error validating authorization header %v\n", err)
			return
		}

		if !isAuthorized {
			w.WriteHeader(http.StatusUnauthorized)
		}

		ctx = ggg.SetCurrentUserId(ctx, userId)
		//ctx = ggg.SetIsAuthorized(ctx, isAuthorized)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateAuthorizationHeader(ctx context.Context, r *http.Request) (int64, bool, error) {
	token, _ := ParseAuthorizationHeader(r)
	//if err != nil {
	//	return 0, false, err
	//}

	log.Println("TOKEN " + token)

	// TODO Parse Token
	return 1, true, nil
}

func ParseAuthorizationHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get(HeaderAuthorization)
	if authHeader == "" {
		return "", errors.Errorf("No %s provided", HeaderAuthorization)
	}

	auth := strings.SplitN(authHeader, " ", 2)
	if len(auth) != 2 || auth[0] != HeaderBearer {
		return "", errors.Errorf("Invalid %s '%s'", HeaderAuthorization, authHeader)
	}

	return auth[1], nil
}

/*
		token, err := validateBearerAuthHeader(ctx, r)
		if err == nil {
			isAuthorized = true
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userIdByte, _ := base64.StdEncoding.DecodeString(claims["id"].(string))
				userId = string(userIdByte[:])
			} else {
				log.Println(err)
			}
		}
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println(w, "Requester ip: %q is not IP:port", r.RemoteAddr)
		}

		ctx = context.WithValue(ctx, "user_id", &userId)
		ctx = context.WithValue(ctx, "requester_ip", &ip)
		ctx = context.WithValue(ctx, "is_authorized", isAuthorized)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
*/

/*
func writeResponse(w http.ResponseWriter, response interface{}, code int) {
	jsonResponse, _ := json.Marshal(response)
	w.WriteHeader(code)
	w.Write(jsonResponse)
}

func validateBasicAuthHeader(r *http.Request) (*model.UserCredentials, error) {
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		return nil, errors.New(gcontext.CredentialsError)
	}
	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return nil, errors.New(gcontext.CredentialsError)
	}
	userCredentials := model.UserCredentials{
		Email:    pair[0],
		Password: pair[1],
	}
	return &userCredentials, nil
}

func validateBearerAuthHeader(ctx context.Context, r *http.Request) (*jwt.Token, error) {
	var tokenString string
	keys, ok := r.URL.Query()["at"]
	if !ok || len(keys) < 1 {
		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || auth[0] != "Bearer" {
			return nil, errors.New(gcontext.CredentialsError)
		}
		tokenString = auth[1]
	} else {
		tokenString = keys[0]
	}
	token, err := ctx.Value("authService").(*service.AuthService).ValidateJWT(&tokenString)
	return token, err
}
*/
