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
		CepOrigem:        cep,
		CepDestino:       "70232090",
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

	var resp responseModel.CalcPriceDeadlineResponse
	err = utils.SoapCallHandleResponse("http://ws.correios.com.br/calculador/CalcPrecoPrazo.asmx", "http://tempuri.org/CalcPrecoPrazo", soapRequest, &resp)

	if err != nil {
		return StatusError{http.StatusInternalServerError, errors.Wrap(err, "Soap call error")}
	}

	return ResponseWriter(res, http.StatusOK, "", resp)
}
