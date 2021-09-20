package responseModel

import "encoding/xml"

type CalcPriceDeadlineResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    CalcPriceDeadlineBody
}

type CalcPriceDeadlineBody struct {
	XMLName                xml.Name `xml:"Body"`
	CalcPrecoPrazoResponse CalcPrecoPrazoResponse
}

type CalcPrecoPrazoResponse struct {
	XMLName              xml.Name `xml:"CalcPrecoPrazoResponse"`
	CalcPrecoPrazoResult CalcPrecoPrazoResult
}

type CalcPrecoPrazoResult struct {
	XMLName  xml.Name   `xml:"CalcPrecoPrazoResult"`
	Servicos []cServico `xml:"Servicos>cServico"`
}

type cServico struct {
	XMLName               xml.Name `xml:"cServico"`
	Codigo                string   `xml:"Codigo"`
	Valor                 string   `xml:"Valor"`
	PrazoEntrega          int      `xml:"PrazoEntrega"`
	ValorMaoPropria       string   `xml:"ValorMaoPropria"`
	ValorAvisoRecebimento string   `xml:"ValorAvisoRecebimento"`
	ValorValorDeclarado   string   `xml:"ValorValorDeclarado"`
	EntregaDomiciliar     string   `xml:"EntregaDomiciliar"`
	EntregaSabado         string   `xml:"EntregaSabado"`
	Erro                  int      `xml:"Erro"`
	MsgErro               string   `xml:"MsgErro"`
	ValorSemAdicionais    string   `xml:"ValorSemAdicionais"`
}
