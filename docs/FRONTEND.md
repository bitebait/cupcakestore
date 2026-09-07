# Frontend: organização e manutenção

A aplicação continua renderizando HTML em Go. Não há bundler, framework JavaScript ou etapa de build frontend. Os recursos usados pelas páginas são locais.

## Onde alterar

| Recurso | Responsabilidade |
| --- | --- |
| `web/css/ui.css` | Cores, foco, componentes compartilhados, autenticação e painel |
| `web/css/storefront.css` | Vitrine, produto, carrinho, checkout e pedidos do cliente |
| `web/js/ui.js` | Menu, senha/confirmar senha, mensagens, CEP e preview de imagem |
| `web/js/storefront.js` | Quantidade, resumo de entrega e contador do carrinho |
| `web/js/admin.js` | Busca e seleção de produto no estoque |
| `views/layouts/store.html` | Estrutura da loja e navegação do cliente |
| `views/layouts/base.html` | Estrutura do painel e navegação administrativa |
| `views/products/form-fields.html` | Campos reutilizados no cadastro e edição de produto |
| `views/stock/product-picker.html` | Busca reutilizada em movimentação e consulta de estoque |
| `views/snippets/cart.html` | Formulário nativo de adicionar produto ou acesso ao login |
| `views/snippets/message.html` | Feedback persistente e acessível |
| `views/snippets/pagination.html` | Navegação que preserva a busca |

O CSS do AdminLTE permanece como base de compatibilidade para os formulários e tabelas existentes. Seu JavaScript, jQuery, jQuery UI, Cleave e demais plugins não são carregados pelas páginas renovadas. Login e cadastro usam somente o CSS próprio. Os arquivos antigos em `web/plugins` não devem receber customizações da aplicação.

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
