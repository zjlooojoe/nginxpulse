<template>
  <header class="page-header">
    <div class="page-title">
      <span class="title-chip">数据日报</span>
      <p class="title-sub">每日概览 · 关键指标 · 行为洞察</p>
    </div>
    <div class="header-actions">
      <WebsiteSelect
        v-model="currentWebsiteId"
        :websites="websites"
        :loading="websitesLoading"
        id="daily-website-selector"
        label="站点"
      />
      <DatePicker
        v-model="currentDate"
        class="daily-date-picker"
        dateFormat="yy-mm-dd"
        updateModelType="string"
        :showClear="false"
        showButtonBar
        :showIcon="true"
      />
      <ThemeToggle />
    </div>
  </header>

  <section class="daily-kpi-grid">
    <div class="card daily-kpi-card">
      <div class="daily-kpi-header">
        <div>
          <div class="daily-kpi-title">浏览量 (PV)</div>
          <div class="daily-kpi-date">{{ currentDate }}</div>
        </div>
        <span class="daily-kpi-icon orange"><i class="ri-pages-line"></i></span>
      </div>
      <div class="daily-kpi-value">{{ kpiMetrics.pv.valueText }}</div>
      <div class="daily-kpi-compare">
        <span>对比前日</span>
        <span class="daily-kpi-delta" :class="kpiMetrics.pv.deltaClass">{{ kpiMetrics.pv.deltaText }}</span>
        <span class="daily-kpi-rate" :class="kpiMetrics.pv.rateClass">{{ kpiMetrics.pv.rateText }}</span>
      </div>
    </div>
    <div class="card daily-kpi-card">
      <div class="daily-kpi-header">
        <div>
          <div class="daily-kpi-title">访客数 (UV)</div>
          <div class="daily-kpi-date">{{ currentDate }}</div>
        </div>
        <span class="daily-kpi-icon green"><i class="ri-user-3-line"></i></span>
      </div>
      <div class="daily-kpi-value">{{ kpiMetrics.uv.valueText }}</div>
      <div class="daily-kpi-compare">
        <span>对比前日</span>
        <span class="daily-kpi-delta" :class="kpiMetrics.uv.deltaClass">{{ kpiMetrics.uv.deltaText }}</span>
        <span class="daily-kpi-rate" :class="kpiMetrics.uv.rateClass">{{ kpiMetrics.uv.rateText }}</span>
      </div>
    </div>
    <div class="card daily-kpi-card">
      <div class="daily-kpi-header">
        <div>
          <div class="daily-kpi-title">会话数</div>
          <div class="daily-kpi-date">{{ currentDate }}</div>
        </div>
        <span class="daily-kpi-icon blue"><i class="ri-chat-3-line"></i></span>
      </div>
      <div class="daily-kpi-value">{{ kpiMetrics.session.valueText }}</div>
      <div class="daily-kpi-compare">
        <span>对比前日</span>
        <span class="daily-kpi-delta" :class="kpiMetrics.session.deltaClass">{{ kpiMetrics.session.deltaText }}</span>
        <span class="daily-kpi-rate" :class="kpiMetrics.session.rateClass">{{ kpiMetrics.session.rateText }}</span>
      </div>
    </div>
    <div class="card daily-kpi-card">
      <div class="daily-kpi-header">
        <div>
          <div class="daily-kpi-title">跳出率</div>
          <div class="daily-kpi-date">{{ currentDate }}</div>
        </div>
        <span class="daily-kpi-icon purple"><i class="ri-run-line"></i></span>
      </div>
      <div class="daily-kpi-value">{{ kpiMetrics.bounce.valueText }}</div>
      <div class="daily-kpi-compare">
        <span>对比前日</span>
        <span class="daily-kpi-delta" :class="kpiMetrics.bounce.deltaClass">{{ kpiMetrics.bounce.deltaText }}</span>
        <span class="daily-kpi-rate" :class="kpiMetrics.bounce.rateClass">{{ kpiMetrics.bounce.rateText }}</span>
      </div>
    </div>
    <div class="card daily-kpi-card">
      <div class="daily-kpi-header">
        <div>
          <div class="daily-kpi-title">平均访问时长</div>
          <div class="daily-kpi-date">{{ currentDate }}</div>
        </div>
        <span class="daily-kpi-icon teal"><i class="ri-time-line"></i></span>
      </div>
      <div class="daily-kpi-value">{{ kpiMetrics.duration.valueText }}</div>
      <div class="daily-kpi-compare">
        <span>对比前日</span>
        <span class="daily-kpi-delta" :class="kpiMetrics.duration.deltaClass">{{ kpiMetrics.duration.deltaText }}</span>
        <span class="daily-kpi-rate" :class="kpiMetrics.duration.rateClass">{{ kpiMetrics.duration.rateText }}</span>
      </div>
    </div>
  </section>

  <section class="daily-mini-grid">
    <div class="card daily-mini-card">
      <div class="daily-mini-title">每个 IP 平均贡献浏览量</div>
      <div class="daily-mini-body">
        <div class="daily-mini-metric">
          <div class="daily-mini-label">变化率</div>
          <div class="daily-mini-value" :class="ipAvg.rateClass">{{ ipAvg.rateText }}</div>
        </div>
        <div class="daily-mini-meta">
          <div>昨日 <span>{{ ipAvg.currentText }}</span></div>
          <div>前日 <span>{{ ipAvg.prevText }}</span></div>
        </div>
      </div>
    </div>
    <div class="card daily-mini-card">
      <div class="daily-mini-title">每个 UV 平均贡献浏览量</div>
      <div class="daily-mini-body">
        <div class="daily-mini-metric">
          <div class="daily-mini-label">变化率</div>
          <div class="daily-mini-value" :class="ipAvg.rateClass">{{ ipAvg.rateText }}</div>
        </div>
        <div class="daily-mini-meta">
          <div>昨日 <span>{{ ipAvg.currentText }}</span></div>
          <div>前日 <span>{{ ipAvg.prevText }}</span></div>
        </div>
      </div>
    </div>
    <div class="card daily-trend-card">
      <div class="daily-trend-header">
        <div class="daily-trend-title">昨日 IP 流量趋势</div>
        <div class="daily-trend-sub">{{ trendSummary }}</div>
      </div>
      <div class="daily-trend-chart">
        <canvas ref="ipChartRef"></canvas>
      </div>
    </div>
  </section>

  <section class="card daily-section">
    <div class="daily-section-header">
      <div class="daily-section-title">
        <span class="section-icon blue"><i class="ri-compass-3-line"></i></span>
        来路分析
      </div>
    </div>
    <div class="daily-section-body daily-source-grid">
      <div class="daily-donut-card">
        <div class="daily-donut">
          <canvas ref="sourceChartRef"></canvas>
        </div>
        <div class="daily-summary-cards">
          <div class="daily-summary-card search">
            <div class="daily-summary-label">搜索引擎</div>
            <div class="daily-summary-value">{{ sourceCards.search.countText }}</div>
            <div class="daily-summary-rate" :class="sourceCards.search.rateClass">{{ sourceCards.search.rateText }}</div>
          </div>
          <div class="daily-summary-card direct">
            <div class="daily-summary-label">直接访问</div>
            <div class="daily-summary-value">{{ sourceCards.direct.countText }}</div>
            <div class="daily-summary-rate" :class="sourceCards.direct.rateClass">{{ sourceCards.direct.rateText }}</div>
          </div>
          <div class="daily-summary-card external">
            <div class="daily-summary-label">外部链接</div>
            <div class="daily-summary-value">{{ sourceCards.external.countText }}</div>
            <div class="daily-summary-rate" :class="sourceCards.external.rateClass">{{ sourceCards.external.rateText }}</div>
          </div>
        </div>
      </div>
      <div class="daily-table-card">
        <div class="daily-tab-bar">
          <button class="daily-tab" :class="{ active: sourceTab === 'referer' }" @click="sourceTab = 'referer'">来路 TOP10</button>
          <button class="daily-tab" :class="{ active: sourceTab === 'search' }" @click="sourceTab = 'search'">搜索引擎</button>
        </div>
        <div class="table-wrapper">
          <table class="ranking-table" v-show="sourceTab === 'referer'">
            <thead>
              <tr>
                <th>来路 TOP10</th>
                <th>IP 数</th>
                <th>IP 变化</th>
                <th>变化率</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="!refererRows.length">
                <td colspan="4">暂无数据</td>
              </tr>
              <tr v-else v-for="row in refererRows" :key="row.label">
                <td>{{ row.label }}</td>
                <td>{{ row.valueText }}</td>
                <td :class="row.deltaClass">{{ row.deltaText }}</td>
                <td :class="row.rateClass">{{ row.rateText }}</td>
              </tr>
            </tbody>
          </table>
          <table class="ranking-table" v-show="sourceTab === 'search'">
            <thead>
              <tr>
                <th>搜索引擎</th>
                <th>IP 数</th>
                <th>IP 变化</th>
                <th>变化率</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="!searchRows.length">
                <td colspan="4">暂无数据</td>
              </tr>
              <tr v-else v-for="row in searchRows" :key="row.label">
                <td>{{ row.label }}</td>
                <td>{{ row.valueText }}</td>
                <td :class="row.deltaClass">{{ row.deltaText }}</td>
                <td :class="row.rateClass">{{ row.rateText }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="daily-summary-strip">{{ sourceSummary }}</div>
      </div>
    </div>
  </section>

  <section class="daily-dual-grid">
    <div class="card daily-section">
      <div class="daily-section-header">
        <div class="daily-section-title">
          <span class="section-icon orange"><i class="ri-pages-line"></i></span>
          内容分析
        </div>
      </div>
      <div class="table-wrapper">
        <table class="ranking-table">
          <thead>
            <tr>
              <th>最受欢迎页面 TOP10</th>
              <th>贡献 IP</th>
              <th>环比变化</th>
              <th>贡献 PV</th>
              <th>环比变化</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!contentRows.length">
              <td colspan="5">暂无数据</td>
            </tr>
            <tr v-else v-for="row in contentRows" :key="row.label">
              <td>{{ row.label }}</td>
              <td>{{ row.uvText }}</td>
              <td :class="row.uvRateClass">{{ row.uvRateText }}</td>
              <td>{{ row.pvText }}</td>
              <td :class="row.pvRateClass">{{ row.pvRateText }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <div class="card daily-section">
      <div class="daily-section-header">
        <div class="daily-section-title">
          <span class="section-icon green"><i class="ri-user-heart-line"></i></span>
          访客分析
        </div>
      </div>
      <div class="daily-visitor-grid">
        <div class="daily-donut">
          <canvas ref="visitorChartRef"></canvas>
        </div>
        <div class="daily-visitor-table">
          <table class="ranking-table">
            <thead>
              <tr>
                <th>访客</th>
                <th>占比</th>
                <th>环比</th>
                <th>平均访问时长</th>
                <th>平均贡献页面浏览量</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="!visitorRows.length">
                <td colspan="5">暂无数据</td>
              </tr>
              <tr v-else v-for="row in visitorRows" :key="row.label">
                <td>{{ row.label }}</td>
                <td>{{ row.shareText }}</td>
                <td :class="row.rateClass">{{ row.rateText }}</td>
                <td>{{ row.durationText }}</td>
                <td>{{ row.avgPvText }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </section>

  <section class="card daily-section">
    <div class="daily-section-header">
      <div class="daily-section-title">
        <span class="section-icon blue"><i class="ri-device-line"></i></span>
        终端分析
      </div>
    </div>
    <div class="daily-device-grid">
      <div class="daily-device-left">
        <div class="daily-donut">
          <canvas ref="deviceChartRef"></canvas>
        </div>
        <div class="daily-device-cards">
          <div class="daily-device-card">
            <div class="daily-device-icon"><i class="ri-computer-line"></i></div>
            <div class="daily-device-label">PC</div>
            <div class="daily-device-value">{{ deviceCards.pc }}</div>
          </div>
          <div class="daily-device-card">
            <div class="daily-device-icon"><i class="ri-apple-line"></i></div>
            <div class="daily-device-label">iOS</div>
            <div class="daily-device-value">{{ deviceCards.ios }}</div>
          </div>
          <div class="daily-device-card">
            <div class="daily-device-icon"><i class="ri-android-line"></i></div>
            <div class="daily-device-label">Android</div>
            <div class="daily-device-value">{{ deviceCards.android }}</div>
          </div>
        </div>
      </div>
      <div class="daily-device-list">
        <div class="daily-list-title">城市 TOP10</div>
        <div class="table-wrapper">
          <table class="ranking-table">
            <thead>
              <tr>
                <th>城市</th>
                <th>占比</th>
                <th>环比</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="!cityRows.length">
                <td colspan="3">暂无数据</td>
              </tr>
              <tr v-else v-for="row in cityRows" :key="row.label">
                <td>{{ row.label }}</td>
                <td>{{ row.shareText }}</td>
                <td :class="row.rateClass">{{ row.rateText }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      <div class="daily-device-list">
        <div class="daily-list-title">浏览器 TOP10</div>
        <div class="table-wrapper">
          <table class="ranking-table">
            <thead>
              <tr>
                <th>浏览器</th>
                <th>占比</th>
                <th>环比</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="!browserRows.length">
                <td colspan="3">暂无数据</td>
              </tr>
              <tr v-else v-for="row in browserRows" :key="row.label">
                <td>{{ row.label }}</td>
                <td>{{ row.shareText }}</td>
                <td :class="row.rateClass">{{ row.rateText }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </section>

  <ParsingOverlay @finished="loadDailyReport" />
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import {
  fetchBrowserStats,
  fetchDeviceStats,
  fetchLocationStats,
  fetchOSStats,
  fetchOverallStats,
  fetchRefererStats,
  fetchSessionSummary,
  fetchTimeSeriesStats,
  fetchUrlStats,
  fetchWebsites,
} from '@/api';
import type { SimpleSeriesStats, TimeSeriesStats, WebsiteInfo } from '@/api/types';
import { formatDate, getUserPreference, saveUserPreference } from '@/utils';
import { Chart } from '@/utils/chartjs';
import ParsingOverlay from '@/components/ParsingOverlay.vue';
import ThemeToggle from '@/components/ThemeToggle.vue';
import WebsiteSelect from '@/components/WebsiteSelect.vue';

const websites = ref<WebsiteInfo[]>([]);
const websitesLoading = ref(true);
const currentWebsiteId = ref('');
const currentDate = ref(formatDate(new Date()));

const overall = ref<Record<string, any> | null>(null);
const sessionSummary = ref<Record<string, any> | null>(null);
const sessionSummaryPrev = ref<Record<string, any> | null>(null);
const timeSeries = ref<TimeSeriesStats | null>(null);
const refererStats = ref<SimpleSeriesStats | null>(null);
const refererPrev = ref<SimpleSeriesStats | null>(null);
const urlStats = ref<SimpleSeriesStats | null>(null);
const urlPrev = ref<SimpleSeriesStats | null>(null);
const deviceStats = ref<SimpleSeriesStats | null>(null);
const devicePrev = ref<SimpleSeriesStats | null>(null);
const osStats = ref<SimpleSeriesStats | null>(null);
const osPrev = ref<SimpleSeriesStats | null>(null);
const browserStats = ref<SimpleSeriesStats | null>(null);
const browserPrev = ref<SimpleSeriesStats | null>(null);
const cityStats = ref<SimpleSeriesStats | null>(null);
const cityPrev = ref<SimpleSeriesStats | null>(null);

const sourceTab = ref<'referer' | 'search'>('referer');

const ipChartRef = ref<HTMLCanvasElement | null>(null);
const sourceChartRef = ref<HTMLCanvasElement | null>(null);
const visitorChartRef = ref<HTMLCanvasElement | null>(null);
const deviceChartRef = ref<HTMLCanvasElement | null>(null);

let ipChart: Chart | null = null;
let sourceChart: Chart | null = null;
let visitorChart: Chart | null = null;
let deviceChart: Chart | null = null;

let dailyRequestId = 0;

const trendSummary = computed(() => {
  if (!timeSeries.value || !timeSeries.value.labels) {
    return '数据就绪后更新峰谷信息';
  }
  const data = timeSeries.value.visitors || [];
  if (!data.length) {
    return '暂无趋势数据';
  }
  const maxIndex = data.indexOf(Math.max(...data));
  const minIndex = data.indexOf(Math.min(...data));
  return `最高峰是 ${maxIndex} 时，最低谷是 ${minIndex} 时`;
});

const kpiMetrics = computed(() => {
  const current = overall.value || {};
  const prev = current.compare?.previous || {};
  const currentSession = sessionSummary.value || {};
  const prevSession = sessionSummaryPrev.value || {};

  return {
    pv: buildMetric(current.pv || 0, prev.pv || 0),
    uv: buildMetric(current.uv || 0, prev.uv || 0),
    session: buildMetric(current.sessionCount || 0, prev.sessionCount || 0),
    bounce: buildPercentMetric(currentSession.bounceRate || 0, prevSession.bounceRate || 0),
    duration: buildDurationMetric(currentSession.avgDurationSeconds || 0, prevSession.avgDurationSeconds || 0),
  }
});

const ipAvg = computed(() => {
  const current = overall.value || {};
  const prev = current.compare?.previous || {};
  const avg = current.uv > 0 ? current.pv / current.uv : 0;
  const prevAvg = prev.uv > 0 ? prev.pv / prev.uv : 0;
  const rate = calcRate(avg, prevAvg);
  return {
    currentText: formatNumber(avg),
    prevText: formatNumber(prevAvg),
    rateText: formatRate(rate),
    rateClass: rateClass(rate),
  }
});

const sourceGroups = computed(() => groupReferers(refererStats.value));
const sourcePrevGroups = computed(() => groupReferers(refererPrev.value));

const sourceCards = computed(() => {
  const current = sourceGroups.value;
  const prev = sourcePrevGroups.value;
  return {
    search: buildSourceCard(current.search, prev.search),
    direct: buildSourceCard(current.direct, prev.direct),
    external: buildSourceCard(current.external, prev.external),
  }
});

const refererRows = computed(() => buildRefererRows(refererStats.value, refererPrev.value, false));
const searchRows = computed(() => buildRefererRows(refererStats.value, refererPrev.value, true));
const sourceSummary = computed(() => buildSourceSummary(refererStats.value, refererPrev.value));
const contentRows = computed(() => buildContentRows(urlStats.value, urlPrev.value));
const visitorRows = computed(() => buildVisitorRows(overall.value, sessionSummary.value));
const deviceCards = computed(() => buildDeviceCards(deviceStats.value, osStats.value));
const browserRows = computed(() => buildSimpleRows(browserStats.value, browserPrev.value));
const cityRows = computed(() => buildSimpleRows(cityStats.value, cityPrev.value));

onMounted(() => {
  initDateFromQuery();
  loadWebsites();
});

onBeforeUnmount(() => {
  destroyCharts();
});

watch(currentWebsiteId, (value) => {
  if (value) {
    saveUserPreference('selectedWebsite', value);
  }
  loadDailyReport();
});

watch(currentDate, (value) => {
  if (value) {
    saveUserPreference('dailyReportDate', value);
  }
  loadDailyReport();
});

watch(timeSeries, (stats) => {
  if (stats) {
    renderTrend(stats);
  }
});

watch(sourceGroups, (groups) => {
  renderSourceDonut(groups);
});

watch(visitorRows, () => {
  renderVisitorDonut();
});

watch(deviceCards, () => {
  renderDeviceDonut();
});

function initDateFromQuery() {
  const queryDate = getDateFromQuery();
  const savedDate = getUserPreference('dailyReportDate', '');
  const defaultDate = queryDate || savedDate || formatDate(new Date());
  currentDate.value = defaultDate;
  if (queryDate) {
    saveUserPreference('dailyReportDate', defaultDate);
  }
}

async function loadWebsites() {
  websitesLoading.value = true;
  try {
    const data = await fetchWebsites();
    websites.value = data || [];
    const saved = getUserPreference('selectedWebsite', '');
    if (saved && websites.value.find((site) => site.id === saved)) {
      currentWebsiteId.value = saved;
    } else if (websites.value.length > 0) {
      currentWebsiteId.value = websites.value[0].id;
    } else {
      currentWebsiteId.value = '';
    }
  } catch (error) {
    console.error('初始化网站失败:', error);
    websites.value = [];
    currentWebsiteId.value = '';
  } finally {
    websitesLoading.value = false;
  }
}

async function loadDailyReport() {
  if (!currentWebsiteId.value || !currentDate.value) {
    return;
  }

  const requestId = ++dailyRequestId;
  const dateStr = currentDate.value;
  const prevDate = shiftDate(dateStr, -1);

  try {
    const [
      overallData,
      sessionData,
      sessionPrevData,
      timeSeriesData,
      refererData,
      refererPrevData,
      urlData,
      urlPrevData,
      deviceData,
      devicePrevData,
      osData,
      osPrevData,
      browserData,
      browserPrevData,
      cityData,
      cityPrevData,
    ] = await Promise.all([
      fetchOverallStats(currentWebsiteId.value, dateStr),
      fetchSessionSummary(currentWebsiteId.value, dateStr),
      fetchSessionSummary(currentWebsiteId.value, prevDate),
      fetchTimeSeriesStats(currentWebsiteId.value, prevDate, 'hourly'),
      fetchRefererStats(currentWebsiteId.value, dateStr, 10),
      fetchRefererStats(currentWebsiteId.value, prevDate, 10),
      fetchUrlStats(currentWebsiteId.value, dateStr, 10),
      fetchUrlStats(currentWebsiteId.value, prevDate, 10),
      fetchDeviceStats(currentWebsiteId.value, dateStr, 10),
      fetchDeviceStats(currentWebsiteId.value, prevDate, 10),
      fetchOSStats(currentWebsiteId.value, dateStr, 10),
      fetchOSStats(currentWebsiteId.value, prevDate, 10),
      fetchBrowserStats(currentWebsiteId.value, dateStr, 10),
      fetchBrowserStats(currentWebsiteId.value, prevDate, 10),
      fetchLocationStats(currentWebsiteId.value, dateStr, 'city', 10),
      fetchLocationStats(currentWebsiteId.value, prevDate, 'city', 10),
    ]);

    if (requestId !== dailyRequestId) {
      return;
    }

    overall.value = overallData;
    sessionSummary.value = sessionData;
    sessionSummaryPrev.value = sessionPrevData;
    timeSeries.value = timeSeriesData;
    refererStats.value = refererData;
    refererPrev.value = refererPrevData;
    urlStats.value = urlData;
    urlPrev.value = urlPrevData;
    deviceStats.value = deviceData;
    devicePrev.value = devicePrevData;
    osStats.value = osData;
    osPrev.value = osPrevData;
    browserStats.value = browserData;
    browserPrev.value = browserPrevData;
    cityStats.value = cityData;
    cityPrev.value = cityPrevData;
  } catch (error) {
    console.error('加载日报失败:', error);
  }
}

function renderTrend(stats: TimeSeriesStats) {
  if (!ipChartRef.value) {
    return;
  }
  if (ipChart) {
    ipChart.destroy();
  }
  const ctx = ipChartRef.value.getContext('2d');
  if (!ctx) {
    return;
  }
  const gradient = ctx.createLinearGradient(0, 0, 0, ipChartRef.value.height || 160);
  gradient.addColorStop(0, 'rgba(30, 123, 255, 0.35)');
  gradient.addColorStop(1, 'rgba(30, 123, 255, 0.05)');
  ipChart = new Chart(ctx, {
    type: 'line',
    data: {
      labels: stats.labels,
      datasets: [
        {
          label: 'IP 流量',
          data: stats.visitors,
          borderColor: '#1e7bff',
          backgroundColor: gradient,
          borderWidth: 2,
          tension: 0.4,
          pointRadius: 0,
          fill: true,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      scales: {
        x: { ticks: { color: '#94a3b8', maxTicksLimit: 12 }, grid: { display: false } },
        y: { ticks: { color: '#94a3b8' }, grid: { color: 'rgba(148, 163, 184, 0.2)' } },
      },
      plugins: { legend: { display: false }, tooltip: { mode: 'index', intersect: false } },
    },
  });
}

function renderSourceDonut(groups: Record<string, number>) {
  if (!sourceChartRef.value) {
    return;
  }
  if (sourceChart) {
    sourceChart.destroy();
  }
  const ctx = sourceChartRef.value.getContext('2d');
  if (!ctx) {
    return;
  }
  sourceChart = new Chart(ctx, {
    type: 'doughnut',
    data: {
      labels: ['搜索引擎', '直接访问', '外部链接'],
      datasets: [
        {
          data: [groups.search, groups.direct, groups.external],
          backgroundColor: ['#1e7bff', '#ff8a3d', '#22c55e'],
          borderWidth: 0,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      cutout: '60%',
      plugins: { legend: { display: false } },
    },
  });
}

function renderVisitorDonut() {
  if (!visitorChartRef.value) {
    return;
  }
  if (visitorChart) {
    visitorChart.destroy();
  }
  const ctx = visitorChartRef.value.getContext('2d');
  if (!ctx) {
    return;
  }
  const current = overall.value || {};
  visitorChart = new Chart(ctx, {
    type: 'doughnut',
    data: {
      labels: ['新访客', '老访客'],
      datasets: [
        {
          data: [current.newVisitorCount || 0, current.returningVisitorCount || 0],
          backgroundColor: ['#1e7bff', '#22c55e'],
          borderWidth: 0,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      cutout: '65%',
      plugins: { legend: { position: 'bottom' } },
    },
  });
}

function renderDeviceDonut() {
  if (!deviceChartRef.value) {
    return;
  }
  if (deviceChart) {
    deviceChart.destroy();
  }
  const ctx = deviceChartRef.value.getContext('2d');
  if (!ctx) {
    return;
  }
  const pcCount = getDeviceCount(deviceStats.value, ['桌面设备']);
  const iosCount = getOsCount(osStats.value, ['ios', 'iphone', 'ipad']);
  const androidCount = getOsCount(osStats.value, ['android']);
  deviceChart = new Chart(ctx, {
    type: 'doughnut',
    data: {
      labels: ['PC', 'iOS', 'Android'],
      datasets: [
        {
          data: [pcCount, iosCount, androidCount],
          backgroundColor: ['#1e7bff', '#22c55e', '#ff8a3d', '#6366f1'],
          borderWidth: 0,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      cutout: '60%',
      plugins: { legend: { position: 'bottom' } },
    },
  });
}

function destroyCharts() {
  if (ipChart) {
    ipChart.destroy();
    ipChart = null;
  }
  if (sourceChart) {
    sourceChart.destroy();
    sourceChart = null;
  }
  if (visitorChart) {
    visitorChart.destroy();
    visitorChart = null;
  }
  if (deviceChart) {
    deviceChart.destroy();
    deviceChart = null;
  }
}

function buildMetric(current: number, prev: number) {
  const delta = current - prev;
  const rate = calcRate(current, prev);
  return {
    valueText: formatNumber(current),
    deltaText: formatSigned(delta),
    rateText: formatRate(rate),
    deltaClass: deltaClass(delta),
    rateClass: rateClass(rate),
  }
}

function buildPercentMetric(current: number, prev: number) {
  const delta = current - prev;
  const rate = calcRate(current, prev);
  return {
    valueText: formatPercent(current),
    deltaText: formatSignedPercent(delta),
    rateText: formatRate(rate),
    deltaClass: deltaClass(delta),
    rateClass: rateClass(rate),
  }
}

function buildDurationMetric(current: number, prev: number) {
  const delta = current - prev;
  const rate = calcRate(current, prev);
  return {
    valueText: formatDuration(current),
    deltaText: formatSignedDuration(delta),
    rateText: formatRate(rate),
    deltaClass: deltaClass(delta),
    rateClass: rateClass(rate),
  }
}

function buildSourceCard(current: number, prev: number) {
  const rate = calcRate(current, prev);
  return {
    countText: formatNumber(current),
    rateText: formatRate(rate),
    rateClass: rateClass(rate),
  }
}

function buildRefererRows(stats: SimpleSeriesStats | null, prevStats: SimpleSeriesStats | null, searchOnly: boolean) {
  const prevMap = buildStatMap(prevStats);
  const keys = stats?.key || [];
  const uvs = stats?.uv || [];
  const rows = keys
    .map((key, idx) => ({ label: key, value: uvs[idx] || 0 }))
    .filter((item) => (searchOnly ? isSearchEngine(item.label) : true));

  return rows.map((item) => {
    const prev = prevMap[item.label] || 0;
    const delta = item.value - prev;
    const rate = calcRate(item.value, prev);
    return {
      label: item.label,
      valueText: formatNumber(item.value),
      deltaText: formatSigned(delta),
      deltaClass: deltaClass(delta),
      rateText: formatRate(rate),
      rateClass: rateClass(rate),
    }
  });
}

function buildContentRows(stats: SimpleSeriesStats | null, prevStats: SimpleSeriesStats | null) {
  const prevPV = buildStatMap(prevStats, 'pv');
  const prevUV = buildStatMap(prevStats, 'uv');

  const keys = stats?.key || [];
  const pvs = stats?.pv || [];
  const uvs = stats?.uv || [];

  return keys.map((key, idx) => {
    const pv = pvs[idx] || 0;
    const uv = uvs[idx] || 0;
    const pvRate = calcRate(pv, prevPV[key] || 0);
    const uvRate = calcRate(uv, prevUV[key] || 0);
    return {
      label: key,
      pvText: formatNumber(pv),
      uvText: formatNumber(uv),
      pvRateText: formatRate(pvRate),
      uvRateText: formatRate(uvRate),
      pvRateClass: rateClass(pvRate),
      uvRateClass: rateClass(uvRate),
    }
  });
}

function buildVisitorRows(overallData: Record<string, any> | null, sessionData: Record<string, any> | null) {
  if (!overallData) {
    return [];
  }
  const newCount = overallData.newVisitorCount || 0;
  const returningCount = overallData.returningVisitorCount || 0;
  const total = newCount + returningCount;

  const prevNew = overallData.prevNewVisitorCount || 0;
  const prevReturning = overallData.prevReturningVisitorCount || 0;
  const avgDuration = sessionData?.avgDurationSeconds || 0;
  const avgPV = overallData.uv ? overallData.pv / overallData.uv : 0;

  const rows = [
    { label: '新访客', count: newCount, prev: prevNew },
    { label: '老访客', count: returningCount, prev: prevReturning },
  ];

  return rows.map((item) => {
    const share = total > 0 ? item.count / total : 0;
    const rate = calcRate(item.count, item.prev);
    return {
      label: item.label,
      shareText: formatPercent(share),
      rateText: formatRate(rate),
      rateClass: rateClass(rate),
      durationText: formatDuration(avgDuration),
      avgPvText: formatNumber(avgPV),
    }
  });
}

function buildDeviceCards(deviceData: SimpleSeriesStats | null, osData: SimpleSeriesStats | null) {
  const pcCount = getDeviceCount(deviceData, ['桌面设备']);
  const iosCount = getOsCount(osData, ['ios', 'iphone', 'ipad']);
  const androidCount = getOsCount(osData, ['android']);
  return {
    pc: formatNumber(pcCount),
    ios: formatNumber(iosCount),
    android: formatNumber(androidCount),
  }
}

function buildSimpleRows(stats: SimpleSeriesStats | null, prevStats: SimpleSeriesStats | null) {
  const prevMap = buildStatMap(prevStats);
  const keys = stats?.key || [];
  const uvs = stats?.uv || [];
  if (!keys.length) {
    return [];
  }

  const total = uvs.reduce((sum, val) => sum + val, 0);
  return keys.map((key, idx) => {
    const value = uvs[idx] || 0;
    const share = total > 0 ? value / total : 0;
    const prev = prevMap[key] || 0;
    const rate = calcRate(value, prev);
    return {
      label: key,
      shareText: formatPercent(share),
      rateText: formatRate(rate),
      rateClass: rateClass(rate),
    }
  });
}

function buildSourceSummary(stats: SimpleSeriesStats | null, prevStats: SimpleSeriesStats | null) {
  if (!stats || !stats.key) {
    return '暂无来路数据';
  }
  const prevMap = buildStatMap(prevStats);
  const keys = stats.key || [];
  const uvs = stats.uv || [];
  if (!keys.length) {
    return '暂无来路数据';
  }
  const diffs = keys.map((key, idx) => ({
    key,
    rate: calcRate(uvs[idx] || 0, prevMap[key] || 0),
  }));
  const rising = diffs.reduce((acc, item) => (item.rate !== null && item.rate > (acc.rate ?? -Infinity) ? item : acc), {
    rate: -Infinity,
  });
  const falling = diffs.reduce((acc, item) => (item.rate !== null && item.rate < (acc.rate ?? Infinity) ? item : acc), {
    rate: Infinity,
  });
  return `上升最快的是 ${rising.key || '-'} ${formatRate(rising.rate)}，下降最快的是 ${falling.key || '-'} ${formatRate(
    falling.rate
  )}`;
}

function calcRate(current: number, prev: number) {
  if (prev === 0) {
    if (current === 0) {
      return 0;
    }
    return null;
  }
  return (current - prev) / prev;
}

function formatNumber(value: number) {
  if (value === null || value === undefined) {
    return '--';
  }
  return Number(value).toLocaleString('zh-CN', { maximumFractionDigits: 2 });
}

function formatPercent(value: number) {
  if (value === null || value === undefined) {
    return '--';
  }
  return `${(value * 100).toFixed(2)}%`;
}

function formatRate(rate: number | null) {
  if (rate === null) {
    return '-';
  }
  return `${rate >= 0 ? '+' : ''}${(rate * 100).toFixed(2)}%`;
}

function formatSigned(value: number) {
  if (value === null || value === undefined) {
    return '--';
  }
  return `${value >= 0 ? '+' : ''}${value}`;
}

function formatSignedPercent(value: number) {
  if (value === null || value === undefined) {
    return '--';
  }
  return `${value >= 0 ? '+' : ''}${(value * 100).toFixed(2)}%`;
}

function formatSignedDuration(value: number) {
  if (value === null || value === undefined) {
    return '--';
  }
  const prefix = value >= 0 ? '+' : '-';
  return `${prefix}${formatDuration(Math.abs(value))}`;
}

function formatDuration(seconds: number) {
  const total = Math.max(0, Math.floor(seconds));
  const hours = Math.floor(total / 3600);
  const minutes = Math.floor((total % 3600) / 60);
  const secs = total % 60;
  return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
}

function deltaClass(delta: number) {
  if (delta > 0) return 'trend-up';
  if (delta < 0) return 'trend-down';
  return 'trend-flat';
}

function rateClass(rate: number | null) {
  if (rate === null) return 'trend-flat';
  if (rate > 0) return 'trend-up';
  if (rate < 0) return 'trend-down';
  return 'trend-flat';
}

function buildStatMap(stats: SimpleSeriesStats | null, field: 'uv' | 'pv' = 'uv') {
  const map: Record<string, number> = {};
  if (!stats || !stats.key) {
    return map;
  }
  const values = stats[field] || [];
  stats.key.forEach((key, idx) => {
    map[key] = values[idx] || 0;
  });
  return map;
}

function getDeviceCount(stats: SimpleSeriesStats | null, names: string[]) {
  const map = buildStatMap(stats);
  return names.reduce((sum, name) => sum + (map[name] || 0), 0);
}

function getOsCount(stats: SimpleSeriesStats | null, keywords: string[]) {
  if (!stats || !stats.key) {
    return 0;
  }
  let total = 0;
  const values = stats.uv || [];
  stats.key.forEach((key, idx) => {
    const lower = String(key || '').toLowerCase();
    if (keywords.some((word) => lower.includes(word))) {
      total += values[idx] || 0;
    }
  });
  return total;
}

function groupReferers(stats: SimpleSeriesStats | null) {
  const keys = stats?.key || [];
  const uvs = stats?.uv || [];
  const groups = { search: 0, direct: 0, external: 0 };
  keys.forEach((key, idx) => {
    const value = uvs[idx] || 0;
    if (isDirectReferer(key)) {
      groups.direct += value;
    } else if (isSearchEngine(key)) {
      groups.search += value;
    } else {
      groups.external += value;
    }
  });
  return groups;
}

function isDirectReferer(value: string) { return value === '直接输入网址访问' || value === '站内访问' || value === '-' || value === ''; }

const searchEngines = ['baidu.', 'google.', 'bing.', 'sogou.', '360.', 'so.com', 'yahoo.', 'duckduckgo.', 'yandex.'];

function isSearchEngine(value: string) {
  const lower = value.toLowerCase();
  return searchEngines.some((engine) => lower.includes(engine));
}

function shiftDate(dateStr: string, offsetDays: number) {
  const date = new Date(dateStr);
  if (Number.isNaN(date.getTime())) {
    return formatDate(new Date());
  }
  date.setDate(date.getDate() + offsetDays);
  return formatDate(date);
}

function getDateFromQuery() {
  const params = new URLSearchParams(window.location.search || '');
  const raw = params.get('date');
  if (!raw) {
    return '';
  }
  if (raw === 'today') {
    return formatDate(new Date());
  }
  if (/^\d{4}-\d{2}-\d{2}$/.test(raw)) {
    return raw;
  }
  return '';
}
</script>

<style scoped lang="scss">
:global(.daily-page) {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.daily-date-picker {
  min-width: 190px;
}

.daily-kpi-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.daily-kpi-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-height: 150px;
}

.daily-kpi-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.daily-kpi-title {
  font-weight: 600;
  font-size: 14px;
}

.daily-kpi-date {
  color: var(--muted);
  font-size: 12px;
  margin-top: 2px;
}

.daily-kpi-icon {
  width: 36px;
  height: 36px;
  border-radius: 12px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  background: var(--panel-muted);
  color: var(--primary);
}

.daily-kpi-icon.orange {
  color: var(--accent);
  background: rgba(245, 158, 11, 0.12);
}

.daily-kpi-icon.green {
  color: var(--green);
  background: rgba(34, 197, 94, 0.12);
}

.daily-kpi-icon.blue {
  color: var(--primary);
  background: rgba(var(--primary-color-rgb), 0.12);
}

.daily-kpi-icon.purple {
  color: #0ea5e9;
  background: rgba(14, 165, 233, 0.12);
}

.daily-kpi-icon.teal {
  color: #14b8a6;
  background: rgba(20, 184, 166, 0.12);
}

.daily-kpi-value {
  font-size: 28px;
  font-weight: 700;
}

.daily-kpi-compare {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: var(--muted);
}

.daily-kpi-delta,
.daily-kpi-rate {
  font-weight: 600;
}

.trend-up {
  color: #ef4444;
}

.trend-down {
  color: #16a34a;
}

.trend-flat {
  color: var(--muted);
}

.daily-mini-grid {
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  gap: 16px;
}

.daily-mini-card {
  grid-column: span 3;
  background: linear-gradient(135deg, rgba(var(--primary-color-rgb), 0.12), rgba(var(--primary-color-rgb), 0.03));
  min-height: 140px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.daily-mini-card:nth-of-type(2) {
  background: linear-gradient(135deg, rgba(14, 165, 233, 0.12), rgba(14, 165, 233, 0.03));
}

.daily-mini-title {
  font-weight: 600;
  font-size: 14px;
}

.daily-mini-body {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
}

.daily-mini-metric {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.daily-mini-label {
  font-size: 12px;
  color: var(--muted);
}

.daily-mini-value {
  font-size: 24px;
  font-weight: 700;
}

.daily-mini-meta {
  font-size: 12px;
  color: var(--muted);
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.daily-trend-card {
  grid-column: span 6;
  min-height: 240px;
}

.daily-trend-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.daily-trend-title {
  font-weight: 600;
}

.daily-trend-sub {
  font-size: 12px;
  color: var(--muted);
}

.daily-trend-chart {
  height: 180px;
}

.daily-source-grid {
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  gap: 16px;
}

.daily-donut-card {
  grid-column: span 4;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.daily-donut {
  height: 220px;
}

.daily-summary-cards {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}

.daily-summary-card {
  border-radius: 14px;
  padding: 14px;
  border: 1px solid var(--border);
  background: var(--panel-muted);
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.daily-summary-card.search {
  border-color: rgba(30, 123, 255, 0.25);
}

.daily-summary-card.direct {
  border-color: rgba(245, 158, 11, 0.2);
}

.daily-summary-card.external {
  border-color: rgba(34, 197, 94, 0.2);
}

.daily-summary-label {
  font-size: 12px;
  color: var(--muted);
}

.daily-summary-value {
  font-size: 18px;
  font-weight: 700;
}

.daily-summary-rate {
  font-size: 12px;
  font-weight: 600;
}

.daily-table-card {
  grid-column: span 8;
}

.daily-tab-bar {
  display: inline-flex;
  gap: 8px;
  padding: 4px;
  background: var(--panel-muted);
  border: 1px solid var(--border);
  border-radius: 999px;
  margin-bottom: 12px;
}

.daily-tab {
  border: none;
  background: transparent;
  padding: 6px 14px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
  color: var(--muted);
  cursor: pointer;
}

.daily-tab.active {
  background: linear-gradient(135deg, var(--primary), var(--primary-strong));
  color: white;
}

.daily-summary-strip {
  margin-top: 10px;
  font-size: 12px;
  color: var(--muted);
}

.daily-dual-grid {
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  gap: 16px;
}

.daily-section {
  grid-column: span 6;
}

.daily-section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.daily-section-title {
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 8px;
}

.daily-visitor-grid {
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  gap: 16px;
}

.daily-visitor-grid .daily-donut {
  grid-column: span 4;
}

.daily-visitor-table {
  grid-column: span 8;
}

.daily-device-grid {
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  gap: 16px;
}

.daily-device-left {
  grid-column: span 4;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.daily-device-cards {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}

.daily-device-card {
  border-radius: 12px;
  border: 1px solid var(--border);
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  background: var(--panel-muted);
}

.daily-device-icon {
  font-size: 18px;
  color: var(--primary);
}

.daily-device-label {
  font-size: 12px;
  color: var(--muted);
}

.daily-device-value {
  font-size: 18px;
  font-weight: 700;
}

.daily-device-list {
  grid-column: span 4;
}

.daily-list-title {
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 8px;
}

@media (max-width: 1200px) {
  .daily-mini-card {
    grid-column: span 6;
  }

  .daily-trend-card {
    grid-column: span 12;
  }

  .daily-donut-card {
    grid-column: span 12;
  }

  .daily-table-card {
    grid-column: span 12;
  }

  .daily-section {
    grid-column: span 12;
  }

  .daily-visitor-grid .daily-donut,
  .daily-visitor-table,
  .daily-device-left,
  .daily-device-list {
    grid-column: span 12;
  }

  .daily-device-cards {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .daily-mini-card {
    grid-column: span 12;
  }

  .daily-device-cards {
    grid-template-columns: 1fr;
  }

  .daily-summary-cards {
    grid-template-columns: 1fr;
  }
}
</style>
