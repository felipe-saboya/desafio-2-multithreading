package main

import (
	"encoding/json"
	"fmt"
	"github.com/felipe-saboya/desafio-2-multithreading/configs"
	"github.com/felipe-saboya/desafio-2-multithreading/internal/dto"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	config, err := configs.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	var channel1 = make(chan dto.Result)
	var channel2 = make(chan dto.Result)

	go GetFromViaCep(config.PostalCodeHosts, channel1)
	go GetFromBrasilApi(config.PostalCodeHosts, channel2)

	select {
	case message := <-channel1: // rabbitmq
		fmt.Printf("API: %s; Endereço: %s", message.Api, message.Address)
	case message := <-channel2: // kafka
		fmt.Printf("API: %s; Endereço: %s", message.Api, message.Address)
	case <-time.After(time.Second):
		println("timeout")
	}
}

func GetFromBrasilApi(p []configs.PostalCodeHost, channel chan<- dto.Result) {
	for _, host := range p {
		if strings.Contains(host.Name, "BrasilApi") {
			var result = &dto.Result{}
			var brasilApi = host.Host

			resp, err := http.Get(brasilApi)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return
			}

			if resp.StatusCode == 200 {
				var brasilApiResult dto.BrasilApi
				err = json.Unmarshal(body, &brasilApiResult)
				if err != nil {
					return
				}

				result.Api = host.Name
				result.Address = brasilApiResult.Street + ", " + brasilApiResult.City + ", " + brasilApiResult.Cep

				channel <- *result
			}
		}
	}
}

func GetFromViaCep(p []configs.PostalCodeHost, channel chan<- dto.Result) {
	for _, host := range p {
		if strings.Contains(host.Name, "ViaCep") {
			var result = &dto.Result{}
			var viaCep = host.Host

			resp, err := http.Get(viaCep)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return
			}

			if resp.StatusCode == 200 {
				var viaCep dto.ViaCep
				err = json.Unmarshal(body, &viaCep)
				if err != nil {
					return
				}

				result.Api = host.Name
				result.Address = viaCep.Logradouro + ", " + viaCep.Localidade + ", " + viaCep.Cep

				channel <- *result
			}
		}
	}
}
