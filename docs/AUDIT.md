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
| Pagamento | Retirada ignorada; dinheiro terminava com erro; cliente HTTP sem timeout | Total calculado no servidor, dinheiro funcional, timeout/limite de resposta/URL e encoding no cliente Pix |
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
| Média | Injeção das dependências na inicialização | Remover dependência global de banco/config das fábricas de rotas, sem duplicar containers |
| Média | Snapshot de nome e descrição dos itens | Pedido histórico não muda quando o catálogo é editado |
| Média | Avaliação com leitor de tela e limpeza dos assets legados | Revisão assistiva completa e remoção dos arquivos vendorizados sem uso após inventário |
| Média | Contrato e confirmação de pagamento Pix | Provedor real validado em ambiente de teste; processamento idempotente de confirmação/reconciliação |

A auditoria não prova ausência de todas as falhas. Os itens acima são riscos e limitações observados, com critérios concretos para próximas mudanças; não foram marcados como concluídos sem evidência.
