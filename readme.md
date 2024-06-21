PROJETO DESAFIO SEMANA 18 IMERSÃO FULLCYCLE
objetivos:

Neste desafio, você deverá criar uma aplicação Golang que gerará uma API REST.

No repositório da imersão existe um arquivo na pasta nextjs-front-end/node-api/data.js. Você irá pegar este arquivo e convertê-lo para JSON e salvar na raiz do seu projeto Golang como “data.json”.

Este arquivo será o banco de dados do seu projeto, a aplicação irá expor algumas rotas e sua aplicação deverá ler o arquivo JSON ao iniciar o programa com o go run main.go e manipular os dados sempre na memória (não salve nada no arquivo JSON).
----------------------------
Sua aplicação deverá expor os endpoints:

GET /events - Listar todos os eventos
GET /events/:eventId - Listar os dados de um evento
GET /events/:eventId/spots - Listar os lugares de um evento
POST /event/:eventId/reserve - Reservar um lugar
----------------------------
Dados:
spots: array com spots, exemplo [A1, B2]
Não poderá reservar o mesmo spot duas vezes, deverá retornar um erro 400

####
iniciando o projeto
----------------------------
bash
go run internal/cmd/main.go
----------------------------

Para testar o projeto pelo arquivo testes.http
instalar a extensão Rest Client no vscode