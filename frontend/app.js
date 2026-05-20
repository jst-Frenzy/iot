const BASE_HOST = window.location.hostname;
const WS_URL = `ws://${BASE_HOST}:8099/ws`;
const API_URL = `http://${BASE_HOST}:8098/api/telemetry`;
const DEVICES_URL = `http://${BASE_HOST}:8098/api/devices`;

/* GLOBAL STATE */
let currentTemperature = 22;
let currentHumidity = 55;
let fanEnabled = false;
let pumpEnabled = false;
let chart = null;
let allData = []; // каждый элемент: { timestamp: Date, temp: number, hum: number }

/* HELPERS */
function formatTime(date) {
  return date.toLocaleTimeString('ru-RU', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  });
}

function toISOString(date) {
  return date.toISOString();
}

/* ELEMENTS */
const tempValue = document.getElementById('temperatureValue');
const humidityValue = document.getElementById('humidityValue');
const tempIndicator = document.getElementById('tempIndicator');
const humIndicator = document.getElementById('humIndicator');
const temperatureState = document.getElementById('temperatureState');
const humidityState = document.getElementById('humidityState');
const fanStatusBadge = document.getElementById('fanStatus');
const pumpStatusBadge = document.getElementById('pumpStatus');
const fanStateTitle = document.getElementById('fanStateTitle');
const fanStateDesc = document.getElementById('fanStateDesc');
const pumpStateTitle = document.getElementById('pumpStateTitle');
const pumpStateDesc = document.getElementById('pumpStateDesc');
const fanBlades = document.getElementById('fanBlades');
const waterLevel = document.getElementById('waterLevel');
const eventFeed = document.getElementById('eventFeed');
const serverStatusText = document.getElementById('serverStatusText');
const devicesCount = document.getElementById('devicesCount');
const miniTemp = document.getElementById('miniTemp');
const miniHum = document.getElementById('miniHum');
const miniFan = document.getElementById('miniFan');
const miniPump = document.getElementById('miniPump');

/* CLOCK */
function updateClock() {
  const clock = document.getElementById('clock');
  if (clock) clock.textContent = new Date().toLocaleTimeString('ru-RU');
}
setInterval(updateClock, 1000);
updateClock();

/* EVENTS */
function addEvent(text) {
  if (!eventFeed) return;
  const div = document.createElement('div');
  div.className = 'event';
  div.innerHTML = `<div class="event-time">${new Date().toLocaleTimeString('ru-RU')}</div><div>${text}</div>`;
  eventFeed.prepend(div);
  while (eventFeed.children.length > 40) eventFeed.removeChild(eventFeed.lastChild);
}

/* MINI INFO (с цветами для ON/OFF) */
function updateMiniInfo() {
  if (miniTemp) miniTemp.innerText = `${Math.round(currentTemperature)}°C`; // целое число
  if (miniHum) miniHum.innerText = `${Math.round(currentHumidity)}%`;
  
  if (miniFan) {
    miniFan.innerText = fanEnabled ? 'ON' : 'OFF';
    miniFan.className = fanEnabled ? 'info-value on' : 'info-value off';
  }
  if (miniPump) {
    miniPump.innerText = pumpEnabled ? 'ON' : 'OFF';
    miniPump.className = pumpEnabled ? 'info-value on' : 'info-value off';
  }
}

/* SENSOR UPDATE */
function updateSensors() {
  // температура – целое число без точки
  if (tempValue) tempValue.innerHTML = `${Math.round(currentTemperature)}<span>°C</span>`;
  if (humidityValue) humidityValue.innerHTML = `${Math.round(currentHumidity)}<span>%</span>`;

  const tempPercent = Math.max(0, Math.min(100, (currentTemperature / 40) * 100));
  const humPercent = Math.max(0, Math.min(100, currentHumidity));
  if (tempIndicator) tempIndicator.style.width = `${tempPercent}%`;
  if (humIndicator) humIndicator.style.width = `${humPercent}%`;

  if (temperatureState) {
    if (currentTemperature >= 28) temperatureState.innerText = 'ВЫШЕ НОРМЫ 🔥';
    else if (currentTemperature <= 18) temperatureState.innerText = 'ПОНИЖЕННАЯ ❄️';
    else temperatureState.innerText = 'НОРМА ✅';
  }
  if (humidityState) {
    if (currentHumidity < 40) humidityState.innerText = 'СУХО 💧';
    else if (currentHumidity > 70) humidityState.innerText = 'ВЫСОКАЯ 🌊';
    else humidityState.innerText = 'НОРМА ✅';
  }
  updateMiniInfo();
  updateAutomation();
}

