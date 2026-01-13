import { createRouter, createWebHistory } from 'vue-router';
import OverviewPage from '@/pages/OverviewPage.vue';
import DailyPage from '@/pages/DailyPage.vue';
import RealtimePage from '@/pages/RealtimePage.vue';
import LogsPage from '@/pages/LogsPage.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'overview',
      component: OverviewPage,
      meta: {
        sidebarLabel: '近期活跃',
        sidebarHint: '15 分钟活跃访客',
        mainClass: '',
      },
    },
    {
      path: '/daily',
      name: 'daily',
      component: DailyPage,
      meta: {
        sidebarLabel: '日报说明',
        sidebarHint: '按天聚合统计并提供趋势解读',
        mainClass: 'daily-page',
      },
    },
    {
      path: '/realtime',
      name: 'realtime',
      component: RealtimePage,
      meta: {
        sidebarLabel: '实时概览',
        sidebarHint: '关注窗口内实时行为趋势',
        mainClass: 'realtime-page',
      },
    },
    {
      path: '/logs',
      name: 'logs',
      component: LogsPage,
      meta: {
        sidebarLabel: '日志查询',
        sidebarHint: '按条件过滤并分页浏览',
        mainClass: 'logs-page',
      },
    },
  ],
  scrollBehavior() {
    return { top: 0 };
  },
});

export default router;
