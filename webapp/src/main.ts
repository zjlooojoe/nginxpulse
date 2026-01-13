import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import PrimeVue from 'primevue/config';
import Button from 'primevue/button';
import Checkbox from 'primevue/checkbox';
import DatePicker from 'primevue/datepicker';
import Dropdown from 'primevue/dropdown';
import InputNumber from 'primevue/inputnumber';
import InputText from 'primevue/inputtext';
import nginxPulsePreset from './styles/primevue-theme';
import 'primeicons/primeicons.css';

import './styles/vendor.scss';
import './styles/index.scss';

const app = createApp(App);
app.use(router);
app.use(PrimeVue, {
  ripple: true,
  theme: {
    preset: nginxPulsePreset,
    options: {
      darkModeSelector: '.dark-mode',
    },
  },
  locale: {
    startsWith: '开始于',
    contains: '包含',
    notContains: '不包含',
    endsWith: '结束于',
    equals: '等于',
    notEquals: '不等于',
    noFilter: '无筛选',
    lt: '小于',
    lte: '小于等于',
    gt: '大于',
    gte: '大于等于',
    dateIs: '日期是',
    dateIsNot: '日期不是',
    dateBefore: '日期早于',
    dateAfter: '日期晚于',
    clear: '清除',
    apply: '应用',
    matchAll: '匹配全部',
    matchAny: '匹配任意',
    addRule: '添加规则',
    removeRule: '移除规则',
    accept: '确定',
    reject: '取消',
    choose: '选择',
    upload: '上传',
    cancel: '取消',
    completed: '完成',
    pending: '进行中',
    fileSizeTypes: ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'],
    dayNames: ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六'],
    dayNamesShort: ['周日', '周一', '周二', '周三', '周四', '周五', '周六'],
    dayNamesMin: ['日', '一', '二', '三', '四', '五', '六'],
    monthNames: ['一月', '二月', '三月', '四月', '五月', '六月', '七月', '八月', '九月', '十月', '十一月', '十二月'],
    monthNamesShort: ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月'],
    chooseYear: '选择年份',
    chooseMonth: '选择月份',
    chooseDate: '选择日期',
    prevDecade: '上十年',
    nextDecade: '下十年',
    prevYear: '上一年',
    nextYear: '下一年',
    prevMonth: '上一月',
    nextMonth: '下一月',
    prevHour: '上一小时',
    nextHour: '下一小时',
    prevMinute: '上一分钟',
    nextMinute: '下一分钟',
    prevSecond: '上一秒',
    nextSecond: '下一秒',
    am: '上午',
    pm: '下午',
    today: '今天',
    weekHeader: '周',
    firstDayOfWeek: 1,
    showMonthAfterYear: true,
    dateFormat: 'yy-mm-dd',
    weak: '弱',
    medium: '中等',
    strong: '强',
    passwordPrompt: '请输入密码',
    emptyFilterMessage: '无可用选项',
    searchMessage: '{0} 个结果可用',
    selectionMessage: '{0} 项已选择',
    emptySelectionMessage: '未选择',
    emptySearchMessage: '没有找到结果',
    emptyMessage: '无数据',
    aria: {
      trueLabel: '是',
      falseLabel: '否',
      nullLabel: '未选择',
      close: '关闭',
      previous: '上一项',
      next: '下一项',
      navigation: '导航',
      scrollTop: '滚动到顶部',
      moveUp: '上移',
      moveDown: '下移',
      moveToTarget: '移动到目标',
      moveToSource: '移动到来源',
      moveAllToTarget: '全部移动到目标',
      moveAllToSource: '全部移动到来源',
      pageLabel: '页',
      firstPageLabel: '第一页',
      lastPageLabel: '最后一页',
      nextPageLabel: '下一页',
      prevPageLabel: '上一页',
      rowsPerPageLabel: '每页行数',
      jumpToPageDropdownLabel: '跳转到页码（下拉）',
      jumpToPageInputLabel: '跳转到页码（输入）',
      selectRow: '选择行',
      unselectRow: '取消选择行',
      expandRow: '展开行',
      collapseRow: '折叠行',
      showFilterMenu: '显示筛选菜单',
      hideFilterMenu: '隐藏筛选菜单',
      filterOperator: '筛选运算符',
      filterConstraint: '筛选条件',
      editRow: '编辑行',
      saveEdit: '保存编辑',
      cancelEdit: '取消编辑',
      listView: '列表视图',
      gridView: '网格视图',
      slide: '幻灯片',
      slideNumber: '{slideNumber}',
      zoomImage: '缩放图片',
      zoomIn: '放大',
      zoomOut: '缩小',
      rotateRight: '向右旋转',
      rotateLeft: '向左旋转',
    },
  },
});
app.component('Button', Button);
app.component('Checkbox', Checkbox);
app.component('DatePicker', DatePicker);
app.component('Dropdown', Dropdown);
app.component('InputNumber', InputNumber);
app.component('InputText', InputText);
app.mount('#app');
