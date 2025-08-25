package utils
import(
	"time"
	"github.com/golang-jwt/jwt/v5"
)
var JwtKey=[]byte("secret_key")//just for development
 func GenerateJWT(username string)(string,error){
	claims:=jwt.MapClaims{}
	claims["username"]=username
	claims["exp"]=time.Now().Add(time.Hour*24).Unix()
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return token.SignedString([]byte("secret_key"))
 }