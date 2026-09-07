# Frontend: organização e manutenção

A aplicação continua renderizando HTML em Go. Não há bundler, framework JavaScript ou etapa de build frontend. Os recursos usados pelas páginas são locais.

## Onde alterar

| Recurso | Responsabilidade |
| --- | --- |
| `web/css/ui.css` | Cores, foco, componentes compartilhados, autenticação e painel |
| `web/css/storefront.css` | Vitrine, produto, carrinho, checkout e pedidos do cliente |
| `web/js/ui.js` | Menu, senha/confirmar senha, mensagens, CEP e preview de imagem |
| `web/js/storefront.js` | Quantidade, resumo de entrega, contador do carrinho e cópia do Pix |
| `web/js/admin.js` | Busca e seleção de produto no estoque |
| `views/layouts/store.html` | Estrutura da loja e navegação do cliente |
| `views/layouts/base.html` | Estrutura do painel e navegação administrativa |
| `views/products/form-fields.html` | Campos reutilizados no cadastro e edição de produto |
| `views/stock/product-picker.html` | Busca reutilizada em movimentação e consulta de estoque |
| `views/snippets/cart.html` | Formulário nativo de adicionar produto ou acesso ao login |
| `views/snippets/quantity.html` | Controle de quantidade compartilhado entre catálogo e carrinho |
| `views/snippets/order-status.html` | Nome do andamento compartilhado entre listas e detalhes, incluindo retirada |
| `views/snippets/pix-payment.html` | QR Code local, Pix Copia e Cola e instruções de confirmação manual |
| `views/snippets/state-options.html` | Opções de UF compartilhadas entre perfil e endereço da loja |
| `views/snippets/store-readiness.html` | Pendências de configuração que impedem recebimento e pagamento |
| `views/snippets/message.html` | Feedback persistente e acessível |
| `views/snippets/pagination.html` | Navegação que preserva a busca |

O CSS do AdminLTE permanece como base de compatibilidade para os formulários e tabelas existentes. Seu JavaScript, jQuery, jQuery UI, Cleave e demais plugins não são carregados pelas páginas renovadas. Login e cadastro usam somente o CSS próprio. Os arquivos vendorizados preservados não devem receber customizações da aplicação; mantenha essas alterações no CSS próprio.

## Convenções

- Preserve os tokens `--cream`, `--cocoa`, `--berry`, `--muted` e `--line` para manter loja e painel consistentes.
- Texto funcional usa `1rem` (16 px na configuração padrão); textos secundários têm mínimo de `.875rem` (14 px). No painel, reutilize `--ui-text`, `--ui-text-small` e `--ui-title`. Não reduza essas fontes no celular: ajuste a disposição dos elementos. Tabelas largas mantêm rolagem dentro do próprio contêiner para preservar a leitura das colunas.
- Use HTML semântico, labels explícitos e formulários nativos. A compra não deve depender da criação de formulários por JavaScript.
- Scripts são externos, carregados com `defer`, e usam `data-*` para encontrar componentes. Evite handlers inline, estado global e IDs novos para cada item.
- Todo formulário POST inclui `_csrf` com `CSRFToken` do contexto raiz. Dentro de `range`, passe o token explicitamente ao partial.
- Renderize conteúdo do usuário com o escape padrão de `html/template`; no JavaScript use `textContent`, não `innerHTML`.
- Use `money` apenas na apresentação. Inputs de preço e atributos `data-*` usados em cálculo mantêm valores numéricos, sem `R$` ou separadores locais.
- Mantenha foco visível, navegação por teclado, áreas de toque confortáveis e respeito a `prefers-reduced-motion`.
- Botões adicionais de quantidade e atualização do total são melhorias progressivas. Os valores definitivos são validados/calculados pelo backend.

## Verificação

```sh
go test ./views ./bootstrap
node --check web/js/ui.js
node --check web/js/storefront.js
node --check web/js/admin.js
```

Na revisão visual, use um banco temporário com dados fictícios. Exercite: vitrine e busca sem resultado; cadastro com confirmação de senha; login; perfil; adicionar/remover produtos; entrega versus retirada; finalização em dinheiro; consulta/cancelamento de pedido; busca de estoque e troca da seleção; upload inválido. Confira também o menu administrativo em tela estreita e a navegação por teclado.

As [capturas de tela](screenshots/) foram obtidas em Chromium com Playwright, em 1440×1000 e 390×844, usando dados fictícios e o placeholder de produto. Não representam produtos ou pedidos reais da loja. A verificação automatizada de acessibilidade usa axe-core (WCAG A/AA) e complementa a inspeção visual; não substitui uma avaliação completa com tecnologias assistivas.

Na verificação final de 07/09/2026, os 16 cenários avaliados não apresentaram violações detectadas pelo axe-core, erros de JavaScript ou overflow horizontal. O percurso incluiu cadastro, conclusão do perfil, compra em dinheiro com retirada e consulta do pedido; no painel, troca da seleção de estoque, entrada de quantidade, rejeição de arquivo inválido e abertura do menu móvel. A suíte `go test -race ./...` também passou após a integração dos templates e do formatter monetário.

Após a revisão de legibilidade, as fontes calculadas de menus, botões, campos, labels e tabelas do painel foram conferidas em 320, 390, 768 e 1440 px de largura. Mantêm 16 px e não provocam overflow da página; tabelas usam rolagem local quando necessário. As capturas refletem essa escala maior.

## Ajustes comerciais e administrativos

