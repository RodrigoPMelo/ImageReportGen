import "./style.css";
import "./app.css";

import { ProcessUploads, RunGeneration, SelectTemplate } from "../wailsjs/go/main/App";

type UploadState = {
  templateName: string;
  uploadedFiles: string[];
  ignoredFiles: string[];
  totalUploads: number;
  status: string;
  busy: boolean;
};

const state: UploadState = {
  templateName: "Nenhum modelo selecionado",
  uploadedFiles: [],
  ignoredFiles: [],
  totalUploads: 0,
  status: "Pronto",
  busy: false,
};

const appRoot = document.querySelector("#app");
if (!appRoot) {
  throw new Error("App root não encontrado");
}

appRoot.innerHTML = `
  <main class="layout">
    <h1>ImageReportGen</h1>
    <section class="panel">
      <div class="row">
        <button id="btn-template" class="btn">Selecionar Modelo .docx</button>
        <span id="template-name" class="muted"></span>
      </div>
      <div id="dropzone" class="dropzone">
        Arraste imagens (.png, .jpg, .jpeg) ou ficheiros .zip para aqui
      </div>
      <div class="row">
        <button id="btn-generate" class="btn btn-primary">Gerar Relatório</button>
        <span id="status" class="status"></span>
      </div>
      <div class="meta">Total de uploads: <strong id="total-uploads">0</strong></div>
      <h3>Ficheiros carregados</h3>
      <ul id="uploaded-list" class="list"></ul>
      <h3>Ficheiros ignorados</h3>
      <ul id="ignored-list" class="list"></ul>
    </section>
  </main>
`;

const btnTemplate = document.getElementById("btn-template") as HTMLButtonElement;
const btnGenerate = document.getElementById("btn-generate") as HTMLButtonElement;
const templateName = document.getElementById("template-name") as HTMLSpanElement;
const statusEl = document.getElementById("status") as HTMLSpanElement;
const dropzone = document.getElementById("dropzone") as HTMLDivElement;
const uploadedList = document.getElementById("uploaded-list") as HTMLUListElement;
const ignoredList = document.getElementById("ignored-list") as HTMLUListElement;
const totalUploads = document.getElementById("total-uploads") as HTMLElement;

function render(): void {
  templateName.textContent = state.templateName;
  statusEl.textContent = state.status;
  btnTemplate.disabled = state.busy;
  btnGenerate.disabled = state.busy;
  totalUploads.textContent = String(state.totalUploads);

  uploadedList.innerHTML = state.uploadedFiles.map((item) => `<li>${item}</li>`).join("");
  ignoredList.innerHTML = state.ignoredFiles.map((item) => `<li>${item}</li>`).join("");
}

function displayName(path: string): string {
  const parts = path.split(/[\\/]/);
  return parts[parts.length - 1] || path;
}

function extractPaths(evt: DragEvent): string[] {
  const files = evt.dataTransfer?.files;
  if (!files) {
    return [];
  }
  const paths: string[] = [];
  for (const file of Array.from(files)) {
    const maybePath = (file as File & { path?: string }).path;
    if (maybePath && maybePath.length > 0) {
      paths.push(maybePath);
    }
  }
  return paths;
}

btnTemplate.addEventListener("click", async () => {
  state.busy = true;
  state.status = "A selecionar template...";
  render();
  try {
    const selected = await SelectTemplate();
    if (selected) {
      state.templateName = selected;
      state.status = "Template selecionado";
    } else {
      state.status = "Seleção cancelada";
    }
  } catch (error) {
    state.status = `Erro ao selecionar template: ${String(error)}`;
  } finally {
    state.busy = false;
    render();
  }
});

dropzone.addEventListener("dragover", (evt) => {
  evt.preventDefault();
  dropzone.classList.add("active");
});

dropzone.addEventListener("dragleave", () => {
  dropzone.classList.remove("active");
});

dropzone.addEventListener("drop", async (evt) => {
  evt.preventDefault();
  dropzone.classList.remove("active");

  const paths = extractPaths(evt);
  if (paths.length === 0) {
    state.status = "Drop sem paths locais disponíveis";
    render();
    return;
  }

  state.busy = true;
  state.status = "A processar uploads...";
  render();
  try {
    const result = await ProcessUploads(paths);
    state.uploadedFiles.push(...result.added.map(displayName));
    state.ignoredFiles.push(...result.ignored.map(displayName));
    state.totalUploads = result.totalUploads;
    state.status = `Uploads processados: +${result.added.length}`;
  } catch (error) {
    state.status = `Erro no processamento: ${String(error)}`;
  } finally {
    state.busy = false;
    render();
  }
});

btnGenerate.addEventListener("click", async () => {
  state.busy = true;
  state.status = "A gerar relatório...";
  render();
  try {
    const result = await RunGeneration();
    state.status = `Relatório gerado: ${displayName(result.outputPath)} (${result.totalImages} imagens)`;
  } catch (error) {
    state.status = `Erro na geração: ${String(error)}`;
  } finally {
    state.busy = false;
    render();
  }
});

render();
