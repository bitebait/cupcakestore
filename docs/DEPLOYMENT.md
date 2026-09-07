# Executar com Docker Compose

Requer Docker Engine ou Docker Desktop com o **plugin oficial Docker Compose v2 ou posterior** e BuildKit. Confira `docker compose version` antes de iniciar. Não é necessário instalar Go, compilador C, Node ou banco na máquina host.

## Desenvolvimento local

```sh
docker compose up --build -d --wait
docker compose ps
```

Abra <http://localhost:8080>. O padrão publica a porta somente em `127.0.0.1`. Para trocar a porta, use `HTTP_PORT=8081 docker compose up -d`. A porta interna permanece 8080 para manter aplicação e healthcheck alinhados.

O `.env` é opcional e serve à interpolação do Compose; ele não é copiado nem montado na imagem. O Compose configura explicitamente SQLite, seu caminho persistente e o modo local, independentemente dos valores legados de `APP_PORT`, `DB_PATH` e `DEV_MODE` do `.env`.

Para criar o primeiro administrador, preencha `ADMIN_EMAIL` e `ADMIN_PASSWORD` no `.env` antes de iniciar. Use uma senha única entre 12 e 72 bytes. Após a criação, remova essas duas variáveis e execute novamente `docker compose up -d`; isso recria o processo sem guardar a senha inicial no ambiente. O seed não promove contas existentes nem substitui senhas.

```sh
docker compose logs --tail=100 app
docker compose stop
docker compose start
```

As sessões são mantidas em memória: reiniciar o processo encerra os logins, mas preserva usuários, carrinhos e pedidos gravados no banco.

## Persistência e atualização

| Volume | Conteúdo | Caminho no contêiner |
| --- | --- | --- |
| `store-data` | SQLite, incluindo arquivos auxiliares de transações | `/data` |
| `product-images` | Fotos enviadas pelo painel e imagens iniciais | `/app/web/images` |

Os volumes pertencem ao projeto Compose. Mantenha o mesmo nome de projeto/diretório nas atualizações, ou fixe `COMPOSE_PROJECT_NAME` no `.env`. Volumes novos recebem a propriedade do usuário da imagem (UID/GID 10001); bind mounts importados precisam permitir leitura e escrita a esse usuário. A criação automática não importa o `gorm.db` ou as fotos existentes na máquina.

Para importar uma instalação antiga, pare o processo antigo e a aplicação no Compose, copie o banco e suas imagens para os respectivos volumes e preserve a propriedade 10001:10001. Não execute duas versões escrevendo simultaneamente no mesmo SQLite. Teste primeiro uma cópia, pois o startup ainda executa `AutoMigrate`.

Para um backup consistente simples, pare a aplicação e copie os dois diretórios juntos:

```sh
mkdir -p backups
docker compose stop app
docker compose cp app:/data backups/data
docker compose cp app:/app/web/images backups/images
docker compose start app
```

Use uma pasta de backup nova por execução e proteja seu acesso, pois contém dados de clientes. Valide a restauração em um projeto Compose separado. Antes de atualizar, faça backup; depois execute:

```sh
docker compose build --pull
docker compose up -d --wait
```

O volume de imagens mantém o caminho legado das fotos. Por isso também preserva os placeholders da primeira imagem instalada. Se esses assets mudarem em uma atualização, copie somente os três arquivos versionados do checkout atualizado:

```sh
docker compose cp web/images/600x400.svg app:/app/web/images/600x400.svg
docker compose cp web/images/cupcake-placeholder.svg app:/app/web/images/cupcake-placeholder.svg
docker compose cp web/images/placeholder.png app:/app/web/images/placeholder.png
```

Os uploads usam nomes aleatórios `.jpg`; os comandos acima não alteram essas fotos. Nunca apague o volume para atualizar assets.

`docker compose down` preserva os volumes. `docker compose down -v` apaga os dados persistentes. Voltar somente à imagem anterior não desfaz alterações de esquema: uma restauração pode exigir o backup correspondente.

