# Auditoria de dependências Go — 07/09/2026

A migração para Fiber 3.5.0 e a atualização das dependências eliminaram os avisos de vulnerabilidade em funções alcançáveis e pacotes importados identificados pelo `govulncheck`.

| Classificação do aviso | Antes | Depois |
| --- | ---: | ---: |
| Função alcançável no grafo estático | 16 | 0 |
| Pacote importado, sem função vulnerável alcançável | 12 | 0 |
| Apenas módulo, sem pacote vulnerável importado | 33 | 1 |

Cada aviso é contado uma vez, na classificação de maior alcance. A análise estática identifica caminhos possíveis; não demonstra exploração de cada aviso.

## Atualizações

| Dependência | Versão anterior | Versão adotada |
| --- | --- | --- |
| Fiber | v2.52.5 | v3.5.0 |
| Fiber HTML templates | v2.0.5 | v2.1.3 |
| Sprig | v3.2.3 | v3.3.0 |
| Validator | v10.15.5 | v10.30.4 |
| GORM | v1.25.10 | v1.31.2 |
| Driver PostgreSQL | v1.5.11 | v1.6.2 |
| Driver SQLite | v1.5.4 | v1.6.0 |
| pgx | v5.5.5 | v5.10.0 |
| fasthttp | v1.52.0 | v1.74.0 |
| golang.org/x/crypto | v0.19.0 | v0.56.0 |
| golang.org/x/image | revisão de 2019 | v0.45.0 |
| golang.org/x/net | v0.21.0 | v0.58.0 |
| golang.org/x/text | v0.14.0 | v0.41.0 |

As versões foram resolvidas pelo proxy Go, sem trocar as bibliotecas de persistência, validação ou templates. `go mod tidy` removeu Fiber v2, o minificador descontinuado no bootstrap e dependências transitivas sem uso. `go.mod` e `go.sum` registram o conjunto completo. O mínimo passou a **Go 1.26.0**, exigido por `golang.org/x/crypto v0.56.0`; use uma versão com os patches atuais da série suportada.

## Avisos encontrados antes

Os avisos alcançáveis incluíam falhas de disponibilidade no parser de formulários e roteamento Fiber, geração previsível de UUID em caso de falha da fonte aleatória, parsing HTML, normalização Unicode, sanitização SQL do pgx e decodificação de imagens TIFF/BMP. Os identificadores, módulos, versões e correções informadas pela base oficial estão no [relatório estruturado](dependency-audit.json).

Depois da atualização permanece somente [GO-2026-5932](https://pkg.go.dev/vuln/GO-2026-5932), sobre o pacote OpenPGP descontinuado de `golang.org/x/crypto`. O projeto usa bcrypt e não importa OpenPGP; o scanner não identificou esse pacote no código compilado. O aviso não possui versão corrigida indicada pela base.

## Reprodução e alcance

```sh
go run golang.org/x/vuln/cmd/govulncheck@v1.7.0 -json ./... > govulncheck.json
```

As duas análises usaram `govulncheck v1.7.0`, a base primária `https://vuln.go.dev` atualizada em `2026-09-02T19:12:04Z` e o toolchain disponível `go1.27.0-X:nodwarf5`, em Linux/amd64. Portanto, o resultado da biblioteca padrão corresponde a esse toolchain, e não constitui validação da versão mínima Go 1.26.0. A CI deve repetir a análise com seu toolchain de implantação.

Este relatório cobre módulos Go. Não audita automaticamente os arquivos JavaScript/CSS vendorizados em `web/plugins` e `web/dist`.