O carrinho permite alterar a quantidade com um formulário POST por produto, preservando o token CSRF. O mesmo partial atende ao catálogo; os botões de incremento só aparecem quando o JavaScript está disponível. O campo numérico e o botão de envio continuam funcionando sem scripts. A quantidade máxima exibida acompanha o estoque consultado; o servidor confere novamente o saldo ao gravar.

O painel de pedidos obtém as próximas etapas de `Order.AvailableTransitions`, que reutiliza a regra de domínio em vez de repetir transições em HTML. Retirada não oferece envio, e cancelamento passa pela página de confirmação. O texto de Pix distingue conferência no banco e cancelamento do pedido de devolução do dinheiro.

A listagem de usuários reúne informações secundárias na mesma célula, preserva a busca no HTML e oferece ações com nomes explícitos. Os contêineres de tabelas são regiões acessíveis por teclado, permitindo rolagem horizontal sem diminuir a fonte. Configurações usam uma largura maior no tablet e links funcionais para a visão geral.

## Recursos de terceiros preservados

A revisão dos templates, código Go, scripts próprios e referências `url(...)` do CSS removeu **2.007 arquivos rastreados sem uso (78.461.173 bytes)** de `web/dist` e `web/plugins`. Saíram JavaScript legado, plugins sem consumidores, CSS alternativo, mapas de código e imagens de demonstração. Arquivos não rastreados não participaram da remoção.

Permanecem 19 recursos usados diretamente ou referenciados pelo CSS, somando 4.402.990 bytes, além de três arquivos de licença:

| Caminho | Uso |
| --- | --- |
| `web/dist/css/adminlte.min.css` | CSS de compatibilidade para formulários e tabelas; inclui Bootstrap 4.6.1 |
| `web/dist/img/logo.png` | Marca nos layouts e na autenticação |
| `web/dist/img/favicon.png` | Ícone servido pelo middleware de favicon |
| `web/plugins/fontawesome-free/css/all.min.css` | Ícones das telas administrativas e formulários |
| `web/plugins/fontawesome-free/webfonts/fa-{brands-400,regular-400,solid-900}.{eot,svg,ttf,woff,woff2}` | 15 arquivos de fontes referenciados pelo CSS preservado |

O CSS do AdminLTE contém somente imagens `data:` embutidas. Seu comentário de sourcemap foi removido porque o mapa deixou de ser distribuído; os comentários de autoria e licença foram preservados. Nenhum JavaScript de terceiros é carregado ou distribuído nesses diretórios.

As licenças foram obtidas das versões originais: [AdminLTE 3.2.0](https://github.com/ColorlibHQ/AdminLTE/blob/v3.2.0/LICENSE), [Bootstrap 4.6.1](https://github.com/twbs/bootstrap/blob/v4.6.1/LICENSE) e [Font Awesome 5.15.4](https://github.com/FortAwesome/Font-Awesome/blob/5.15.4/LICENSE.txt). Cópias acompanham os arquivos em `web/dist/LICENSE.AdminLTE.txt`, `web/dist/LICENSE.Bootstrap.txt` e `web/plugins/fontawesome-free/LICENSE.txt`.

## Recebimento e Pix

O checkout oferece entrega apenas quando o serviço está ativo e o cliente possui endereço válido. Retirada depende de endereço, cidade e UF válidos da loja; o cliente precisa somente dos dados de contato. Essas condições usam os métodos do domínio (`HasContactDetails`, `HasDeliveryAddress`, `IsPickupAvailable` e `IsPixAvailable`). Quando falta uma opção, a interface apresenta a correção possível; sem recebimento ou pagamento, a finalização fica indisponível. O backend repete as validações antes de reservar e pagar.

Antes do pagamento, a página mostra o perfil atual e as configurações atuais da loja. Depois da escolha, o pedido exibe os dados registrados na operação. Pedidos em aberto apresentam pagamento e recebimento como escolhas pendentes, sem anunciar uma forma já confirmada.

O Pix permanece na página do pedido. `PixQR` contém PNG em base64 e o template usa o prefixo literal `data:image/png;base64,`, mantendo o escape do HTML. `PixString` fica em um campo somente leitura, utilizável para copiar manualmente sem JavaScript. O botão de cópia usa Clipboard API quando disponível e orienta a cópia manual caso falhe. Um erro de validação do payload oculta tanto o código quanto o QR.

A aprovação do Pix exige um formulário administrativo separado, com checkbox obrigatório de conferência do recebimento no extrato. A interface não oferece aprovação pelo seletor comum de status. Após confirmação, o cliente solicita cancelamento ou devolução à loja; `CanCustomerCancel` controla as ações exibidas e o backend protege a operação. Cancelamento administrativo não significa reembolso automático.

Para validar esta etapa, inclua cenários de retirada sem endereço do cliente, loja sem endereço de retirada, entrega sem endereço completo, ausência de pagamento, Pix local no celular, cópia manual sem JavaScript, confirmação pelo administrador e cancelamento após Pix confirmado. As UFs do perfil e da loja usam `select`; automações devem usar `selectOption`, preservando os testes de tipografia de 16 px.

Na validação desta etapa em 07/09/2026, o percurso real no Chromium confirmou 15 verificações funcionais, incluindo Pix local sem JavaScript, cópia exata do payload, bloqueio de dinheiro após alteração da taxa entre revisão e envio, preservação do Pix já emitido e confirmação/cancelamento conforme o perfil. As sete telas e estados avaliados pelo axe-core não apresentaram violações detectadas, erros de JavaScript, requisições externas, violações de CSP ou overflow horizontal. Os controles visíveis no celular mantiveram 16 px. As capturas `pix-payment-desktop.png`, `pix-payment-mobile.png` e `checkout-mobile.png` usam somente dados fictícios; nenhuma transação bancária foi realizada.
