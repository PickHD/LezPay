package requester

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PickHD/LezPay/auth/internal/config"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// SendMailRequest consist request send mail
	SendMailRequest struct {
		Subject          string   `json:"subject"`
		To               []string `json:"to"`
		Cc               string   `json:"cc"`
		CcTitle          string   `json:"cc_title"`
		Username         string   `json:"username"`
		Link             string   `json:"link"`
		UserType         string   `json:"user_type"`
		VerificationType string   `json:"verification_type"`
	}

	// SendMailResponse consist response from send mail
	SendMailResponse struct {
		Message string `json:"message"`
	}

	// UtilityRequester is an interface that has all the function to be implemented inside utility requester
	UtilityRequester interface {
		SendMail(ctx context.Context, req *SendMailRequest) (*SendMailResponse, error)
	}

	// UtilityRequesterImpl is an app utility struct that consists of all the dependencies needed for utility requester
	UtilityRequesterImpl struct {
		Context    context.Context
		Config     *config.Configuration
		Logger     *logrus.Logger
		Tracer     *trace.TracerProvider
		HTTPClient *http.Client
	}
)

// NewUtilityRequester return new instances utility requester
func NewUtilityRequester(ctx context.Context, config *config.Configuration, logger *logrus.Logger, tracer *trace.TracerProvider, httpCli *http.Client) *UtilityRequesterImpl {
	return &UtilityRequesterImpl{
		Context:    ctx,
		Config:     config,
		Logger:     logger,
		Tracer:     tracer,
		HTTPClient: httpCli,
	}
}

func (ur *UtilityRequesterImpl) SendMail(ctx context.Context, req *SendMailRequest) (*SendMailResponse, error) {
	tr := ur.Tracer.Tracer("Auth-SendMail Requester")
	_, span := tr.Start(ctx, "Start SendMail")
	defer span.End()

	jsonData, err := json.Marshal(req)
	if err != nil {
		ur.Logger.Error("UtilityRequesterImpl.SendMail json.Marshal ERROR ", err)

		return nil, err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", ur.Config.Service.UtilityURL, "send"), bytes.NewBuffer(jsonData))
	if err != nil {
		ur.Logger.Error("UtilityRequesterImpl.SendMail http.NewRequest ERROR ", err)

		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := ur.HTTPClient.Do(request)
	if err != nil {
		ur.Logger.Error("UtilityRequesterImpl.SendMail httpClient.Do ERROR ", err)

		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ur.Logger.Error("UtilityRequesterImpl.SendMail httpClient.Do ERROR ", err)

		return nil, err
	}

	var responses SendMailResponse

	err = json.NewDecoder(resp.Body).Decode(&responses)
	if err != nil {
		ur.Logger.Error("UtilityRequesterImpl.SendMail json.NewDecoder.Decode ERROR ", err)

		return nil, err
	}

	ur.Logger.Info("GOT RESPONSES FROM UTILITY SERVICES", responses)

	return &responses, nil
}
