const form = document.querySelector("#loop-form");
const submitButton = document.querySelector("#submit-button");
const result = document.querySelector("#result");

let currentAudioUrl = null;
let progressTimer = null;
let progressStartedAt = null;

const legacyHistoryKey = "youtube-audio-looper-history";
const requestHistoryKey = "youtube-audio-looper-request-history";
const lastResultKey = "youtube-audio-looper-last-result";
const maxHistoryItems = 10;

const progressSteps = [
  { at: 0, label: "Preparing request" },
  { at: 8, label: "Connecting to YouTube" },
  { at: 22, label: "Downloading selected segment" },
  { at: 68, label: "Creating audio loop" },
  { at: 88, label: "Finalizing MP3" },
];

function setLoading(isLoading) {
  submitButton.disabled = isLoading;
  submitButton.querySelector("span:last-child").textContent = isLoading
    ? "Creating..."
    : "Create MP3";
}

function clearAudioUrl() {
  if (currentAudioUrl) {
    URL.revokeObjectURL(currentAudioUrl);
    currentAudioUrl = null;
  }
}

function readJSON(key, fallback) {
  try {
    return JSON.parse(localStorage.getItem(key) || JSON.stringify(fallback));
  } catch {
    return fallback;
  }
}

function writeJSON(key, value) {
  try {
    localStorage.setItem(key, JSON.stringify(value));
  } catch (error) {
    console.warn("Could not save local data", error);
  }
}

function readRequestHistory() {
  return readJSON(requestHistoryKey, []);
}

function writeRequestHistory(items) {
  writeJSON(requestHistoryKey, items.slice(0, maxHistoryItems));
}

function readLastResult() {
  return readJSON(lastResultKey, null);
}

function blobToDataUrl(blob) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.addEventListener("load", () => resolve(reader.result));
    reader.addEventListener("error", () => reject(reader.error));
    reader.readAsDataURL(blob);
  });
}

