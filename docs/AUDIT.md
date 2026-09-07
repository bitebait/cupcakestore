# Auditoria de engenharia e modernização

Auditoria iniciada em 7 de setembro de 2026. Objetivo: preservar a loja de cupcakes e o painel administrativo, simplificando o código e corrigindo segurança, integridade e manutenção. A análise foi dividida entre backend, frontend e infraestrutura, com revisão cruzada e integração final.

## Critério de arquitetura

Manter o monólito Go com HTML renderizado no servidor. Não adicionar microsserviços, uma SPA ou abstrações sem necessidade. As rotas compõem dependências, controllers tratam HTTP, services implementam casos de uso, repositories persistem os dados e models descrevem o domínio. Serviços e repositórios devem evoluir sem importar Fiber; sessão, cookies e redirecionamentos pertencem à borda HTTP.

Operações comerciais com vários writes precisam de transação explícita. O checkout grava reserva, pedido, vínculo do carrinho e endereço juntos; cancelar é idempotente. Validação do formulário complementa a validação de domínio, sem permitir editar IDs, permissões, saldos ou totais enviados pelo cliente.

## Achados tratados

| Área | Falha identificada | Correção e evidência |
| --- | --- | --- |
| Autorização | Redirecionamento não interrompia handlers; sessão conservava privilégios revogados | Permissões relidas no banco, bloqueio encerra requisição; regressões de conta inativa, removida e administrador rebaixado |
| Formulários | Binding direto de entidades aceitava IDs, permissões, estoque e preço total | Campos permitidos explícitos; testes de cadastro e produto e fluxo real de pagamento adulterado |
| Autenticação | Sessão não rotacionada; alteração de senha aceitava senha curta | Renovação no login, validação de senha atual/tamanho, e-mail normalizado, limitação de tentativas |
| Cadastro | Usuário e detalhes do perfil podiam ficar parcialmente gravados | Criação em transação |
| Carrinho | Erros de escrita ignorados, quantidades negativas e edição após checkout | Validação e serialização de escritas; carrinho finalizado imutável |
| Estoque | Saídas excediam saldo e cancelamento repetido criava estoque | Débito condicional atômico; devolução única; testes concorrentes e rollback |
| Pedidos | Dono do carrinho não validado antes da criação; status sem regra de avanço | Propriedade conferida na transação, transições válidas e preservação de preço/endereço |
| Pagamento | Retirada ignorada; dinheiro terminava com erro; emissão dependia de gerador externo | Total calculado no servidor, dinheiro funcional, Pix BR Code local e confirmação manual explícita |
| Upload | Nome sem extensão causava panic; imagem comprimida podia alocar memória excessiva | Nome aleatório e JPEG de saída, validação de conteúdo, 4 MiB e 16 milhões de pixels |
| CSRF | Mutações em GET e proteção dependente de formulários inconsistentes | POST para logout, remoção e criação de checkout; token ligado à sessão |
| Configuração | Import exigia `.env`; DSN PostgreSQL inválido; seed com credencial pública | Configuração explícita validada, ambiente sem `.env`, DSN escapado e admin opt-in |
| Seeds | Edição do e-mail da loja criava outro registro | Seed idempotente e transacional, Pix desativado até configurar |
| Frontend | Rastreador fixo inclusive em login; JS duplicado; navegação e formulários quebrados | Remoção do rastreamento, assets locais, validação nativa, rótulos e mensagens acessíveis |
| Qualidade | Ausência de testes e automação | Testes isolados de domínio/HTTP/templates; CI build/race/vet/gofmt |

## Migração efetiva para Fiber 3.5

O projeto agora requer Go 1.26 ou superior, com patches atuais. Além da atualização de módulos, o código foi migrado para:

- `fiber.Ctx` como interface, `fiber.Query[int]` e `fiber.Locals[T]` para leitura tipada.
- `Bind().Body` com DTOs restritos e `StructValidator` compartilhado. Os serviços continuam validando invariantes do domínio.
- Sessões por `session.NewWithStore`/`FromContext`, sem `Store.Get`/`Save` manual. Identidade armazenada como ID; permissões sempre consultadas no banco.
- Renovação da sessão e remoção do token CSRF no login; token antigo e cookie anterior são rejeitados em testes.
- CSRF com `extractors.FromForm` e helpers oficiais de contexto.
- `Redirect().Status(...).To(...)`, preservando os contratos HTTP existentes.
- `static.New` e `ListenConfig`, preservando HTTPS e shutdown. Compressão só para assets públicos, sem comprimir HTML contendo token CSRF.
- Logging padrão `log/slog` nos repositories e autenticação sem HTTP nos services. Registro de login atualiza apenas datas, sem sobrescrever permissões de um snapshot antigo.
- Updates de usuário/perfil com colunas explícitas, sem fallback de `Save` para criar registros.

