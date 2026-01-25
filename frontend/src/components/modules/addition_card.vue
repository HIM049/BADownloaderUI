<template>
    <div class="w-[60%] mx-auto mb-5">
        <var-paper :elevation="2" >
            <h2 class="text-xl font-bold ml-5 my-3">{{ title }}</h2>
            <slot />
        </var-paper>
    </div>
</template>

<script setup>
const props = defineProps(['title'])

</script>