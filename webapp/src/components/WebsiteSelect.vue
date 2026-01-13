<template>
  <div class="select-group">
    <label v-if="label" class="select-label" :for="selectId">{{ label }}</label>
    <Dropdown
      :inputId="selectId"
      v-model="selectedValue"
      class="website-dropdown"
      :options="websites"
      optionLabel="name"
      optionValue="id"
      :disabled="disabled"
      :placeholder="placeholderText"
      :loading="loading"
      :emptyMessage="emptyText"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { WebsiteInfo } from '@/api/types';

const props = withDefaults(
  defineProps<{
    modelValue: string;
    websites: WebsiteInfo[];
    label?: string;
    id?: string;
    loading?: boolean;
    loadingText?: string;
    emptyText?: string;
  }>(),
  {
    label: '站点',
    id: 'website-selector',
    loading: false,
    loadingText: '加载中...',
    emptyText: '没有可用的网站',
  }
);

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void;
}>();

const selectId = computed(() => props.id || 'website-selector');
const disabled = computed(() => props.loading || props.websites.length === 0);
const selectedValue = computed({
  get: () => props.modelValue,
  set: (value: string) => emit('update:modelValue', value),
});
const placeholderText = computed(() => {
  if (props.loading) {
    return props.loadingText;
  }
  if (!props.websites.length) {
    return props.emptyText;
  }
  return '请选择站点';
});
</script>
