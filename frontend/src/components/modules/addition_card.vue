<template>
    <body class="add-body" >
        <var-paper :elevation="2" >
            <h2>{{ title }}</h2>
            <slot />
        </var-paper>
    </body>
</template>

<script setup>
const props = defineProps(['title'])

</script>

<style>
    .add-body {
        width: 60%;
        margin: 0px auto 20px auto;;
    }
</style>