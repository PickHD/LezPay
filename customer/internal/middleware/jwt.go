package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/PickHD/LezPay/customer/internal/helper"
	"github.com/PickHD/LezPay/customer/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type (
	// DecodePayloadData consists decoded payload data
	DecodePayloadData struct {
		UserID uint64 `json:"user_id"`
		Email  string `json:"email"`
	}
)

const (
	payloadUserID  string = "user_id"
	payloadEmail   string = "email"
	payloadExpires string = "exp"

	JWTExpire time.Duration = time.Duration(7) * time.Hour
)

// ValidateJWTMiddleware responsible to validating jwt in header each request
func ValidateJWTMiddleware(ctx *fiber.Ctx) error {
	// validate JWT coming from request, if valid decode into a struct
	decodedPayload, err := validate(ctx)
	if err != nil {
		return helper.NewResponses[any](ctx, fiber.StatusUnauthorized, fmt.Sprintf("Unauthorized access, reason : %s", err.Error()), err, nil, nil)
	}

	// pass decoded payload into ctx.Locals()
	ctx.Locals(model.KeyJWTValidAccess, decodedPayload)

	// going to next handler..
	return ctx.Next()
}

// validate will checking validity of signed JWT token from request in
func validate(ctx *fiber.Ctx) (DecodePayloadData, error) {
	header := ctx.Get("Authorization", "")
	if !strings.Contains(header, "Bearer") {
		return DecodePayloadData{}, model.NewError(model.NotFound, "Token not found")
	}

	getToken := strings.Replace(header, "Bearer ", "", -1)
	validToken, err := jwt.Parse(getToken, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, model.NewError(model.Validation, "Invalid token")
		}
		return []byte(helper.GetEnvString("JWT_SECRET")), nil
	})
	if err != nil {
		return DecodePayloadData{}, err
	}

	claims := validToken.Claims.(jwt.MapClaims)

	// Check is token expired or not
	if expInt, ok := claims[payloadExpires].(float64); ok {
		now := time.Now().Unix()
		if now > int64(expInt) {
			return DecodePayloadData{}, model.NewError(model.Validation, "Token expired")
		}
	} else {
		return DecodePayloadData{}, model.NewError(model.Type, "type assertion payload exp error")
	}

	decodePayload, err := insertPayload(claims)
	if err != nil {
		return DecodePayloadData{}, err
	}

	return decodePayload, nil
}

// Extract will extracting payload data from ctx.Locals
func Extract(data interface{}) (DecodePayloadData, error) {
	extractData, ok := data.(DecodePayloadData)
	if !ok {
		return DecodePayloadData{}, model.NewError(model.Type, "type assertion extract data error")
	}

	return extractData, nil
}

// insertPayload will inserting data from decoded payload into defined struct
func insertPayload(claims jwt.MapClaims) (DecodePayloadData, error) {
	decodePayloadData := DecodePayloadData{}

	if userID, ok := claims[payloadUserID].(float64); ok {
		decodePayloadData.UserID = uint64(userID)
	} else {
		return DecodePayloadData{}, model.NewError(model.Type, "type assertion user_id error")
	}

	if userEmail, ok := claims[payloadEmail].(string); ok {
		decodePayloadData.Email = userEmail
	} else {
		return DecodePayloadData{}, model.NewError(model.Type, "type assertion email error")
	}

	return decodePayloadData, nil
}
