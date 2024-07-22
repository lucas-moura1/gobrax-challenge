# gobrax-challenge

Este projeto é uma solução para o desafio gobrax, focado no cadastro e vinculação de motoristas e veículos. Abaixo estão detalhadas as entidades, funcionalidades e instruções para execução do projeto.

## Entidades

### Motorista
- `name`: Nome do motorista.
- `lastName`: Último nome do motorista.
- `email`: E-mail de contato do motorista.
- `phone`: Telefone de contato do motorista.
- `license`: Documento da CNH.
- `licenseType`: Qual categoria da CNH.

### Veículo
- `brand`: Marca do veículo.
- `vehicleModel`: Modelo do veículo.
- `year`: Ano de fabricação do veículo.
- `plate`: Placa do veículo.


## Funcionalidades

- **Gestão de Motorista**:

    - Criação (`POST /drivers`)
```json
Ex:
{
    "name": "John",
    "lastName": "Deo",
    "email": "john@test.com",
    "phone": "21984736452",
    "license": "928843839",
    "licenseType": "B"
}
```
    - Adicionar veículo (`POST /drivers/{id}/vehicle`)
```json
Ex:
{
	"plate": "HIJ-1231",
	"brand": "Ford",
	"vehicleModel": "Focus",
	"year": 2007
}
```
    - Listagem (`GET /drivers`)
    - Obter por ID (`GET /drivers/{id}`)
    - Atualização (`PATCH /drivers/{id}`)
    - Remoção (`DELETE /drivers/{id}`)

- **Gestão de Veículos**:

    - Listagem (`GET /vehicles`)
    - Obter por ID (`GET /vehicles/{id}`)
    - Atualização (`PATCH /vehicles/{id}`)
    - Remoção (`DELETE /vehicles/{id}`)

## Como Executar o Projeto

Deve ter:
- ***Docker*** e ***docker-compose*** instalado na máquina.

### Para executar

- ```git clone <url_repositorio>``` : clonar o repositório;
- ```docker compose up```: rodar a aplicação

### Para rodar os testes unitários
- `go test ./...`

Para acessar a API diretamente é preciso acessar ```http://localhost:8080``` + o endPoint.