Fontes de migração: [guia Fiber](https://docs.gofiber.io/whats_new/), [sessões](https://docs.gofiber.io/middleware/session/), [CSRF](https://docs.gofiber.io/middleware/csrf/) e [validação](https://docs.gofiber.io/guide/validation/). O código e as assinaturas também foram conferidos no módulo v3.5.0 instalado, pois a documentação Next pode acompanhar versões posteriores.

A [auditoria de dependências](DEPENDENCY_AUDIT.md) registra versões e evidências: 16 avisos alcançáveis e 12 em pacotes importados antes; zero nessas duas categorias depois. Um aviso permanece somente em OpenPGP, pacote não importado. Esse resultado não cobre automaticamente bibliotecas JavaScript/CSS vendorizadas.

## Renovação do frontend

Vitrine, produto, carrinho, checkout, pedidos e autenticação foram redesenhados; o painel recebeu navegação e componentes coerentes com a loja. Os scripts inline e dependências de jQuery/Bootstrap JS/jQuery UI/Cleave foram substituídos por três módulos vanilla compartilhados. Formulários de produto e busca de estoque agora reutilizam partials. O CSS existente é mantido apenas onde dá suporte aos formulários/tabelas legados.

A validação em Chromium cobriu compra completa, cadastro, perfil, estoque, upload e menu móvel, com capturas desktop/celular. Axe-core encontrou contrastes insuficientes, corrigidos durante a revisão. O [guia frontend](FRONTEND.md) documenta a estrutura e como manter as convenções.

## Segunda revisão: comércio, composição e implantação

Uma nova rodada com três agentes revisou comércio/Pix, experiência administrativa e implantação. A integração manteve o monólito e usou composição explícita em vez de introduzir containers de dependências, eventos ou novos frameworks.

- **Regressão de entrega corrigida:** `OrderRepository.Update` grava somente o status. `UpdatePayment` valida e grava a escolha de pagamento e seus valores. Alterar taxa ou disponibilidade de entrega não impede o administrador de avançar um pedido Pix já emitido. Testes também garantem que o caminho de status ignora tentativas de reescrever termos financeiros e que uma nova escolha continua validando a taxa atual.
- **Pix:** QR Code e Copia e Cola gerados localmente, com valor fixo, referência do pedido e CRC16. Chaves são normalizadas/validadas, incluindo CNPJ alfanumérico; nome e cidade usam os mesmos limites em configuração e geração. Reabrir um Pix pendente reutiliza seus dados, sem recalcular seus termos. A leitura valida estrutura, checksum, valor e referência antes de mostrar o código; links externos legados não são usados. Alterações de forma de pagamento/entrega são bloqueadas após emissão; entrega desativada exige escolha explícita de retirada.
- **Carrinho:** nova ação de quantidade reutiliza a escrita transacional existente, obtém o proprietário da sessão e valida produto, saldo e congelamento após checkout. Campos adulterados não escolhem outro carrinho. O formulário funciona sem JavaScript.
- **Estoque:** entrada extrema podia ultrapassar a capacidade de um inteiro e corromper a representação do saldo no SQLite. A atualização agora limita a soma atomicamente; a regressão reproduziu a falha antes da correção e verifica saldo/histórico preservados depois.
- **Painel:** próximas etapas vêm de `Order.AvailableTransitions`, que reutiliza a regra de domínio. Retirada não permite envio; cancelamento tem confirmação e explica o tratamento separado de valores Pix. Listagem de usuários e formulários de configuração foram simplificados; o contraste do indicador de conta ativa foi corrigido.
- **Administração:** transações e bloqueios preservam pelo menos um administrador ativo, impedindo sua exclusão, desativação ou perda de permissão. Escritas com versões antigas são rejeitadas. A configuração usa atualização otimista com colunas explícitas, sem recriar registros excluídos por um `Save` antigo.
- **Indicadores:** pedidos concluídos incluem somente os efetivamente entregues ou retirados. Pedidos enviados permanecem na relação de pedidos em andamento até a entrega.
- **Composição:** `bootstrap/routes.go` monta serviços/repositórios compartilhados uma vez e injeta controllers nas funções de registro. Removidos structs, interface e constructors de roteadores que não acrescentavam comportamento. Rotas e middleware de autenticação não acessam mais o banco global. O teste de autorização usa a composição completa e verifica uma consulta de conta por requisição.
- **Navegador:** Helmet do Fiber aplica CSP sem `unsafe-inline`/`unsafe-eval`, proteção contra enquadramento e detecção incorreta de MIME. HTML com tokens CSRF e JSON recebem `private, no-store`; assets continuam cacheáveis. Formulários permitem somente destinos da própria loja, inclusive o Pix local.
- **Implantação:** Docker multi-stage com CGO, runtime sem root, volumes, raiz somente leitura, probes sem sessão e override HTTPS. A CI constrói a imagem e testa readiness e vitrine; o [guia operacional](DEPLOYMENT.md) cobre persistência, atualização e certificados.
- **Assets:** removidos 2.007 arquivos vendorizados sem uso (78.461.173 bytes), preservando as 19 dependências efetivas de CSS, fontes e imagens e suas licenças. A compressão usa `io/fs` e cache em memória: o teste real do contêiner identificou e corrigiu o 404 de CSS que ocorria ao tentar criar cache na raiz somente leitura.

Validação integrada: `go test -race ./...`, testes no navegador com CSP ativa, edição de quantidade seguida de checkout, finalização administrativa de retirada e login/carrinho com JavaScript desativado. A checagem automatizada WCAG A/AA complementa a inspeção visual.

O Pix local elimina chamadas a geradores externos: gerar instruções não cria uma cobrança remota órfã. O projeto não consulta o DICT nem confirma liquidação bancária; validação sintática não prova registro ou titularidade da chave. A confirmação é manual após conferir o extrato. Um QR estático já copiado não pode ser revogado pelo cancelamento local e pode receber pagamentos repetidos ou tardios, que exigem conciliação/devolução pela loja. Os testes incluem vetor oficial do BCB, checksum, formatos de chave, limites monetários e divergência entre código e pedido. Uma verificação independente com ZXing decodificou o PNG gerado e recuperou exatamente o vetor oficial. Veja [operação e limites do Pix](PIX.md).

Fontes dos cabeçalhos: [Helmet Fiber](https://docs.gofiber.io/next/middleware/helmet/) e [CSP form-action](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/form-action). As opções usadas foram conferidas no módulo Fiber 3.5 instalado.

## Compatibilidade de dados e operação

- Não apagamos ou reinicializamos o banco do usuário. Testes operam exclusivamente em bancos temporários.
- O novo campo `orders.stock_reserved` é verdadeiro para reservas criadas pelo fluxo transacional. Pedidos legados começam com falso, pois o sistema anterior não permite inferir se houve reserva ou devolução duplicada. Reconcilie esses pedidos com os movimentos históricos antes de marcar uma reserva legada como existente; não execute uma atualização geral desse campo.
- Seeds não redefinem contas antigas. Instalações legadas precisam revisar a conta criada com credencial pública.
- SQLite mantém uma conexão para serializar writes e aplicar foreign keys. PostgreSQL tem DSN testado; os cenários comerciais concorrentes ainda precisam ser executados em um PostgreSQL real.
- Sessões continuam em memória: reiniciar o processo encerra logins. A migração do Fiber também exige novo login.

## Próximas melhorias com escopo definido

| Prioridade | Trabalho | Critério de conclusão |
| --- | --- | --- |
| Alta | Migrações versionadas e ensaio em cópia de banco legado | Roll-forward documentado e reconciliação de reservas anteriores sem perda de dados |
| Alta | Expiração de reservas abandonadas | Job idempotente cancela somente pedidos elegíveis e restitui estoque uma única vez |
| Alta | Testes comerciais em PostgreSQL | Mesmos cenários de concorrência, rollback e integridade passando no banco suportado |
| Média | Valores monetários em centavos inteiros | Migração preserva totais históricos e elimina aritmética monetária em float64 |
| Média | Sessões persistentes quando necessário | Reinício e múltiplas instâncias mantêm sessão/CSRF com expiração consistente |
| Média | Listagens distinguirem vazio de erro | Interfaces devolvem erro, controllers apresentam falha sem aparentar lista vazia |
| Média | Snapshot de nome e descrição dos itens | Pedido histórico não muda quando o catálogo é editado |
| Média | Avaliação com leitor de tela | Revisão assistiva completa além dos testes automatizados de acessibilidade |
| Opcional | Confirmação automática de pagamento Pix | Integração com o PSP da loja, credenciais fora do código, eventos autenticados e conciliação idempotente testados em homologação |

A auditoria não prova ausência de todas as falhas. Os itens acima são riscos e limitações observados, com critérios concretos para próximas mudanças; não foram marcados como concluídos sem evidência.
