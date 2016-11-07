package handlers

import (
    "github.com/dgrijalva/jwt-go"
    "net/http"
    "time"
    "github.com/auth0/go-jwt-middleware"
    // "encoding/base64"
    // "os"
)

var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
    ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
        //get auth0 client secret secret from env variable
        // decoded, err := base64.URLEncoding.DecodeString(os.Getenv("AUTH0_CLIENT_SECRET"))
        // if err != nil {
        //     return nil, err
        // }
        // return decoded, nil
        return mySigningKey, nil
    },
    SigningMethod: jwt.SigningMethodHS256,
})

var mySigningKey = []byte("secret")

var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    /* Create the token */
    token := jwt.New(jwt.SigningMethodHS256)

    /* Set token claims */

    claims := make(jwt.MapClaims)
    claims["admin"] = true
    claims["name"] = "Ado Kukic"
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

    /* Sign the token with our secret */
    tokenString, _ := token.SignedString(mySigningKey)

    /* Finally, write the token to the browser window */
    w.Write([]byte(tokenString))
})
