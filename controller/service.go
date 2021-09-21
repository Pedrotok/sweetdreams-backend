package controller

import (
	controller "SweetDreams/controller/requestModel"
	"SweetDreams/controller/responseModel"
	"SweetDreams/controller/utils"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetDeliveryInfo(db *mongo.Database, res http.ResponseWriter, req *http.Request) error {
	cep := req.FormValue("cep")
	amountString := req.FormValue("amount")
	_, err := strconv.ParseInt(amountString, 10, 64)

	if cep == "" {
		return StatusError{http.StatusBadRequest, errors.Wrap(err, "empty CEP")}
	}

	if err != nil {
		return StatusError{http.StatusBadRequest, errors.Wrap(err, "invalid amount")}
	}

	soapRequest := &controller.CalcPriceDeadlineRequest{
		XMLNs:            "http://tempuri.org/",
		Empresa:          "",
		Senha:            "",
		Servico:          "04510",
		CepOrigem:        "70232090",
		CepDestino:       cep,
		Peso:             "1",
		Formato:          1,
		Comprimento:      45,
		Largura:          45,
		Altura:           15,
		Diametro:         45,
		MaoPropria:       "N",
		ValorDeclarado:   0,
		AvisoRecebimento: "S",
	}

	var calcPriceDeadlineResponse responseModel.CalcPriceDeadlineResponse
	err = utils.SoapCallHandleResponse("http://ws.correios.com.br/calculador/CalcPrecoPrazo.asmx",
		"http://tempuri.org/CalcPrecoPrazo", soapRequest, &calcPriceDeadlineResponse)
	if err != nil {
		return StatusError{http.StatusInternalServerError, errors.Wrap(err, "Soap call error")}
	}

	resp, err := parsePriceDeadlineResponse(calcPriceDeadlineResponse)

	return ResponseWriter(res, http.StatusOK, "", resp)
}

func parsePriceDeadlineResponse(data responseModel.CalcPriceDeadlineResponse) (*responseModel.DeliveryInfoResponse, error) {
	services := data.Body.CalcPrecoPrazoResponse.CalcPrecoPrazoResult.Servicos
	if len(services) == 0 {
		return nil, errors.New("Correios didn't return any services")
	}

	msgError := services[0].MsgErro
	if msgError != "" {
		return nil, errors.New(msgError)
	}

	deliveryInfoResponse := &responseModel.DeliveryInfoResponse{
		DeliveryType: getDeliveryTypeFromCode(services[0].Codigo),
		Value:        utils.ParsePriceFromStringToInt(services[0].Valor),
		DeliveryTime: services[0].PrazoEntrega,
	}
	return deliveryInfoResponse, nil
}

func getDeliveryTypeFromCode(code string) string {
	if code == "4510" {
		return "PAC"
	} else if code == "4014" {
		return "SEDEX"
	}

	return ""
}
