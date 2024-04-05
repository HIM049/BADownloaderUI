<template>
  <MainPage />
</template>

<script setup>
import { StyleProvider, Themes } from '@varlet/ui'
import MainPage from './components/main_page.vue'

const lightTheme = {
    '--color-primary': '#9c4146',
    '--color-on-primary': '#ffffff',
    '--color-primary-container': '#ffdad9',
    '--color-on-primary-container': '#40000a',
    '--color-info': '#775656',
    '--color-on-info': '#ffffff',
    '--color-info-container': '#ffdad9',
    '--color-on-info-container': '#2c1516',
    '--color-warning': '#755a2f',
    '--color-on-warning': '#ffffff',
    '--color-warning-container': '#ffddaf',
    '--color-on-warning-container': '#281800',
    '--color-danger': '#ba1a1a',
    '--color-on-danger': '#ffffff',
    '--color-danger-container': '#ffdad6',
    '--color-on-danger-container': '#410002',
    '--color-body': '#fffbff',
    '--color-text': '#201a1a',
    '--color-on-surface-variant': '#524343',
    '--color-outline': '#857373',
    '--color-inverse-surface': '#362f2f',
    '--paper-background': '#fbf1f1'
}
const darkTheme = {
    '--color-primary': '#ffb3b4',
    '--color-on-primary': '#5f131c',
    '--color-primary-container': '#7e2a30',
    '--color-on-primary-container': '#ffdad9',
    '--color-info': '#e6bdbc',
    '--color-on-info': '#44292a',
    '--color-info-container': '#5d3f3f',
    '--color-on-info-container': '#ffdad9',
    '--color-warning': '#e5c18d',
    '--color-on-warning': '#422c05',
    '--color-warning-container': '#5b421a',
    '--color-on-warning-container': '#ffddaf',
    '--color-danger': '#ffb4ab',
    '--color-on-danger': '#690005',
    '--color-danger-container': '#93000a',
    '--color-on-danger-container': '#ffb4ab',
    '--color-body': '#201a1a',
    '--color-text': '#ece0df',
    '--color-on-surface-variant': '#d7c1c1',
    '--color-outline': '#a08c8c',
    '--color-inverse-surface': '#ece0df'
}

const mergedThemeLight = { ...Themes.md3Light, ...lightTheme };
StyleProvider(mergedThemeLight)

// const mergedThemeDark = { ...Themes.md3Dark, ...darkTheme };
// StyleProvider(mergedThemeDark)
</script>

<style>
div.var-paper {
    padding: 20px;
    border-radius: 20px;
}
</style>