function formatDate(value) {
  return new Intl.DateTimeFormat(undefined, {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(value));
}

function getHistoryTitle(item) {
  try {
    const parsedURL = new URL(item.request.youtube_url);
    return parsedURL.searchParams.get("v") || parsedURL.hostname;
  } catch {
    return item.request.youtube_url;
  }
}

function fillFormFromRequest(request) {
  form.elements.youtube_url.value = request.youtube_url;
  form.elements.start.value = request.start;
  form.elements.end.value = request.end;
  form.elements.duration.value = request.duration;
}

function saveRequestHistoryItem(request) {
  const items = readRequestHistory();
  const dedupedItems = items.filter((item) => {
    return JSON.stringify(item.request) !== JSON.stringify(request);
  });
  const nextItems = [
    {
      id: String(Date.now()),
      createdAt: new Date().toISOString(),
      request,
    },
    ...dedupedItems,
  ];

  writeRequestHistory(nextItems);
}

async function saveLastResult(blob, request) {
  const dataUrl = await blobToDataUrl(blob);
  writeJSON(lastResultKey, {
    createdAt: new Date().toISOString(),
    request,
    dataUrl,
  });
}

function migrateLegacyHistory() {
  if (localStorage.getItem(requestHistoryKey)) {
    return;
  }

  const legacyItems = readJSON(legacyHistoryKey, []);
  if (legacyItems.length === 0) {
    return;
  }

  writeRequestHistory(
    legacyItems.map((item) => ({
      id: item.id || String(Date.now()),
      createdAt: item.createdAt || new Date().toISOString(),
      request: item.request,
    })),
  );

  const latestLegacyItem = legacyItems[0];
  if (latestLegacyItem?.dataUrl) {
    writeJSON(lastResultKey, {
      createdAt: latestLegacyItem.createdAt || new Date().toISOString(),
      request: latestLegacyItem.request,
      dataUrl: latestLegacyItem.dataUrl,
    });
  }
}

function stopProgress() {
  if (progressTimer) {
    clearInterval(progressTimer);
    progressTimer = null;
  }
}

function getProgressLabel(progress) {
  return progressSteps.reduce((current, step) => {
    return progress >= step.at ? step.label : current;
  }, progressSteps[0].label);
}

function updateProgress(progress) {
  const progressValue = Math.min(Math.round(progress), progress >= 100 ? 100 : 99);
  const track = result.querySelector(".progress-track");
  const bar = result.querySelector("#progress-bar");
  const percent = result.querySelector("#progress-percent");
  const status = result.querySelector("#progress-status");
  const elapsed = result.querySelector("#progress-elapsed");

  if (!track || !bar || !percent || !status || !elapsed || !progressStartedAt) {
    return;
  }

  const elapsedSeconds = Math.floor((Date.now() - progressStartedAt) / 1000);
  track.setAttribute("aria-valuenow", String(progressValue));
  bar.style.width = `${progressValue}%`;
  percent.textContent = `${progressValue}%`;
  status.textContent = getProgressLabel(progressValue);
  elapsed.textContent = `${elapsedSeconds}s elapsed`;
}

function startProgress() {
  stopProgress();
  progressStartedAt = Date.now();

  result.className = "result result-loading";
  result.innerHTML = `
    <div class="progress-header">
      <h2>Creating audio</h2>
      <span id="progress-percent">0%</span>
    </div>
    <div class="progress-track" role="progressbar" aria-valuemin="0" aria-valuemax="100">
      <div id="progress-bar" class="progress-bar"></div>
    </div>
    <div class="progress-meta">
      <p id="progress-status">Preparing request</p>
      <span id="progress-elapsed">0s elapsed</span>
    </div>
  `;

  updateProgress(2);

  progressTimer = setInterval(() => {
    const elapsedSeconds = (Date.now() - progressStartedAt) / 1000;
    const progress = Math.min(92, 6 + elapsedSeconds * 2.2);
    updateProgress(progress);
  }, 500);
}

function showError(message) {
  stopProgress();
  clearAudioUrl();
  result.className = "result result-error";
  result.innerHTML = `
    <h2>Could not create audio</h2>
    <p>${escapeHtml(message)}</p>
  `;
}

function showEmptyState() {
  clearAudioUrl();
  result.className = "result";
  result.innerHTML = `
    <div class="empty-state">
      <p>Result audio will appear here after processing.</p>
    </div>
  `;
}

function showAudioSource(audioSource, request) {
  stopProgress();

  result.className = "result result-success";
  result.innerHTML = `
    <h2>Looped audio is ready</h2>
    <audio controls src="${audioSource}"></audio>
    <a class="download-link" href="${audioSource}" download="looped-audio.mp3">
      Download MP3
    </a>
    ${request ? `<p class="result-details">${escapeHtml(request.start)}-${escapeHtml(request.end)} to ${escapeHtml(request.duration)}</p>` : ""}
    ${renderHistory()}
  `;
}

function showAudio(blob, request) {
  clearAudioUrl();
  currentAudioUrl = URL.createObjectURL(blob);
  showAudioSource(currentAudioUrl, request);
}

function showStoredResult(item) {
  clearAudioUrl();
  fillFormFromRequest(item.request);
  showAudioSource(item.dataUrl, item.request);
}

function showRequestHistoryItem(item) {
  fillFormFromRequest(item.request);
}

function renderHistory() {
  const items = readRequestHistory();

  if (items.length === 0) {
    return "";
  }

  return `
    <section class="history">
      <div class="history-header">
        <h3>Recent requests</h3>
        <button class="history-clear" type="button" data-clear-history>Clear</button>
      </div>
      <div class="history-list">
        ${items
          .map(
            (item) => `
              <article class="history-item">
                <button class="history-play" type="button" data-history-id="${escapeHtml(item.id)}">
                  <span>${escapeHtml(getHistoryTitle(item))}</span>
                  <small>${escapeHtml(item.request.start)}-${escapeHtml(item.request.end)} - ${escapeHtml(item.request.duration)}</small>
                  <small>${escapeHtml(formatDate(item.createdAt))}</small>
                </button>
              </article>
            `,
          )
          .join("")}
      </div>
    </section>
  `;
}

function escapeHtml(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

async function getErrorMessage(response) {
  const contentType = response.headers.get("content-type") || "";

  if (contentType.includes("application/json")) {
    const data = await response.json();
    return data.error || "Request failed";
  }

  const text = await response.text();
  return text || "Request failed";
}

form.addEventListener("submit", async (event) => {
  event.preventDefault();

  const formData = new FormData(form);
  const payload = {
    youtube_url: formData.get("youtube_url").trim(),
    start: formData.get("start").trim(),
    end: formData.get("end").trim(),
    duration: formData.get("duration").trim(),
  };

  setLoading(true);
  startProgress();

  try {
    const response = await fetch("/api/v1/audio/loop", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      showError(await getErrorMessage(response));
      return;
    }

    const blob = await response.blob();
    updateProgress(100);
    saveRequestHistoryItem(payload);
    await saveLastResult(blob, payload);
    showAudio(blob, payload);
  } catch (error) {
    showError(error.message || "Network error");
  } finally {
    setLoading(false);
  }
});

result.addEventListener("click", (event) => {
  const historyButton = event.target.closest("[data-history-id]");
  const clearButton = event.target.closest("[data-clear-history]");

  if (historyButton) {
    const item = readRequestHistory().find((historyItem) => historyItem.id === historyButton.dataset.historyId);
    if (item) {
      showRequestHistoryItem(item);
    }
  }

  if (clearButton) {
    writeRequestHistory([]);
    const latestResult = readLastResult();
    if (latestResult) {
      showStoredResult(latestResult);
    } else {
      showEmptyState();
    }
  }
});

migrateLegacyHistory();

const latestResult = readLastResult();

if (latestResult) {
  showStoredResult(latestResult);
}
