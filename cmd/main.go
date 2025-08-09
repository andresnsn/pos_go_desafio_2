package main

import (
	"encoding/json"
	"fmt"
	"fullcycle_desafios_go_2/internal/domain"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/getcep", getCepHandler)

	//Estou subindo um servidor apenas para ser o gatilho da aplicação, que deve responder a busca por /getcep
	fmt.Println("Servidor rodando localmente em http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}

func getCepHandler(w http.ResponseWriter, r *http.Request) {

	viaCepChannel := make(chan domain.ViaCep)

	brasilApiChannel := make(chan domain.BrasilAPI)

	go getViaCepApi(viaCepChannel)

	go getBrasilApi(brasilApiChannel)

	select {

	case viaCep := <-viaCepChannel:
		fmt.Printf("Foi pelo ViaCep, estes são os dados %+v\n", viaCep)

	case brasilApi := <-brasilApiChannel:
		fmt.Printf("Foi pelo BrasilApi, estes são os dados %+v\n", brasilApi)

	case <-time.After(1 * time.Second):
		fmt.Println("Nenhuma das APIs respondeu a tempo!")

	}

}

func getViaCepApi(ch chan<- domain.ViaCep) {

	resp, err := http.Get("https://viacep.com.br/ws/03180001/json")

	if err != nil {
		fmt.Println("Erro ao realizar requisição para ViaCep: ", err)
	}

	defer resp.Body.Close()

	var viacep domain.ViaCep
	err = json.NewDecoder(resp.Body).Decode(&viacep)

	if err != nil {
		fmt.Println("Erro ao converter para JSON: ", err)
	}

	ch <- viacep

	fmt.Println("Devolvi pro canal de ViaCep")

}

func getBrasilApi(ch chan<- domain.BrasilAPI) {

	resp, err := http.Get("https://brasilapi.com.br/api/cep/v1/01153000+cep")

	if err != nil {
		fmt.Println("Erro ao realizar a requisição para BrasilAPI: ", err)
	}

	defer resp.Body.Close()

	var brasilApi domain.BrasilAPI

	err = json.NewDecoder(resp.Body).Decode(&brasilApi)

	if err != nil {
		fmt.Println("Erro ao converter para JSON: ", err)
	}

	ch <- brasilApi

	fmt.Println("Devolvi pro canal de BrasilAPI")
}
