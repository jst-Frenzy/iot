const WS_URL = 'ws://localhost:8099/ws';
const API_URL = 'http://localhost:8098/api/telemetry';
const DEVICES_URL = 'http://localhost:8098/api/devices';

/* =========================================
   GLOBAL STATE
========================================= */

let currentTemperature = 22;
let currentHumidity = 55;
let fanEnabled = false;
let pumpEnabled = false;
let chart = null;

/* =========================================
   ELEMENTS (СООТВЕТСТВУЕТ ВАШЕМУ HTML)
========================================= */

const tempValue = document.getElementById('temperatureValue');
const humidityValue = document.getElementById('humidityValue');
const tempFill = document.getElementById('tempFill');
const humidityFill = document.getElementById('humidityFill');
const temperatureState = document.getElementById('temperatureState');
const humidityState = document.getElementById('humidityState');
const fanStatusBadge = document.getElementById('fanStatus');
const pumpStatusBadge = document.getElementById('pumpStatus');
const fanStateText = document.getElementById('fanStateText');
const pumpStateText = document.getElementById('pumpStateText');
const fanImage = document.getElementById('fanImage');
const fanGlow = document.getElementById('fanGlow');
const waterDrops = document.getElementById('waterDrops');
const eventFeed = document.getElementById('eventFeed');
const wsStatusSpan = document.getElementById('wsStatus');
const serverStatusText = document.getElementById('serverStatusText');
const devicesCountSpan = document.getElementById('devicesCount');

/* =========================================
   CLOCK
========================================= */

function updateClock() {
    const clock = document.getElementById('clock');
    if (clock) {
        clock.textContent = new Date().toLocaleTimeString('ru-RU');
    }
}
setInterval(updateClock, 1000);
updateClock();

/* =========================================
   EVENTS FEED
========================================= */

function addEvent(text, type = 'info') {
    const div = document.createElement('div');
    div.className = 'event';
    div.innerHTML = `
        <div class="event-time">${new Date().toLocaleTimeString()}</div>
        <div>${text}</div>
    `;
    
    if (eventFeed) {
        eventFeed.prepend(div);
        if (eventFeed.children.length > 40) {
            eventFeed.removeChild(eventFeed.lastChild);
        }
    }
    
    console.log(`[EVENT] ${text}`);
}

/* =========================================
   SENSOR UPDATE
========================================= */

function updateSensors() {
    // Update values
    if (tempValue) {
        tempValue.innerHTML = `${currentTemperature.toFixed(1)}<span style="font-size:24px">°C</span>`;
    }
    if (humidityValue) {
        humidityValue.innerHTML = `${currentHumidity.toFixed(0)}<span style="font-size:24px">%</span>`;
    }
    
    // Thermometer fill
    const tempPercent = Math.min(100, Math.max(0, (currentTemperature / 40) * 100));
    if (tempFill) tempFill.style.height = `${tempPercent}%`;
    if (humidityFill) humidityFill.style.height = `${currentHumidity}%`;
    
    // Temperature status
    if (temperatureState) {
        if (currentTemperature >= 28) {
            temperatureState.className = 'sensor-status red-status';
            temperatureState.innerText = 'ВЫШЕ НОРМЫ 🔥';
        } else if (currentTemperature <= 18) {
            temperatureState.className = 'sensor-status red-status';
            temperatureState.innerText = 'ПОНИЖЕННАЯ ❄️';
        } else {
            temperatureState.className = 'sensor-status red-status';
            temperatureState.innerText = 'НОРМА ✅';
        }
    }
    
    // Humidity status
    if (humidityState) {
        if (currentHumidity < 40) {
            humidityState.className = 'sensor-status blue-status';
            humidityState.innerText = 'НИЗКАЯ 💧';
        } else if (currentHumidity > 70) {
            humidityState.className = 'sensor-status blue-status';
            humidityState.innerText = 'ВЫСОКАЯ 🌊';
        } else {
            humidityState.className = 'sensor-status blue-status';
            humidityState.innerText = 'НОРМА ✅';
        }
    }
    
    updateAutomation();
}

/* =========================================
   AUTOMATION RULES
========================================= */

