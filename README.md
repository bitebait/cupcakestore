# Cupcake Store

Projeto Integrador Transdisciplinar em Engenharia de Software II - UNICID - Cruzeiro Sul Virtual

_Esse é um projeto que funcionará como uma loja online de cupcakes para uma pequena empresa. 
Ele faz parte de um trabalho acadêmico e utiliza conceitos aprendidos durante o curso. 
Tenha em mente que este é um projeto acadêmico e não atende aos requisitos para ser utilizado em produção._

 #### 🔥 Sinta-se à vontade para contribuir com o código (; 🔥

## Como rodar o projeto *local*?

Clone o repositório:
~~~sh
git clone https://github.com/bitebait/cupcakestore.git
~~~

Navegue até a pasta do projeto:
~~~sh
cd cupcakestore/
~~~

Crie um novo arquivo .env com base no .env.example e atualize suas configurações:
~~~sh
cp .env.example .env 
~~~

Atualize os módulos:
~~~go
go mod tidy
~~~

Rode o projeto:
~~~go
go run .
~~~

### Informações Adicionais

- **Link da Solução em Funcionamento:** [Cupcake Store](https://cupcakestore.schwaab.me:2053/store)
- **Usuário DEMO ADMIN**: `admin@admin.com` / `admin@admin.com`
- **Linguagem Back-end**: Golang
- **Front-end**: HTML+CSS+JS ([AdminLTE Bootstrap Admin Dashboard](https://adminlte.io/))
- **Banco de Dados**: Sqlite3 (usando gorm – Golang ORM)
- **Hospedagem**: Linode (VPS) + Cloudflare
- **Plataforma**: Web (responsivo para tablet, smartphone e web)

### Estrutura do Projeto

A estrutura do projeto é organizada da seguinte forma:

- `bootstrap`: *Contém arquivos relacionados à inicialização do projeto.*
- `config`: *Responsável pelas configurações do ambiente.*
- `controllers`: *Engloba os controladores da aplicação.*
- `database`: *Arquivos relativos ao banco de dados, incluindo scripts de inicialização.*
- `docs`: *Documentação do projeto.*
- `middlewares`: *Implementação de middlewares, como controle de autenticação.*
- `models`: *Define os modelos de dados utilizados na aplicação.*
- `repositories`: *Responsável pelo acesso e manipulação dos dados.*
- `routers`: *Configuração das rotas da aplicação.*
- `services`: *Serviços oferecidos pela aplicação.*
- `session`: *Gerenciamento de sessões de usuário.*
- `utils`: *Utilitários diversos.*
- `views`: *Templates e arquivos relacionados à visualização da aplicação.*
- `web`: *Recursos web, como favicons, imagens, assets, etc.*

### Tecnologias Utilizadas

- **Linguagens**: Go, JavaScript, CSS, HTML
- **Frameworks e Bibliotecas**: [GO Fiber Framework](https://github.com/gofiber/fiber) & [GORM](https://gorm.io/index.html) (ORM para Golang)
- **Front-end**: HTML+CSS+JS ([AdminLTE Bootstrap Admin Dashboard](https://adminlte.io/))

### Autoria

Este projeto foi desenvolvido por William Schwaab (<william@schwaab.me>) como parte do Projeto Integrador Transdisciplinar em Engenharia de Software II - UNICID - Cruzeiro Sul Virtual.

Para mais informações, consulte a [documentação](https://github.com/bitebait/cupcakestore/tree/main/docs).


## Imagens

- **Loja:**
  ![Loja](https://github.com/bitebait/cupcakestore/blob/main/docs/store.png)

- **Painel de Admin:**
  ![Painel de Admin](https://github.com/bitebait/cupcakestore/blob/main/docs/dashboard.png)
