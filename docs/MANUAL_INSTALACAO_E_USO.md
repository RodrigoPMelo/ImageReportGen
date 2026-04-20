# ImageReportGen — Manual de Instalação e Uso

Bem-vindo. Este guia explica, passo a passo, como instalar e utilizar o **ImageReportGen** no seu computador, **sem precisar de conhecimentos de informática avançados**.

---

## 1. Visão geral (O que é o ImageReportGen?)

O **ImageReportGen** é um programa para **Windows** que ajuda a criar **relatórios com muitas fotografias** de forma rápida e organizada. Em vez de colar imagem a imagem num documento Word e alinhar tudo à mão, o programa **faz essa montagem por si**, respeitando o **aspeto oficial** do documento da sua empresa (cabeçalho, rodapé e marca de água).

A dor que resolve é simples: **poupar tempo** na preparação de relatórios fotográficos (por exemplo, visitas a obra, inspeções ou álbuns internos), mantendo um resultado **profissional** e **consistente** com o molde que já utiliza no Word.

> **Importante:** O programa **funciona sem ligação à Internet**. As suas fotografias e o seu documento-molde **permanecem no seu computador** — **não são enviados** para serviços na nuvem nem para terceiros.

---

## 2. Instalação (simples e direta)

Não existe um “instalador” com muitos ecrãs. O que recebe é uma **pasta** já pronta para utilizar.

1. **Descompacte** o ficheiro que lhe foi entregue (se vier num arquivo `.zip`), com o botão direito do rato → **Extrair tudo…** (ou equivalente).
2. Arraste a pasta resultante para um sítio **fácil de encontrar**, por exemplo:
   - **Documentos**, ou  
   - **Área de trabalho**.
3. Abra essa pasta e faça **duplo clique** no ficheiro **`ImageReportGen.exe`** para iniciar o programa.

> **Dica:** Se o Windows mostrar um aviso de segurança na primeira vez, veja a secção **“O programa não abre ou o Windows bloqueia”** mais abaixo, no separador de perguntas frequentes.

### (Opcional) Criar um atalho na Área de trabalho

1. Localize o ficheiro **`ImageReportGen.exe`** dentro da pasta da aplicação.
2. Clique com o **botão direito** do rato sobre o ficheiro.
3. Escolha **Mostrar mais opções** (se aparecer) e depois **Enviar para** → **Área de trabalho (criar atalho)**.  
   No dia seguinte, pode abrir o programa diretamente a partir do atalho.

---

## 3. Passo a passo de uso (o fluxo de trabalho)

Quando o programa abrir, verá uma janela com um título, botões e uma **zona central** para arrastar ficheiros. Siga esta ordem:

### Passo 1 — Preparar e escolher o “molde” (documento Word)

1. No Microsoft Word, prepare o documento que serve de **molde** da empresa: página com **cabeçalho**, **rodapé** e **marca de água** como costuma fazer para relatórios oficiais.
2. Guarde esse documento como **`.docx`** (formato normal do Word atual).
3. No **ImageReportGen**, clique no botão **“Selecionar Modelo .docx”**.
4. Na janela que o Windows abre, escolha o ficheiro do molde e confirme.

**O que o programa espera do molde**

- Precisa de um documento com o **aspeto da página** já definido (logótipo, cores, cabeçalho, rodapé, marca de água, etc.).
- **Não é obrigatório** colocar tabelas vazias ou “caixas” para as fotos — o programa trata da disposição das imagens ao gerar o relatório.

### Passo 2 — Adicionar as fotografias

1. Selecione no Explorador de ficheiros as imagens que quer incluir (pode ser uma a uma ou várias de seguida).
2. **Arraste** esses ficheiros para a **zona central** da janela (a área com o texto a indicar que pode largar imagens ou um `.zip`).

**Tipos de ficheiro aceites**

- Imagens: **`.jpg`**, **`.jpeg`** ou **`.png`**.
- Também pode largar um ficheiro **`.zip`** que contenha fotografias com essas extensões.

