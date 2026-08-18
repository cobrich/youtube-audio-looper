const form = document.querySelector("#loop-form");
const submitButton = document.querySelector("#submit-button");
const result = document.querySelector("#result");

let currentAudioUrl = null;
let progressTimer = null;
let progressStartedAt = null;

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

function showAudio(blob) {
  stopProgress();
  clearAudioUrl();
  currentAudioUrl = URL.createObjectURL(blob);

  result.className = "result result-success";
  result.innerHTML = `
    <h2>Looped audio is ready</h2>
    <audio controls src="${currentAudioUrl}"></audio>
    <a class="download-link" href="${currentAudioUrl}" download="looped-audio.mp3">
      Download MP3
    </a>
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
    showAudio(blob);
  } catch (error) {
    showError(error.message || "Network error");
  } finally {
    setLoading(false);
  }
});
