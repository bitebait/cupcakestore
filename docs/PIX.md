# Pix local com confirmação manual

A loja gera QR Code e Pix Copia e Cola no próprio backend. Não envia dados do pedido a sites geradores, não redireciona o cliente para pagar fora da loja e não exige uma conta em gateway. O código contém instruções para o aplicativo bancário; a geração ou exibição desse código não comprova que dinheiro foi recebido.

## Configuração

No painel de configuração, informe o tipo de chave, a chave Pix registrada da loja, o nome do recebedor e a cidade. Verifique a chave e sua titularidade no aplicativo do banco antes de ativar o Pix. A aplicação valida formato e dígitos verificadores, mas não consulta o DICT.

Aceitamos CPF, CNPJ numérico ou alfanumérico, e-mail, celular no formato internacional com `+` e chave aleatória. Máscaras usuais de CPF/CNPJ/celular são normalizadas; e-mail fica em minúsculas. O código mantém valor com duas casas decimais e uma referência estável por pedido. Nome/cidade têm acentos normalizados para ASCII e são limitados a 25/15 caracteres; texto vazio ou não representável é rejeitado. Essas regras seguem os campos do [BR Code](https://www.bcb.gov.br/content/estabilidadefinanceira/spb_docs/ManualBRCode.pdf) e do [Manual de Iniciação do Pix](https://www.bcb.gov.br/content/estabilidadefinanceira/pix/Regulamento_Pix/II_ManualdePadroesparaIniciacaodoPix.pdf). O CNPJ alfanumérico usa o [cálculo de dígitos da Receita Federal](https://www.gov.br/receitafederal/pt-br/centrais-de-conteudo/publicacoes/documentos-tecnicos/cnpj/manual-dv-cnpj.pdf).

O banco identifica o recebedor pela chave registrada. O cliente deve conferir o beneficiário e o valor na tela do banco antes de autorizar a transferência.

## Compra e atendimento

O cliente escolhe entrega ou retirada e confirma Pix. A página do pedido mostra QR Code, código para copiar e valor. O conteúdo permanece igual ao reabrir o pedido, mesmo que a configuração ou taxa de entrega da loja tenha mudado depois. Forma de pagamento e entrega ficam fixadas após emitir as instruções.

O pedido aguarda conferência. O administrador verifica o crédito no extrato do banco, comparando valor e identificador da transação, e só então confirma o recebimento no painel. Um comprovante enviado pelo cliente ou a informação “já paguei” não substitui a conferência do crédito. A confirmação administrativa registra responsável e horário; o andamento da preparação continua pelas etapas disponíveis no pedido.

Emissão e renderização trabalham sem rede. O backend verifica estrutura, checksum, moeda, chave e correspondência de valor/referência antes de apresentar instruções armazenadas. Pedidos legados com código inválido, sem valor fixo ou com referência divergente devem ser tratados pela loja; seus links externos antigos não são reabertos automaticamente.

## Limites e conciliação

Não há confirmação bancária automática, consulta de saldo, devolução ou cancelamento remoto. Cancelar o pedido restitui a reserva de estoque conforme as regras locais, mas não revoga um QR Code que o cliente já salvou. Pagamentos duplicados, feitos após cancelamento ou com valores divergentes precisam ser conciliados e, quando aplicável, devolvidos pela loja através do banco. Não existe expiração bancária imposta por este QR estático.

Se houver necessidade futura de conciliação automática, integre a API do PSP usado pela loja. Isso exige conta habilitada e credenciais configuradas na infraestrutura, sem compartilhá-las no chat ou versioná-las. Uma integração deve autenticar eventos, consultar o pagamento no PSP, conferir valor/recebedor/identificador, deduplicar recebimentos e recuperar notificações perdidas. Como referência técnica, a [API Pix do BCB](https://github.com/bacen/pix-api) define o contrato, e a documentação oficial da Efí exemplifica [cobranças imediatas](https://dev.efipay.com.br/docs/api-pix/cobrancas-imediatas/) e [webhooks autenticados](https://dev.efipay.com.br/docs/api-pix/webhooks/); isso não representa a seleção ou contratação desse provedor.

## Manutenção e testes

`models/pix_key.go` concentra normalização/validação. `services/pix_local.go` monta o payload e renderiza PNG com o encoder [go-qrcode v2.2.5](https://github.com/yeqown/go-qrcode/releases/tag/v2.2.5), sem biblioteca de efeitos gráficos. A referência e os valores pertencem ao pedido, e o QR é uma representação desses dados.

```sh
go test ./models ./services ./controllers ./repositories
```

As regressões verificam o vetor publicado pelo BCB, CRC16, nomes acentuados, CNPJ alfanumérico, chaves inválidas, limites de valor/referência, PNG/borda branca e divergência entre código e pedido. Na validação da migração, um decoder independente ZXing recuperou exatamente os 137 caracteres do vetor oficial a partir do PNG de 342 × 342 pixels; as ferramentas dessa verificação ficaram fora das dependências da aplicação. Testes de software não demonstram registro da chave nem liquidação real; a implantação deve conferir esses dados no banco da loja.
