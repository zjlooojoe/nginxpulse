<template>
  <header class="page-header">
    <div class="page-title">
      <span class="title-chip">实时访问</span>
      <p class="title-sub">分钟级监控 · 行为分布 · 快速识别异常</p>
    </div>
    <div class="header-actions">
      <WebsiteSelect
        v-model="currentWebsiteId"
        :websites="websites"
        :loading="websitesLoading"
        id="realtime-website-selector"
        label="站点"
      />
      <div class="realtime-range">
        <button
          v-for="option in windowOptions"
          :key="option"
          class="realtime-range-btn"
          :class="{ active: currentWindow === option }"
          @click="setWindow(option)"
        >
          {{ option }}分钟
        </button>
      </div>
      <ThemeToggle />
    </div>
  </header>

  <section class="realtime-grid">
    <div class="card realtime-card realtime-overview">
      <div class="realtime-card-header">
        <div class="realtime-card-title">
          <span class="section-icon blue"><i class="ri-radar-line"></i></span>
          <span>{{ activeTitle }}</span>
        </div>
      </div>
      <div class="realtime-metric">
        <div class="realtime-value">{{ formatCount(activeCount) }}</div>
        <div class="realtime-mini-bars">
          <span v-for="(bar, index) in activeBars" :key="index" :class="{ active: bar.active }" :style="{ height: `${bar.height}px` }"></span>
        </div>
      </div>
      <div class="realtime-subtitle">{{ deviceSubtitle }}</div>
      <div class="realtime-device-cards">
        <div class="realtime-device-card">
          <div class="realtime-device-icon"><i class="ri-computer-line"></i></div>
          <div class="realtime-device-label">PC端</div>
          <div class="realtime-device-count">{{ formatCount(deviceStats.pc.count) }}</div>
          <div class="realtime-device-rate">{{ formatPercent(deviceStats.pc.percent) }}</div>
        </div>
        <div class="realtime-device-card">
          <div class="realtime-device-icon"><i class="ri-smartphone-line"></i></div>
          <div class="realtime-device-label">移动端</div>
          <div class="realtime-device-count">{{ formatCount(deviceStats.mobile.count) }}</div>
          <div class="realtime-device-rate">{{ formatPercent(deviceStats.mobile.percent) }}</div>
        </div>
        <div class="realtime-device-card">
          <div class="realtime-device-icon"><i class="ri-shield-line"></i></div>
          <div class="realtime-device-label">其他</div>
          <div class="realtime-device-count">{{ formatCount(deviceStats.other.count) }}</div>
          <div class="realtime-device-rate">{{ formatPercent(deviceStats.other.percent) }}</div>
        </div>
      </div>
    </div>

    <div class="card realtime-card">
      <div class="realtime-card-header">
        <div class="realtime-card-title">
          <span class="section-icon blue"><i class="ri-compass-3-line"></i></span>
          来路
        </div>
      </div>
      <div class="realtime-top">
        <span class="realtime-rank">NO.1</span>
        <div class="realtime-top-title">{{ topReferer.name }}</div>
        <div class="realtime-top-meta">
          <span class="realtime-top-count">{{ formatCount(topReferer.count) }}</span>
          <span class="realtime-top-rate">{{ formatPercent(topReferer.percent) }}</span>
        </div>
      </div>
      <div class="realtime-mini-bars">
        <span v-for="(bar, index) in activeBars" :key="index" :class="{ active: bar.active }" :style="{ height: `${bar.height}px` }"></span>
      </div>
      <div class="table-wrapper">
        <table class="ranking-table realtime-table">
          <thead>
            <tr>
              <th>来路</th>
              <th class="realtime-count-col">访客数</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!refererItems.length">
              <td colspan="2">暂无数据</td>
            </tr>
            <tr v-else v-for="item in refererItems" :key="item.name">
              <td>{{ item.name }}</td>
              <td class="realtime-count-col">{{ formatCount(item.count) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="card realtime-card">
      <div class="realtime-card-header">
        <div class="realtime-card-title">
          <span class="section-icon orange"><i class="ri-pages-line"></i></span>
          受访页面
        </div>
      </div>
      <div class="realtime-top">
        <span class="realtime-rank">NO.1</span>
        <div class="realtime-top-title">{{ topPage.name }}</div>
        <div class="realtime-top-meta">
          <span class="realtime-top-count">{{ formatCount(topPage.count) }}</span>
          <span class="realtime-top-rate">{{ formatPercent(topPage.percent) }}</span>
        </div>
      </div>
      <div class="realtime-mini-bars">
        <span v-for="(bar, index) in activeBars" :key="index" :class="{ active: bar.active }" :style="{ height: `${bar.height}px` }"></span>
      </div>
      <div class="table-wrapper">
        <table class="ranking-table realtime-table">
          <thead>
            <tr>
              <th>受访页面标题</th>
              <th class="realtime-count-col">查看次数</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!pageItems.length">
              <td colspan="2">暂无数据</td>
            </tr>
            <tr v-else v-for="item in pageItems" :key="item.name">
              <td>{{ item.name }}</td>
              <td class="realtime-count-col">{{ formatCount(item.count) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="card realtime-card">
      <div class="realtime-card-header">
        <div class="realtime-card-title">
          <span class="section-icon orange"><i class="ri-login-circle-line"></i></span>
          入口页面
        </div>
      </div>
      <div class="realtime-top">
        <span class="realtime-rank">NO.1</span>
        <div class="realtime-top-title">{{ topEntry.name }}</div>
        <div class="realtime-top-meta">
          <span class="realtime-top-count">{{ formatCount(topEntry.count) }}</span>
          <span class="realtime-top-rate">{{ formatPercent(topEntry.percent) }}</span>
        </div>
      </div>
      <div class="realtime-mini-bars">
        <span v-for="(bar, index) in activeBars" :key="index" :class="{ active: bar.active }" :style="{ height: `${bar.height}px` }"></span>
      </div>
      <div class="table-wrapper">
        <table class="ranking-table realtime-table">
          <thead>
            <tr>
              <th>入口页面</th>
              <th class="realtime-count-col">入口次数</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!entryItems.length">
              <td colspan="2">暂无数据</td>
            </tr>
            <tr v-else v-for="item in entryItems" :key="item.name">
              <td>{{ item.name }}</td>
              <td class="realtime-count-col">{{ formatCount(item.count) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="card realtime-card">
      <div class="realtime-card-header">
        <div class="realtime-card-title">
          <span class="section-icon green"><i class="ri-global-line"></i></span>
          浏览器
        </div>
      </div>
      <div class="realtime-top">
        <span class="realtime-rank">NO.1</span>
        <div class="realtime-top-title">{{ topBrowser.name }}</div>
        <div class="realtime-top-meta">
          <span class="realtime-top-count">{{ formatCount(topBrowser.count) }}</span>
          <span class="realtime-top-rate">{{ formatPercent(topBrowser.percent) }}</span>
        </div>
      </div>
      <div class="realtime-mini-bars">
        <span v-for="(bar, index) in activeBars" :key="index" :class="{ active: bar.active }" :style="{ height: `${bar.height}px` }"></span>
      </div>
      <div class="table-wrapper">
        <table class="ranking-table realtime-table">
          <thead>
            <tr>
              <th>浏览器</th>
              <th class="realtime-count-col">访客数</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!browserItems.length">
              <td colspan="2">暂无数据</td>
            </tr>
            <tr v-else v-for="item in browserItems" :key="item.name">
              <td>{{ item.name }}</td>
              <td class="realtime-count-col">{{ formatCount(item.count) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="card realtime-card">
      <div class="realtime-card-header">
        <div class="realtime-card-title">
          <span class="section-icon blue"><i class="ri-map-pin-2-line"></i></span>
          城市/地区
        </div>
      </div>
      <div class="realtime-top">
        <span class="realtime-rank">NO.1</span>
        <div class="realtime-top-title">{{ topCity.name }}</div>
        <div class="realtime-top-meta">
          <span class="realtime-top-count">{{ formatCount(topCity.count) }}</span>
          <span class="realtime-top-rate">{{ formatPercent(topCity.percent) }}</span>
        </div>
      </div>
      <div class="realtime-mini-bars">
        <span v-for="(bar, index) in activeBars" :key="index" :class="{ active: bar.active }" :style="{ height: `${bar.height}px` }"></span>
      </div>
      <div class="table-wrapper">
        <table class="ranking-table realtime-table">
          <thead>
            <tr>
              <th>城市/地区</th>
              <th class="realtime-count-col">访客数</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!cityItems.length">
              <td colspan="2">暂无数据</td>
            </tr>
            <tr v-else v-for="item in cityItems" :key="item.name">
              <td>{{ item.name }}</td>
              <td class="realtime-count-col">{{ formatCount(item.count) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </section>

  <ParsingOverlay @finished="loadRealtime" />
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { fetchRealtimeStats, fetchWebsites } from '@/api';
import type { RealtimeSeriesItem, RealtimeStats, WebsiteInfo } from '@/api/types';
import { getUserPreference, saveUserPreference } from '@/utils';
import ParsingOverlay from '@/components/ParsingOverlay.vue';
import ThemeToggle from '@/components/ThemeToggle.vue';
import WebsiteSelect from '@/components/WebsiteSelect.vue';

const websites = ref<WebsiteInfo[]>([]);
const websitesLoading = ref(true);
const currentWebsiteId = ref('');
const windowOptions = [5, 15, 30];
const currentWindow = ref(30);

const activeCount = ref(0);
const activeSeries = ref<number[]>([]);
const deviceBreakdown = ref<RealtimeSeriesItem[]>([]);
const referers = ref<RealtimeSeriesItem[]>([]);
const pages = ref<RealtimeSeriesItem[]>([]);
const entryPages = ref<RealtimeSeriesItem[]>([]);
const browsers = ref<RealtimeSeriesItem[]>([]);
const locations = ref<RealtimeSeriesItem[]>([]);

let refreshTimer: number | null = null;

const activeTitle = computed(() => `最近${currentWindow.value}分钟活跃访客数`);
const deviceSubtitle = computed(() => `最近${currentWindow.value}分钟的终端设备`);

const activeBars = computed(() => {
  const values = Array.isArray(activeSeries.value) && activeSeries.value.length
    ? activeSeries.value
    : new Array(currentWindow.value).fill(0);
  const maxVal = Math.max(1, ...values);
  return values.map((value) => {
    const ratio = value / maxVal;
    return {
      height: Math.max(6, Math.round(ratio * 24)),
      active: value > 0,
    }
  });
});

const deviceStats = computed(() => {
  const breakdown = deviceBreakdown.value || [];
  const pc = breakdown.find((item) => item.name === 'PC端') || { count: 0, percent: 0 };
  const mobile = breakdown.find((item) => item.name === '移动端') || { count: 0, percent: 0 };
  const other = breakdown.find((item) => item.name === '其他') || { count: 0, percent: 0 };
  return { pc, mobile, other };
});

const topReferer = computed(() => getTopItem(referers.value));
const topPage = computed(() => getTopItem(pages.value));
const topEntry = computed(() => getTopItem(entryPages.value));
const topBrowser = computed(() => getTopItem(browsers.value));
const topCity = computed(() => getTopItem(locations.value));

const refererItems = computed(() => referers.value || []);
const pageItems = computed(() => pages.value || []);
const entryItems = computed(() => entryPages.value || []);
const browserItems = computed(() => browsers.value || []);
const cityItems = computed(() => locations.value || []);

onMounted(() => {
  initWindowFromPreference();
  loadWebsites();
});

onBeforeUnmount(() => {
  stopAutoRefresh();
});

watch(currentWebsiteId, (value) => {
  if (value) {
    saveUserPreference('selectedWebsite', value);
  }
  loadRealtime();
  restartAutoRefresh();
});

watch(currentWindow, (value) => {
  saveUserPreference('realtimeWindow', String(value));
  loadRealtime();
  restartAutoRefresh();
});

function initWindowFromPreference() {
  const queryWindow = getWindowFromQuery();
  const savedWindow = parseInt(getUserPreference('realtimeWindow', '30'), 10);
  const preferred = Number.isFinite(queryWindow) ? queryWindow : savedWindow;
  if (windowOptions.includes(preferred)) {
    currentWindow.value = preferred;
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

async function loadRealtime() {
  if (!currentWebsiteId.value) {
    return;
  }
  try {
    const data: RealtimeStats = await fetchRealtimeStats(currentWebsiteId.value, currentWindow.value);
    activeCount.value = data.activeCount || 0;
    activeSeries.value = data.activeSeries || [];
    deviceBreakdown.value = data.deviceBreakdown || [];
    referers.value = data.referers || [];
    pages.value = data.pages || [];
    entryPages.value = data.entryPages || [];
    browsers.value = data.browsers || [];
    locations.value = data.locations || [];
  } catch (error) {
    console.error('加载实时数据失败:', error);
  }
}

function startAutoRefresh() {
  if (refreshTimer) {
    return;
  }
  refreshTimer = window.setInterval(() => {
    loadRealtime();
  }, 30000);
}

function stopAutoRefresh() {
  if (!refreshTimer) {
    return;
  }
  window.clearInterval(refreshTimer);
  refreshTimer = null;
}

function restartAutoRefresh() {
  stopAutoRefresh();
  if (currentWebsiteId.value) {
    startAutoRefresh();
  }
}

function setWindow(value: number) {
  if (currentWindow.value === value) {
    return;
  }
  currentWindow.value = value;
}

function getWindowFromQuery() {
  const params = new URLSearchParams(window.location.search || '');
  const raw = params.get('window');
  if (!raw) {
    return Number.NaN;
  }
  const parsed = parseInt(raw, 10);
  return Number.isFinite(parsed) ? parsed : Number.NaN;
}

function getTopItem(items: RealtimeSeriesItem[] = []) {
  if (!items.length) {
    return { name: '-', count: 0, percent: 0 };
  }
  return items[0];
}

function formatCount(value: number) { return Number(value || 0).toLocaleString('zh-CN'); }

function formatPercent(value: number) { return `${(Number(value || 0) * 100).toFixed(2)}%`; }
</script>

<style scoped lang="scss">
:global(.realtime-page) {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.realtime-range {
  display: inline-flex;
  gap: 8px;
  padding: 6px;
  border-radius: 999px;
  background: var(--panel-muted);
  border: 1px solid var(--border);
  box-shadow: var(--shadow-soft);
}

.realtime-range-btn {
  border: none;
  background: transparent;
  padding: 6px 14px;
  border-radius: 999px;
  font-weight: 600;
  font-size: 12px;
  color: var(--muted);
  cursor: pointer;
}

.realtime-range-btn.active {
  background: linear-gradient(135deg, var(--primary), var(--primary-strong));
  color: white;
  box-shadow: 0 8px 16px rgba(var(--primary-color-rgb), 0.28);
}

.realtime-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 18px;
}

.realtime-card {
  display: flex;
  flex-direction: column;
  gap: 14px;
  min-height: 260px;
}

.realtime-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.realtime-card-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 700;
}

.realtime-metric {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.realtime-value {
  font-size: 28px;
  font-weight: 700;
}

.realtime-mini-bars {
  display: grid;
  grid-auto-flow: column;
  grid-auto-columns: 1fr;
  gap: 4px;
  height: 28px;
  align-items: end;
}

.realtime-mini-bars span {
  display: block;
  width: 100%;
  min-height: 6px;
  border-radius: 6px;
  background: var(--panel-muted);
}

.realtime-mini-bars span.active {
  background: linear-gradient(180deg, rgba(30, 123, 255, 0.65), rgba(30, 123, 255, 0.15));
}

.realtime-subtitle {
  font-size: 12px;
  color: var(--muted);
}

.realtime-device-cards {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}

.realtime-device-card {
  background: linear-gradient(160deg, rgba(255, 255, 255, 0.7), var(--panel-muted));
  border-radius: 14px;
  border: 1px solid var(--border);
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

:global(body.dark-mode) .realtime-device-card {
  background: linear-gradient(160deg, rgba(255, 255, 255, 0.04), var(--panel-muted));
}

.realtime-device-icon {
  font-size: 18px;
  color: var(--primary);
}

.realtime-device-label {
  font-size: 12px;
  color: var(--muted);
}

.realtime-device-count {
  font-size: 18px;
  font-weight: 700;
}

.realtime-device-rate {
  font-size: 12px;
  color: var(--muted);
}

.realtime-top {
  display: flex;
  align-items: center;
  gap: 12px;
}

.realtime-rank {
  padding: 4px 8px;
  border-radius: 999px;
  background: rgba(245, 158, 11, 0.18);
  color: var(--accent);
  font-size: 11px;
  font-weight: 700;
}

.realtime-top-title {
  flex: 1;
  font-weight: 600;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.realtime-top-meta {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.realtime-top-count {
  font-size: 20px;
  font-weight: 700;
}

.realtime-top-rate {
  font-size: 12px;
  color: var(--muted);
}

.realtime-table th,
.realtime-table td {
  padding: 10px 12px;
}

.realtime-count-col {
  text-align: right;
}

@media (max-width: 1200px) {
  .realtime-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 900px) {
  .realtime-grid {
    grid-template-columns: 1fr;
  }

  .realtime-device-cards {
    grid-template-columns: 1fr;
  }
}
</style>