Se largar outro tipo de ficheiro, o nome pode aparecer na lista **“Ficheiros ignorados”** — isso é normal; o programa só processa o que reconhece como imagem ou ficheiro compactado com imagens.

### Passo 3 — Gerar o relatório

1. Quando estiver pronto, clique em **“Gerar Relatório”**.
2. Observe a **linha de estado** (texto ao lado dos botões): lá aparece se está a processar, se correu bem ou se surgiu algum problema.
3. Enquanto o programa estiver a trabalhar, os botões podem ficar **temporariamente indisponíveis** — espere até a mensagem voltar a **“Pronto”** ou mostrar o resultado.

**Onde fica o relatório final**

- O programa cria um ficheiro com o nome **`relatorio_gerado.docx`**.
- Esse ficheiro é gravado na **pasta de trabalho** em que o programa está a correr (por exemplo, por vezes a mesma pasta de onde abriu o programa ou a pasta onde está o executável — depende de como foi iniciado).
- Se não vir o ficheiro de imediato, abra o **Explorador de ficheiros**, use a **caixa de pesquisa** e procure por **`relatorio_gerado.docx`**.

> **Nota:** O programa **não pergunta** “onde guardar” nesta versão — o nome do ficheiro é sempre o mesmo. Se já existir um `relatorio_gerado.docx` na mesma pasta, o Word ou o Windows podem pedir para substituir ou sugerir outro nome ao abrir — trate disso como costuma fazer com qualquer documento.

---

## 4. Dicas de ouro e regras do molde

- **Fotos em modo horizontal (“paisagem”):** o programa coloca até **três fotografias por página**, **uma por linha** (três linhas, uma coluna).
- **Fotos em modo vertical (“retrato”):** o programa coloca até **quatro fotografias por página**, numa **grelha de duas colunas e duas linhas** (2×2).

> **Lembrete:** Se tiver **muitas** fotografias, a geração pode demorar alguns instantes. **Não feche a janela** do programa enquanto a mensagem de estado indicar que está a processar ou enquanto os botões estiverem indisponíveis.

---

## 5. Solução de problemas (FAQ básico)

### “O programa não abre ou o Windows bloqueia.”

Na primeira execução, o Windows pode mostrar um aviso de que o ficheiro não é “frequentemente transferido”. Para uma aplicação **interna** da sua organização:

1. Clique em **Mais informações** (se aparecer).
2. Depois clique em **Executar na mesma** ou **Sim**, conforme o ecrã.

Se a sua empresa tiver regras de segurança rígidas, peça ajuda ao **departamento de informática** — eles sabem como autorizar programas internos.

### “O relatório final ficou sem o cabeçalho (ou sem o rodapé / marca de água).”

O programa **usa o molde** que escolheu. Se o cabeçalho, rodapé ou marca de água **não estiverem definidos no próprio Word** no documento-molde, o resultado também não os terá “por magia”. Abra o molde no Word, confirme que tudo aparece em **Visualizar impressão** e volte a guardar o `.docx` antes de o selecionar de novo no **ImageReportGen**.

### “Coloquei as fotos, mas deu erro de formato.”

Confirme que as imagens são **`.jpg`**, **`.jpeg`** ou **`.png`**. Outros formatos (por exemplo, ficheiros de telemóvel em formatos proprietários) precisam de ser **convertidos** para um destes formatos antes de os largar no programa.

### “Onde está o ficheiro `relatorio_gerado.docx`?”

Use a **pesquisa** do Explorador de ficheiros pelo nome **`relatorio_gerado.docx`**. Se gerar vários relatórios no mesmo dia, considere **renomear** ou **mover** o ficheiro para outra pasta logo após cada geração, para não sobrescrever o anterior por engano.

---

*Obrigado por utilizar o ImageReportGen. Em caso de dúvida, fale com a pessoa que lhe entregou o programa ou com o suporte interno da sua empresa.*