function updateAutomation() {
    let changed = false;
    
    // FAN: включаем при t > 28, выключаем при t < 18
    if (currentTemperature > 28 && !fanEnabled) {
        fanEnabled = true;
        changed = true;
        addEvent('🌀 ВЕНТИЛЯТОР ВКЛЮЧЕН • высокая температура');
    } else if (currentTemperature < 18 && fanEnabled) {
        fanEnabled = false;
        changed = true;
        addEvent('🛑 ВЕНТИЛЯТОР ВЫКЛЮЧЕН • температура в норме');
    }
    
    // PUMP: включаем при влажности < 40, выключаем при > 70
    if (currentHumidity < 40 && !pumpEnabled) {
        pumpEnabled = true;
        changed = true;
        addEvent('💧 НАСОС ВКЛЮЧЕН • низкая влажность почвы');
    } else if (currentHumidity > 70 && pumpEnabled) {
        pumpEnabled = false;
        changed = true;
        addEvent('🛑 НАСОС ВЫКЛЮЧЕН • влажность в норме');
    }
    
    if (changed) {
        renderDevices();
    }
}

/* =========================================
   DEVICES UI RENDER
========================================= */

function renderDevices() {
    // Fan UI
    if (fanEnabled) {
        if (fanStatusBadge) {
            fanStatusBadge.className = 'status-badge on';
            fanStatusBadge.innerText = 'ON';
        }
        if (fanStateText) fanStateText.innerText = 'Работает';
        if (fanImage) fanImage.classList.add('fan-active');
        if (fanGlow) fanGlow.classList.add('glow-active');
    } else {
        if (fanStatusBadge) {
            fanStatusBadge.className = 'status-badge off';
            fanStatusBadge.innerText = 'OFF';
        }
        if (fanStateText) fanStateText.innerText = 'Ожидание';
        if (fanImage) fanImage.classList.remove('fan-active');
        if (fanGlow) fanGlow.classList.remove('glow-active');
    }
    
    // Pump UI
    if (pumpEnabled) {
        if (pumpStatusBadge) {
            pumpStatusBadge.className = 'status-badge on';
            pumpStatusBadge.innerText = 'ON';
        }
        if (pumpStateText) pumpStateText.innerText = 'Полив';
        if (waterDrops) waterDrops.classList.add('active');
    } else {
        if (pumpStatusBadge) {
            pumpStatusBadge.className = 'status-badge off';
            pumpStatusBadge.innerText = 'OFF';
        }
        if (pumpStateText) pumpStateText.innerText = 'Ожидание';
        if (waterDrops) waterDrops.classList.remove('active');
    }
}

/* =========================================
   CHART INIT
========================================= */

function initChart() {
    const ctx = document.getElementById('telemetryChart');
    if (!ctx) return;
    
    chart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: [],
            datasets: [
                {
                    label: 'Температура °C',
                    data: [],
                    borderColor: '#ff5f5f',
                    backgroundColor: 'rgba(255,95,95,0.1)',
                    fill: true,
                    tension: 0.4,
                    pointRadius: 3
                },
                {
                    label: 'Влажность %',
                    data: [],
                    borderColor: '#38bdf8',
                    backgroundColor: 'rgba(56,189,248,0.1)',
                    fill: true,
                    tension: 0.4,
                    pointRadius: 3
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: {
                legend: {
                    labels: { color: '#d7eaff' }
                }
            },
            scales: {
                x: {
                    ticks: { color: '#7c9ab5' },
                    grid: { color: 'rgba(255,255,255,0.05)' }
                },
                y: {
                    ticks: { color: '#7c9ab5' },
                    grid: { color: 'rgba(255,255,255,0.05)' }
                }
            }
        }
    });
}

/* =========================================
   ADD CHART POINT
========================================= */

function addChartPoint(temperature, humidity) {
    if (!chart) return;
    
    const time = new Date().toLocaleTimeString();
    
    chart.data.labels.push(time);
    chart.data.datasets[0].data.push(temperature);
    chart.data.datasets[1].data.push(humidity);
    
    if (chart.data.labels.length > 30) {
        chart.data.labels.shift();
        chart.data.datasets[0].data.shift();
        chart.data.datasets[1].data.shift();
    }
    
    chart.update();
}

