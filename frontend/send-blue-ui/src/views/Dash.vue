<template>
  <VChart class="chart" :option="option" />
</template>

<script setup lang="ts">
import { use } from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { PieChart } from "echarts/charts";
import { TitleComponent, TooltipComponent, LegendComponent } from "echarts/components";
import VChart, { THEME_KEY } from "vue-echarts";
import { ref, provide, watch } from "vue";

use([CanvasRenderer, PieChart, TitleComponent, TooltipComponent, LegendComponent]);

const props = defineProps({
    clientId: {
      type: String,
      required: true,
    },
    });
const leads = ref<string[]>([]);
async function fetchLeads(clientId: string) {
    try {
    const base = import.meta.env.VITE_API_BASE || ''
    const res = await fetch(`${base}/clients/leads?client_id=${clientId}`)
    if (!res.ok) throw new Error('Network')
    const data = await res.json()
    leads.value = data.leads || []
    } catch (e) {
    console.error('fetchLeads error', e)
    }
}
    watch(
      () => props.clientId,
      (newClientId) => {
        console.log("Client ID changed:", newClientId);
        fetchLeads(newClientId);
      }
        // You can add additional logic here to handle client ID changes
    );
watch(leads, (newLeads) => {
    console.log("Leads updated:", newLeads);
});
provide(THEME_KEY, "dark");

const option = ref({
  title: {
    text: "Traffic Sources",
    left: "center",
  },
  tooltip: {
    trigger: "item",
    formatter: "{a} <br/>{b} : {c} ({d}%)",
  },
  legend: {
    orient: "vertical",
    left: "left",
    data: ["Direct", "Email", "Ad Networks", "Video Ads", "Search Engines"],
  },
  series: [
    {
      name: "Traffic Sources",
      type: "pie",
      radius: "55%",
      center: ["50%", "60%"],
      data: [
        { value: 335, name: "Direct" },
        { value: 310, name: "Email" },
        { value: 234, name: "Ad Networks" },
        { value: 135, name: "Video Ads" },
        { value: 1548, name: "Search Engines" },
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
</script>

<style scoped>
.chart {
  height: 400px;
}
</style>