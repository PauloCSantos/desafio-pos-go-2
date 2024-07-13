package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type BrasilAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}
type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	var cep string = "13720000"
	var url1 string = "https://brasilapi.com.br/api/cep/v1/" + cep
	var url2 string = "http://viacep.com.br/ws/" + cep + "/json/"

	c1 := make(chan BrasilAPI)
	c2 := make(chan ViaCep)

	go func() {
		req, err := http.Get(url2)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
		}
		defer req.Body.Close()
		res, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao ler a resposta: %v\n", err)
		}
		var data ViaCep
		err = json.Unmarshal(res, &data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)
		}
		c2 <- data
	}()
	go func() {
		req, err := http.Get(url1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
		}
		defer req.Body.Close()
		res, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao ler a resposta: %v\n", err)
		}
		var data BrasilAPI
		err = json.Unmarshal(res, &data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)

		}
		c1 <- data
	}()

	select {
	case msg := <-c1:
		fmt.Println("BrasilApi respondeu primeiro")
		fmt.Printf("%+v\n", msg)
	case msg := <-c2:
		fmt.Println("ViaCep respondeu primeiro")
		fmt.Printf("%+v\n", msg)
	case <-time.After(time.Second * 1):
		fmt.Println("Erro de timeout: nenhuma api retornou a tempo")
	}
}
