<script setup>
import navBar from './components/navbar.vue'
import mainPage from './components/main_page.vue'
import settingPage from './components/setting_page.vue'
import { ref } from 'vue'


const isSetting = ref(false)

function switchPages() {
  isSetting.value = isSetting.value ? false : true
}
</script>

<template>

  <navBar @page-switch="switchPages"/>
  <section id="app-function">
    <keep-alive>
      <component :is="isSetting ? settingPage : mainPage"></component>
    </keep-alive>
  </section>

</template>