/* AUTOMATION */
function updateAutomation() {
  let changed = false;
  if (currentTemperature > 28 && !fanEnabled) {
    fanEnabled = true;
    changed = true;
    addEvent(`🌀 Вентиляция включена • Температура: ${Math.round(currentTemperature)}°C`);
  } else if (currentTemperature < 18 && fanEnabled) {
    fanEnabled = false;
    changed = true;
    addEvent(`🛑 Вентиляция отключена • Температура: ${Math.round(currentTemperature)}°C`);
  }
  if (currentHumidity < 40 && !pumpEnabled) {
    pumpEnabled = true;
    changed = true;
    addEvent(`💦 Полив включен • Влажность почвы: ${Math.round(currentHumidity)}%`);
  } else if (currentHumidity > 70 && pumpEnabled) {
    pumpEnabled = false;
    changed = true;
    addEvent(`🛑 Полив отключен • Влажность почвы: ${Math.round(currentHumidity)}%`);
  }
  if (changed) renderDevices();
}

/* DEVICES UI */
function renderDevices() {
  if (fanEnabled) {
    fanStatusBadge.className = 'device-badge on';
    fanStatusBadge.innerText = 'ON';
    fanStateTitle.innerText = 'Охлаждение';
    fanStateDesc.innerText = 'Отключится при < 18°C';
    fanBlades.classList.add('fan-spin');
  } else {
    fanStatusBadge.className = 'device-badge off';
    fanStatusBadge.innerText = 'OFF';
    fanStateTitle.innerText = 'Ожидание';
    fanStateDesc.innerText = 'Включится при > 28°C';
    fanBlades.classList.remove('fan-spin');
  }
  if (pumpEnabled) {
    pumpStatusBadge.className = 'device-badge on';
    pumpStatusBadge.innerText = 'ON';
    pumpStateTitle.innerText = 'Полив активен';
    pumpStateDesc.innerText = 'Отключится при > 70%';
    waterLevel.classList.add('pump-pulse');
  } else {
    pumpStatusBadge.className = 'device-badge off';
    pumpStatusBadge.innerText = 'OFF';
    pumpStateTitle.innerText = 'Ожидание';
    pumpStateDesc.innerText = 'Включится при < 40%';
    waterLevel.classList.remove('pump-pulse');
  }
  updateMiniInfo();
}

/* CHART INIT */
function initChart() {
  const ctx = document.getElementById('telemetryChart');
  if (!ctx) return;
  chart = new Chart(ctx, {
    type: 'line',
    data: { labels: [], datasets: [
      { label: 'Температура °C', data: [], borderColor: '#ff5f5f', backgroundColor: 'rgba(255,95,95,0.14)', fill: true, tension: 0.35, pointRadius: 2, borderWidth: 2 },
      { label: 'Влажность %', data: [], borderColor: '#38bdf8', backgroundColor: 'rgba(56,189,248,0.14)', fill: true, tension: 0.35, pointRadius: 2, borderWidth: 2 }
    ] },
    options: {
      responsive: true, maintainAspectRatio: false,
      interaction: { mode: 'index', intersect: false },
      plugins: { legend: { labels: { color: '#d9ecff', font: { family: 'Inter' } } } },
      scales: {
        x: { ticks: { color: '#8aa6c0' }, grid: { color: 'rgba(255,255,255,0.04)' }, type: 'category' },
        y: { ticks: { color: '#8aa6c0' }, grid: { color: 'rgba(255,255,255,0.04)' } }
      }
    }
  });
}

let currentZoomHours = 1;
function updateChartZoom() {
  if (!chart || allData.length === 0) return;
  const now = new Date();
  const from = new Date(now.getTime() - currentZoomHours * 60 * 60 * 1000);
  const filtered = allData.filter(d => d.timestamp >= from && d.timestamp <= now);
  if (filtered.length === 0) return;
  chart.data.labels = filtered.map(d => formatTime(d.timestamp));
  chart.data.datasets[0].data = filtered.map(d => d.temp);
  chart.data.datasets[1].data = filtered.map(d => d.hum);
  chart.update();
}

