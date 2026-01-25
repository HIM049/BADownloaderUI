<template>
    <body class="w-[60%] mx-auto my-5" >
        <var-paper :elevation="2" >
            <div class="flex items-center justify-between">
                <h2 class="text-xl font-bold ml-5 my-3">{{ title }}</h2>
                <div><slot name="actions"/></div>
            </div>
            <slot />
        </var-paper>
    </body>
</template>

<script setup>
const props = defineProps(['title'])

</script>