/* =========================================
   LOAD HISTORY
========================================= */

async function loadHistory(hours = 24) {
    try {
        const now = new Date();
        const from = new Date(now.getTime() - hours * 60 * 60 * 1000);
        
        // Можно загрузить историю, если нужно
        console.log(`Loading history for last ${hours} hours`);
        
    } catch (error) {
        console.error('History load error:', error);
    }
}

/* =========================================
   LOAD DEVICES
========================================= */

async function loadDevices() {
    try {
        const response = await fetch(DEVICES_URL);
        if (response.ok) {
            const data = await response.json();
            console.log('Devices:', data.devices);
            if (devicesCountSpan) {
                devicesCountSpan.textContent = data.devices?.length || 4;
            }
            addEvent(`📡 Устройства подключены: ${data.devices?.join(', ') || '4 устройства'}`);
        }
    } catch (error) {
        console.error('Devices load error:', error);
        if (devicesCountSpan) devicesCountSpan.textContent = '4';
    }
}

/* =========================================
   WEBSOCKET
========================================= */

let ws = null;

function connectWebSocket() {
    ws = new WebSocket(WS_URL);
    
    ws.onopen = () => {
        console.log('WebSocket connected');
        addEvent('✅ WebSocket подключен к серверу');
        if (wsStatusSpan) {
            wsStatusSpan.textContent = 'CONNECTED';
            wsStatusSpan.className = 'live-value connected';
        }
        if (serverStatusText) {
            serverStatusText.textContent = 'Сервер подключен';
        }
    };
    
    ws.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            console.log('WS data:', data);
            
            const temp = Number(data.Temperature);
            const wet = Number(data.Wet);
            
            if (!isNaN(temp)) currentTemperature = temp;
            if (!isNaN(wet)) currentHumidity = wet;
            
            updateSensors();
            addChartPoint(currentTemperature, currentHumidity);
            
        } catch (error) {
            console.error('WS parse error:', error);
        }
    };
    
    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        addEvent('❌ Ошибка WebSocket соединения');
        if (wsStatusSpan) {
            wsStatusSpan.textContent = 'ERROR';
            wsStatusSpan.className = 'live-value';
        }
    };
    
    ws.onclose = () => {
        console.log('WebSocket disconnected');
        addEvent('⚠️ Соединение потеряно, переподключение...');
        if (wsStatusSpan) {
            wsStatusSpan.textContent = 'RECONNECTING';
            wsStatusSpan.className = 'live-value';
        }
        setTimeout(connectWebSocket, 3000);
    };
}

/* =========================================
   RANGE BUTTONS
========================================= */

function initRangeButtons() {
    const buttons = document.querySelectorAll('.range-btn');
    buttons.forEach(btn => {
        btn.addEventListener('click', () => {
            buttons.forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            const hours = parseInt(btn.dataset.hours);
            loadHistory(hours);
        });
    });
}

/* =========================================
   MOCK DATA (FALLBACK)
========================================= */

function startMockData() {
    let temp = 22;
    let humidity = 55;
    let increasing = true;
    
    setInterval(() => {
        if (increasing) {
            temp += Math.random() * 1.5;
            if (temp > 35) increasing = false;
        } else {
            temp -= Math.random() * 1.5;
            if (temp < 10) increasing = true;
        }
        
        humidity += (Math.random() - 0.5) * 4;
        humidity = Math.min(85, Math.max(25, humidity));
        
        currentTemperature = Math.round(temp * 10) / 10;
        currentHumidity = Math.round(humidity);
        
        updateSensors();
        addChartPoint(currentTemperature, currentHumidity);
    }, 5000);
}

/* =========================================
   START
========================================= */

// Инициализация
initChart();
updateSensors();
renderDevices();
initRangeButtons();
loadDevices();
loadHistory(24);
connectWebSocket();

// Fallback: если через 5 секунд нет данных, включаем мок
setTimeout(() => {
    if (currentTemperature === 22 && currentHumidity === 55) {
        console.log('No WebSocket data, using mock data');
        startMockData();
        addEvent('🔧 Используются демо-данные (сервер не обнаружен)');
    }
}, 5000);

addEvent('🌱 Система автоматической теплицы запущена');