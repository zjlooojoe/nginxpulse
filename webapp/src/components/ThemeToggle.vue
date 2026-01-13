<template>
  <button class="theme-toggle" type="button" title="切换主题" :aria-pressed="isDark" @click="toggle">
    <span class="light-icon">🌙</span>
    <span class="dark-icon">☀️</span>
  </button>
</template>

<script setup lang="ts">
import { computed, inject } from 'vue';

type ThemeContext = {
  isDark: { value: boolean };
  toggle: () => void;
}

const theme = inject<ThemeContext>('theme', null);
const isDark = computed(() => theme?.isDark.value ?? false);

function toggle() {
  theme?.toggle();
}
</script>
