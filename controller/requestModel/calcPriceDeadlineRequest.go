package controller

import "encoding/xml"

type CalcPriceDeadlineRequest struct {
	XMLName          xml.Name `xml:"CalcPrecoPrazo"`
	XMLNs            string   `xml:"xmlns,attr"`
	Empresa          string   `xml:"nCdEmpresa"`
	Senha            string   `xml:"sDsSenha"`
	Servico          string   `xml:"nCdServico"`
	CepOrigem        string   `xml:"sCepOrigem"`
	CepDestino       string   `xml:"sCepDestino"`
	Peso             string   `xml:"nVlPeso"`
	Formato          int      `xml:"nCdFormato"`
	Comprimento      int      `xml:"nVlComprimento"`
	Altura           int      `xml:"nVlAltura"`
	Largura          int      `xml:"nVlLargura"`
	Diametro         int      `xml:"nVlDiametro"`
	MaoPropria       string   `xml:"sCdMaoPropria"`
	ValorDeclarado   int      `xml:"nVlValorDeclarado"`
	AvisoRecebimento string   `xml:"sCdAvisoRecebimento"`
}
