<template>
    <!-- <div class="dash-card"> -->
    <div ref="chartContainer" v-show="showChart" class="chart-wrapper">
  <VChart class="chart" :option="option" :key="chartKey"/>
    </div>
  <!-- </div> -->
</template>

<script setup lang="ts">
import { use } from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { PieChart } from "echarts/charts";
import { TitleComponent, TooltipComponent, LegendComponent, GridComponent, DatasetComponent } from "echarts/components";
import VChart, { THEME_KEY } from "vue-echarts";
import { ref, provide, watch, onMounted, nextTick } from "vue";
import { clientStore } from "@/stores/client.ts";
import type { RefSymbol } from "@vue/reactivity";


const chartKey = ref(0);
// const chartReady = ref(false)
const chartContainer = ref<HTMLElement | null>(null);

const showChart = ref(false);

use([CanvasRenderer, PieChart, TitleComponent, TooltipComponent, LegendComponent, GridComponent, DatasetComponent]);
provide(THEME_KEY, "dark");


    const option = ref({
    title: {
        text: clientStore().currClient.name +"'s Message stats",
        left: "center",
    },
    tooltip: {
        trigger: "item",
        formatter: "{a} <br/>{b} : {c} ({d}%)",
    },
    legend: {
        orient: "vertical",
        left: "left",
        data: ["QUEUED", "SENT"],
    },
    series: [
        {
        name: clientStore().currClient.name,
        type: "pie",
        radius: "55%",
        center: ["50%", "60%"],
        data: [
            { value: clientStore().currClient.messageQueue.length, name: "QUEUED" },
            { value: clientStore().currClient.allMessagesSent.length, name: "SENT" },
        ],
        emphasis: {
            itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: "rgba(0, 0, 0, 0.5)",
            },
        },
        },
    ],
    });

watch(() => [clientStore().currClient.name, clientStore().currClient.messageQueue.length, clientStore().currClient.allMessagesSent.length], () => {
 
 option.value = {
    title: {
        text: clientStore().currClient.name +"'s Message stats",
        left: "center",
    },
    tooltip: {
        trigger: "item",
        formatter: "{a} <br/>{b} : {c} ({d}%)",
    },
    legend: {
        orient: "vertical",
        left: "left",
        data: ["QUEUED", "SENT"],
    },
    series: [
        {
        name: clientStore().currClient.name,
        type: "pie",
        radius: "55%",
        center: ["50%", "60%"],
        data: [
            { value: clientStore().currClient.messageQueue.length, name: "QUEUED" },
            { value: clientStore().currClient.allMessagesSent.length, name: "SENT" },
        ],
        emphasis: {
            itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: "rgba(0, 0, 0, 0.5)",
            },
        },
        },
    ],
    }
 chartKey.value += 1; // Triggers re-render
}, { immediate: true });
onMounted(async () => {
  await nextTick(() => {
  }).then(() => {
    window.dispatchEvent(new Event('resize'))
    showChart.value = true;
    chartKey.value++
  });

});

</script>

<style scoped>
.dash-card {
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 12px;
  padding: 20px;
  margin: 20px auto;
  width: 5000px;
}

.chart {
  width: 100%;
  height: 100%;

}
.chart-wrapper {
  width: 100%;
  min-height: 400px;
  min-width: 500px;
}
</style>