## HTTPS para uma instalação pública

O arquivo base usa HTTP e cookies de desenvolvimento. Para publicar com TLS diretamente no Fiber, combine o override fornecido. No ambiente de implantação, configure:

```dotenv
BIND_ADDRESS=0.0.0.0
HTTP_PORT=443
TLS_DOMAIN=loja.example.com
TLS_DIRECTORY=/caminho/absoluto/certificados-da-loja
```

O diretório precisa existir e conter `fullchain.pem` e `privkey.pem`, legíveis pelo UID/GID 10001. Use uma cópia dedicada dos certificados reais e permissões restritas para a chave; não copie symlinks quebrados de outra árvore de certificados. O diretório é montado apenas para leitura e não entra na imagem.

```sh
docker compose -f compose.yaml -f compose.production.yaml config --quiet
docker compose -f compose.yaml -f compose.production.yaml up --build -d --wait
```

Esse modo define `DEV_MODE=false`, desativa recarga de templates e ativa cookies `Secure`. A porta pública 443 aponta para o listener TLS 8080, permitindo que o processo continue sem privilégios. O healthcheck usa HTTPS, verifica o certificado e seu nome DNS e acessa o loopback do contêiner; não usa `--insecure` e não depende de resolver o domínio na rede externa. A renovação/emissão de certificados pertence à infraestrutura; reinicie o serviço após renová-los. Não há listener HTTP na porta 80 ou emissão ACME automática neste Compose.

Se já houver proxy reverso, configure-o para acessar o upstream HTTPS com verificação do certificado, restringindo a publicação do upstream ao host/rede necessários. Não use `DEV_MODE=true` como atalho para terminar TLS no proxy: isso desativa cookies `Secure`. O projeto não habilita confiança automática em `X-Forwarded-*`; uma futura opção de TLS offload deve definir e testar explicitamente os proxies confiáveis.

## Operação e verificações

O build usa Go 1.26 com patches da tag oficial e CGO; o runtime usa Debian da mesma geração para compatibilidade com SQLite. A imagem contém apenas binário, templates, assets e bibliotecas/ferramentas de runtime. Atualizações da imagem base requerem rebuild com `--pull`.

O serviço executa como usuário não privilegiado, com raiz somente leitura, capacidades removidas e diretório temporário limitado. SQLite e imagens são as únicas áreas persistentes graváveis. Os arquivos estáticos usam `io/fs` e cache de compressão em memória, sem criar arquivos de cache ao lado dos assets. HTML com tokens CSRF permanece sem compressão.

`GET /livez` verifica o processo HTTP. `GET /readyz` verifica a conexão com o banco com timeout de dois segundos. Os endpoints não criam sessões nem expõem mensagens internas. O healthcheck do Compose usa `/readyz`; o estado `unhealthy` é diagnóstico, não substitui monitoramento externo nem reinicia automaticamente um contêiner ainda em execução.

Após iniciar, confira loja, login, inclusão de foto no painel, checkout em dinheiro e persistência após recriação do contêiner. O Pix exige configuração real da loja e validação operacional própria; o Compose não cria uma conta de provedor de pagamento. Mantenha uma única réplica enquanto SQLite e sessões em memória forem utilizados.

Validação realizada em 07/09/2026 com Docker Engine 29.7.2 e Compose oficial 5.5.1: build CGO, probes HTTP/HTTPS sem sessão, loja/login/assets comprimidos em raiz somente leitura, processo UID 10001, upload pelo painel e preservação do hash da foto após recriação. O teste HTTPS usou certificado descartável com nome DNS verificado e confirmou cookies `Secure`/`HttpOnly`; não valida DNS público, emissão ou renovação de certificados da implantação real.

Referências: [Docker Compose — serviços](https://docs.docker.com/reference/compose-file/services/), [persistência com volumes](https://docs.docker.com/engine/storage/volumes/) e [healthchecks do Fiber](https://docs.gofiber.io/middleware/healthcheck/).
