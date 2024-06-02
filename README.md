# Cupcake Store

Projeto Integrador Transdisciplinar em Engenharia de Software II - UNICID - Cruzeiro Sul Virtual

## Como rodar o projeto *local*?

Clone o repositório:
~~~sh
git clone https://github.com/bitebait/cupcakestore.git
~~~

Navegue até a pasta do projeto:
~~~sh
cd cupcakestore/
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
- **Banco de Dados**: Sqlite3 (usando gorm – Golang ORM)
- **Hospedagem**: Linode (VPS) + Cloudflare
- **Plataforma**: Web (responsivo para tablet, smartphone e web)

### Estrutura do Projeto

A estrutura do projeto é organizada da seguinte forma:

- `bootstrap`
- `config`
- `controllers`
- `database`
- `docs`
- `middlewares`
- `models`
- `repositories`
- `routers`
- `services`
- `session`
- `utils`
- `views`
- `web`

### Tecnologias Utilizadas

- **Linguagens**: Go, JavaScript, CSS, HTML
- **Frameworks e Bibliotecas**: gorm (ORM para Golang)

### Autoria

Este projeto foi desenvolvido por William Schwaab como parte do Projeto Integrador Transdisciplinar em Engenharia de Software II - UNICID - Cruzeiro Sul Virtual.

Para mais informações, consulte a [documentação](https://github.com/bitebait/cupcakestore/tree/main/docs).


## Imagens

- **Loja:**
  ![Loja](https://github.com/bitebait/cupcakestore/blob/main/docs/store.png)

- **Painel de Admin:**
  ![Painel de Admin](https://github.com/bitebait/cupcakestore/blob/main/docs/dashboard.png)
