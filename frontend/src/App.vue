<template>
  <MainPage />
</template>

<script setup>
import MainPage from './app/main_page.vue'
import { ref, onBeforeMount } from 'vue'
import { StyleProvider, Themes } from '@varlet/ui'
import { GetTheme } from '../wailsjs/go/main/App'

const ThemeLightPink = {
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
    '--paper-background': '#fbf1f1',
    '--page-background': '#ffeaea',
    '--color-primary-font': '#fc729a'
}
// const darkPink = {
//     '--color-primary': '#ffb3b4',
//     '--color-on-primary': '#5f131c',
//     '--color-primary-container': '#7e2a30',
//     '--color-on-primary-container': '#ffdad9',
//     '--color-info': '#e6bdbc',
//     '--color-on-info': '#44292a',
//     '--color-info-container': '#5d3f3f',
//     '--color-on-info-container': '#ffdad9',
//     '--color-warning': '#e5c18d',
//     '--color-on-warning': '#422c05',
//     '--color-warning-container': '#5b421a',
//     '--color-on-warning-container': '#ffddaf',
//     '--color-danger': '#ffb4ab',
//     '--color-on-danger': '#690005',
//     '--color-danger-container': '#93000a',
//     '--color-on-danger-container': '#ffb4ab',
//     '--color-body': '#201a1a',
//     '--color-text': '#ece0df',
//     '--color-on-surface-variant': '#d7c1c1',
//     '--color-outline': '#a08c8c',
//     '--color-inverse-surface': '#ece0df'
// }

const ThemeLightBlue = {
    "--color-primary": "#006688",
    "--color-on-primary": "#ffffff",
    "--color-primary-container": "#c2e8ff",
    "--color-on-primary-container": "#001e2b",
    "--color-info": "#4e616d",
    "--color-on-info": "#ffffff",
    "--color-info-container": "#d1e5f3",
    "--color-on-info-container": "#091e28",
    "--color-warning": "#5f5a7d",
    "--color-on-warning": "#ffffff",
    "--color-warning-container": "#e5deff",
    "--color-on-warning-container": "#1c1736",
    "--color-danger": "#ba1a1a",
    "--color-on-danger": "#ffffff",
    "--color-danger-container": "#ffdad6",
    "--color-on-danger-container": "#410002",
    "--color-body": "#fbfcfe",
    "--color-text": "#191c1e",
    "--color-on-surface-variant": "#40484d",
    "--color-outline": "#71787d",
    "--color-inverse-surface": "#2e3133",
    '--paper-background': '#e0f3ff',
    '--page-background': '#c3e8ff',
    '--color-primary-font': '#fc729a'
}
//   "dark": {
//     "--color-primary": "#76d1ff",
//     "--color-on-primary": "#003548",
//     "--color-primary-container": "#004d67",
//     "--color-on-primary-container": "#c2e8ff",
//     "--color-info": "#b5c9d7",
//     "--color-on-info": "#20333d",
//     "--color-info-container": "#364954",
//     "--color-on-info-container": "#d1e5f3",
//     "--color-warning": "#c9c1ea",
//     "--color-on-warning": "#312c4c",
//     "--color-warning-container": "#484364",
//     "--color-on-warning-container": "#e5deff",
//     "--color-danger": "#ffb4ab",
//     "--color-on-danger": "#690005",
//     "--color-danger-container": "#93000a",
//     "--color-on-danger-container": "#ffb4ab",
//     "--color-body": "#191c1e",
//     "--color-text": "#e1e2e5",
//     "--color-on-surface-variant": "#c0c7cd",
//     "--color-outline": "#8a9297",
//     "--color-inverse-surface": "#e1e2e5"
//   }

const mergedThemeLight = { ...Themes.md3Light, ...ThemeLightPink };
StyleProvider(mergedThemeLight)

const Theme = ref();

onBeforeMount(() => {
    GetTheme().then(result => {
        switch (result) {
            case 'lightBlue':
                Theme.value = ThemeLightBlue;
                break;

            case 'lightPink':
            default:
                Theme.value = ThemeLightPink;
                break;
        }
    const mergedThemeLight = { ...Themes.md3Light, ...Theme.value };
    StyleProvider(mergedThemeLight)
    });
});


// const mergedThemeDark = { ...Themes.md3Dark, ...darkTheme };
// StyleProvider(mergedThemeDark)
</script>

<style>
div.var-paper {
    padding: 20px;
    border-radius: 20px;
}
</style>