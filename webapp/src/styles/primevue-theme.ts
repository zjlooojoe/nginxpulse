import { definePreset } from '@primeuix/themes';
import Aura from '@primeuix/themes/aura';

const primaryScale = {
  50: 'color-mix(in srgb, var(--primary), #ffffff 92%)',
  100: 'color-mix(in srgb, var(--primary), #ffffff 85%)',
  200: 'color-mix(in srgb, var(--primary), #ffffff 70%)',
  300: 'color-mix(in srgb, var(--primary), #ffffff 55%)',
  400: 'color-mix(in srgb, var(--primary), #ffffff 30%)',
  500: 'var(--primary)',
  600: 'var(--primary-strong)',
  700: 'color-mix(in srgb, var(--primary-strong), #000000 12%)',
  800: 'color-mix(in srgb, var(--primary-strong), #000000 24%)',
  900: 'color-mix(in srgb, var(--primary-strong), #000000 40%)',
  950: 'color-mix(in srgb, var(--primary-strong), #000000 55%)',
};

const surfaceLight = {
  0: 'var(--panel)',
  50: 'var(--bg)',
  100: 'var(--panel-muted)',
  200: 'var(--border)',
  300: 'color-mix(in srgb, var(--border), var(--text) 12%)',
  400: 'color-mix(in srgb, var(--muted), #ffffff 20%)',
  500: 'var(--muted)',
  600: 'color-mix(in srgb, var(--muted), var(--text) 45%)',
  700: 'color-mix(in srgb, var(--text), #ffffff 20%)',
  800: 'color-mix(in srgb, var(--text), #000000 5%)',
  900: 'var(--text)',
  950: 'color-mix(in srgb, var(--text), #000000 25%)',
};

const surfaceDark = {
  0: 'var(--text)',
  50: 'color-mix(in srgb, var(--text), #000000 18%)',
  100: 'var(--muted)',
  200: 'color-mix(in srgb, var(--muted), #000000 25%)',
  300: 'color-mix(in srgb, var(--muted), #000000 40%)',
  400: 'color-mix(in srgb, var(--border), #ffffff 10%)',
  500: 'var(--border)',
  600: 'color-mix(in srgb, var(--border), #000000 20%)',
  700: 'var(--panel)',
  800: 'var(--panel-muted)',
  900: 'var(--bg)',
  950: 'color-mix(in srgb, var(--bg), #000000 20%)',
};

const nginxPulsePreset = definePreset(Aura, {
  semantic: {
    primary: primaryScale,
    colorScheme: {
      light: {
        surface: surfaceLight,
        primary: {
          color: 'var(--primary)',
          contrastColor: '#ffffff',
          hoverColor: 'var(--primary-strong)',
          activeColor: 'color-mix(in srgb, var(--primary-strong), #000000 15%)',
        },
        highlight: {
          background: 'var(--primary-soft)',
          focusBackground: 'color-mix(in srgb, var(--primary-soft), var(--primary) 20%)',
          color: 'var(--primary-strong)',
          focusColor: 'var(--primary-strong)',
        },
        text: {
          color: 'var(--text)',
          hoverColor: 'var(--text)',
          mutedColor: 'var(--muted)',
          hoverMutedColor: 'color-mix(in srgb, var(--muted), var(--text) 20%)',
        },
        content: {
          background: 'var(--panel)',
          hoverBackground: 'var(--panel-muted)',
          borderColor: 'var(--border)',
          color: 'var(--text)',
          hoverColor: 'var(--text)',
        },
        formField: {
          background: 'var(--input-bg)',
          disabledBackground: 'var(--panel-muted)',
          filledBackground: 'var(--panel-muted)',
          filledHoverBackground: 'var(--panel-muted)',
          filledFocusBackground: 'var(--panel-muted)',
          borderColor: 'var(--border)',
          hoverBorderColor: 'color-mix(in srgb, var(--border), var(--primary) 25%)',
          focusBorderColor: 'var(--primary)',
          invalidBorderColor: '#ef4444',
          color: 'var(--text)',
          disabledColor: 'var(--muted)',
          placeholderColor: 'var(--muted)',
          invalidPlaceholderColor: '#ef4444',
          floatLabelColor: 'var(--muted)',
          floatLabelFocusColor: 'var(--primary)',
          floatLabelActiveColor: 'var(--muted)',
          floatLabelInvalidColor: '#ef4444',
          iconColor: 'var(--muted)',
          shadow: 'var(--shadow-soft)',
        },
        overlay: {
          select: {
            background: 'var(--panel)',
            borderColor: 'var(--border)',
            color: 'var(--text)',
          },
          popover: {
            background: 'var(--panel)',
            borderColor: 'var(--border)',
            color: 'var(--text)',
          },
          modal: {
            background: 'var(--panel)',
            borderColor: 'var(--border)',
            color: 'var(--text)',
          },
        },
      },
      dark: {
        surface: surfaceDark,
        primary: {
          color: 'var(--primary)',
          contrastColor: 'var(--bg)',
          hoverColor: 'var(--primary-strong)',
          activeColor: 'color-mix(in srgb, var(--primary-strong), #000000 15%)',
        },
        highlight: {
          background: 'var(--primary-soft)',
          focusBackground: 'color-mix(in srgb, var(--primary-soft), var(--primary) 20%)',
          color: 'var(--text)',
          focusColor: 'var(--text)',
        },
        text: {
          color: 'var(--text)',
          hoverColor: 'var(--text)',
          mutedColor: 'var(--muted)',
          hoverMutedColor: 'color-mix(in srgb, var(--muted), #ffffff 20%)',
        },
        content: {
          background: 'var(--panel)',
          hoverBackground: 'var(--panel-muted)',
          borderColor: 'var(--border)',
          color: 'var(--text)',
          hoverColor: 'var(--text)',
        },
        formField: {
          background: 'var(--input-bg)',
          disabledBackground: 'var(--panel-muted)',
          filledBackground: 'var(--panel-muted)',
          filledHoverBackground: 'var(--panel-muted)',
          filledFocusBackground: 'var(--panel-muted)',
          borderColor: 'var(--border)',
          hoverBorderColor: 'color-mix(in srgb, var(--border), var(--primary) 25%)',
          focusBorderColor: 'var(--primary)',
          invalidBorderColor: '#f87171',
          color: 'var(--text)',
          disabledColor: 'var(--muted)',
          placeholderColor: 'var(--muted)',
          invalidPlaceholderColor: '#f87171',
          floatLabelColor: 'var(--muted)',
          floatLabelFocusColor: 'var(--primary)',
          floatLabelActiveColor: 'var(--muted)',
          floatLabelInvalidColor: '#f87171',
          iconColor: 'var(--muted)',
          shadow: 'var(--shadow-soft)',
        },
        overlay: {
          select: {
            background: 'var(--panel)',
            borderColor: 'var(--border)',
            color: 'var(--text)',
          },
          popover: {
            background: 'var(--panel)',
            borderColor: 'var(--border)',
            color: 'var(--text)',
          },
          modal: {
            background: 'var(--panel)',
            borderColor: 'var(--border)',
            color: 'var(--text)',
          },
        },
      },
    },
  },
});

export default nginxPulsePreset;
