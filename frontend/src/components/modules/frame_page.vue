<template>
    <body class="page-body" >
        <var-paper :elevation="2" >
            <div style="display: flex; align-items: center; justify-content: space-between;">
                <h2>{{ title }}</h2>
                <div><slot name="actions"/></div>
            </div>
            <slot />
        </var-paper>
    </body>
</template>

<script setup>
const props = defineProps(['title'])

</script>

<style>
    .page-body {
        width: 60%;
        margin: 20px auto;
    }
</style>