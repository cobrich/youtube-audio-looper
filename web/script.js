const form = document.querySelector("#loop-form");
const submitButton = document.querySelector("#submit-button");
const result = document.querySelector("#result");

let currentAudioUrl = null;

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

function showError(message) {
  clearAudioUrl();
  result.className = "result result-error";
  result.innerHTML = `
    <h2>Could not create audio</h2>
    <p>${escapeHtml(message)}</p>
  `;
}

function showAudio(blob) {
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

  result.className = "result result-loading";
  result.innerHTML = `
    <div class="spinner" aria-hidden="true"></div>
    <p>Downloading, cutting, and looping audio...</p>
  `;

  setLoading(true);

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
    showAudio(blob);
  } catch (error) {
    showError(error.message || "Network error");
  } finally {
    setLoading(false);
  }
});