function addDataPoint(timestamp, temperature, humidity) {
  allData.push({ timestamp, temp: temperature, hum: humidity });
  allData.sort((a, b) => a.timestamp - b.timestamp);
  updateChartZoom();
}

/* ЗАГРУЗКА ИСТОРИИ (один раз за 24 часа) */
async function fetchDeviceData(deviceName, from, to) {
  const url = `${API_URL}?device_name=${encodeURIComponent(deviceName)}&from=${toISOString(from)}&to=${toISOString(to)}`;
  try {
    const resp = await fetch(url);
    if (!resp.ok) return [];
    const data = await resp.json();
    return Array.isArray(data) ? data : [];
  } catch (e) {
    console.error(`Ошибка загрузки ${deviceName}:`, e);
    return [];
  }
}

async function loadFullHistory(hours = 24) {
  const now = new Date();
  const from = new Date(now.getTime() - hours * 60 * 60 * 1000);
  const to = now;
  const [tempRaw, humRaw] = await Promise.all([
    fetchDeviceData('temperature', from, to),
    fetchDeviceData('humidity', from, to)
  ]);
  const tempMap = new Map();
  const humMap = new Map();
  tempRaw.forEach(item => { if (item.timestamp && typeof item.value === 'number') tempMap.set(item.timestamp, item.value); });
  humRaw.forEach(item => { if (item.timestamp && typeof item.value === 'number') humMap.set(item.timestamp, item.value); });
  const allTimestamps = new Set([...tempMap.keys(), ...humMap.keys()]);
  const sortedTimestamps = Array.from(allTimestamps).sort((a,b) => new Date(a) - new Date(b));
  allData = [];
  for (const ts of sortedTimestamps) {
    const date = new Date(ts);
    const temp = tempMap.get(ts);
    const hum = humMap.get(ts);
    if (temp !== undefined && hum !== undefined) {
      allData.push({ timestamp: date, temp, hum });
    }
  }
  updateChartZoom();
  addEvent(`📈 Загружена история за последние ${hours} ч (${allData.length} точек)`);
}

/* ЗАГРУЗКА УСТРОЙСТВ (счётчик) */
async function loadDevices() {
  try {
    const resp = await fetch(DEVICES_URL);
    if (!resp.ok) return;
    const data = await resp.json();
    if (devicesCount && data.devices) devicesCount.textContent = data.devices.length;
  } catch (e) { console.error('Devices error:', e); }
}

/* КНОПКИ ПЕРИОДОВ (только масштаб, без логов) */
function initPeriodButtons() {
  const buttons = document.querySelectorAll('.period-btn');
  const setActive = (activeBtn) => {
    buttons.forEach(btn => btn.classList.remove('active'));
    activeBtn.classList.add('active');
  };
  buttons.forEach(btn => {
    btn.addEventListener('click', () => {
      const hours = parseFloat(btn.dataset.hours);
      if (isNaN(hours) || hours <= 0) return;
      currentZoomHours = hours;
      setActive(btn);
      updateChartZoom();
      // событие в лог НЕ добавляем
    });
  });
}

/* WEBSOCKET */
let ws = null;
function connectWebSocket() {
  ws = new WebSocket(WS_URL);
  ws.onopen = () => {
    serverStatusText.textContent = 'Сервер подключен';
    addEvent('✅ WebSocket подключен');
  };
  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      const temp = Number(data.Temperature);
      const wet = Number(data.Wet);
      if (!isNaN(temp)) currentTemperature = temp;
      if (!isNaN(wet)) currentHumidity = wet;
      updateSensors();
      addDataPoint(new Date(), currentTemperature, currentHumidity);
    } catch (e) { console.error('WS parse error:', e); }
  };
  ws.onerror = () => { serverStatusText.textContent = 'Ошибка сервера'; };
  ws.onclose = () => {
    serverStatusText.textContent = 'Переподключение...';
    setTimeout(connectWebSocket, 3000);
  };
}

/* INIT */
window.addEventListener('DOMContentLoaded', async () => {
  initChart();
  updateSensors();
  renderDevices();
  loadDevices();
  await loadFullHistory(24);
  initPeriodButtons();
  connectWebSocket();
  addEvent('🌱 Система автоматической теплицы запущена');
});