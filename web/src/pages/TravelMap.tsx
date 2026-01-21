import { useEffect, useRef } from "react";

export default function TravelMap() {
  const chartRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    let chartInstance: any = null;

    const handleResize = () => {
      chartInstance?.resize();
    };

    const initChart = async () => {
      if (!(window as any).echarts) {
        console.error("ECharts is not loaded");
        return;
      }

      if (!chartRef.current) return;

      try {
        // Fetch China Map Data
        const response = await fetch(
          "https://geo.datav.aliyun.com/areas_v3/bound/100000_full.json",
        );
        const chinaJson = await response.json();

        // Register Map
        (window as any).echarts.registerMap("china", chinaJson);

        // Initialize Chart
        chartInstance = (window as any).echarts.init(chartRef.current);

        const data = [
          { name: "北京", value: [116.407387, 39.904179] },
          { name: "南昌", value: [115.892151, 28.676493] },
          { name: "西安", value: [108.94866, 34.22245] },
          { name: "华山", value: [110.08752, 34.56608] },
          { name: "银川", value: [106.206479, 38.502621] },
          { name: "厦门", value: [118.103886, 24.489231] },
          { name: "武夷山", value: [118.036655, 27.756515] },
          { name: "东莞", value: [113.760234, 23.051271] },
          { name: "广州", value: [113.30765, 23.120049] },
          { name: "上海", value: [121.487899, 31.249162] },
          { name: "无锡", value: [120.318665, 31.501063] },
          { name: "苏州", value: [120.619907, 31.317987] },
          { name: "杭州", value: [120.219375, 30.259244] },
          { name: "绍兴", value: [120.592467, 30.002365] },
          { name: "嘉兴", value: [120.760428, 30.773992] },
          { name: "湖州", value: [120.137243, 30.877925] },
          { name: "中卫", value: [105.196754, 37.521124] },
          { name: "阿拉善盟左旗", value: [105.706422, 38.844814] },
          { name: "洪湖", value: [113.461212, 29.827365] },
          { name: "武汉", value: [114.3162, 30.581084] },
          { name: "咸宁", value: [114.300061, 29.880657] },
          { name: "泗洪", value: [118.22861, 33.40972] },
          { name: "日照", value: [119.52685, 35.41691] },
          { name: "香港", value: [114.1095, 22.3964] },
        ];
        const selectBg = "#fff";
        const regions = [
          { name: "北京市", itemStyle: { areaColor: selectBg } },
          { name: "上海市", itemStyle: { areaColor: selectBg } },
          { name: "湖北省", itemStyle: { areaColor: selectBg } },
          { name: "广东省", itemStyle: { areaColor: selectBg } },
          { name: "福建省", itemStyle: { areaColor: selectBg } },
          { name: "浙江省", itemStyle: { areaColor: selectBg } },
          { name: "江苏省", itemStyle: { areaColor: selectBg } },
          { name: "宁夏回族自治区", itemStyle: { areaColor: selectBg } },
          { name: "内蒙古自治区", itemStyle: { areaColor: selectBg } },
          { name: "山东省", itemStyle: { areaColor: selectBg } },
          { name: "江西省", itemStyle: { areaColor: selectBg } },
          { name: "陕西省", itemStyle: { areaColor: selectBg } },
          { name: "香港特别行政区", itemStyle: { areaColor: selectBg } },
        ];

        const option = {
          backgroundColor: "#fff", // tailwind gray-100
          title: {
            text: "山海漫记，皆是旅途",
            left: "center",
            top: 20,
            textStyle: {
              color: "#374151", // tailwind gray-700
              fontSize: 20,
            },
          },
          geo: {
            map: "china",
            roam: true,
            label: {
              emphasis: {
                show: false,
              },
            },
            regions: regions,
            itemStyle: {
              normal: {
                areaColor: "#e5e7eb", // tailwind gray-200
                borderColor: "#9ca3af", // tailwind gray-400
              },
              emphasis: {
                areaColor: "#d1d5db", // tailwind gray-300
              },
            },
          },
          tooltip: {
            trigger: "item",
          },
          series: [
            {
              name: "足迹",
              type: "effectScatter",
              coordinateSystem: "geo",
              data: data,
              symbol: "circle",
              symbolSize: 6,
              // showEffectOn: "render",
              // rippleEffect: {
              //   brushType: "stroke",
              //   scale: 3,
              // },
              // hoverAnimation: true,
              label: {
                formatter: "{b}",
                position: "right",
                show: true,
                color: "#1f2937", // tailwind gray-800
                fontSize: 10,
              },
              itemStyle: {
                color: "#60a5fa", // tailwind blue-400
                shadowBlur: 10,
                shadowColor: "rgba(96, 165, 250, 0.5)",
              },
              zlevel: 1,
            },
          ],
        };

        chartInstance.setOption(option);
        window.addEventListener("resize", handleResize);
      } catch (error) {
        console.error("Failed to load map data", error);
      }
    };

    // If echarts is not loaded yet (e.g. async script), we might want to wait or retry.
    // For now, assume it's loaded as we put it in head without defer/async or before main.
    initChart();

    return () => {
      window.removeEventListener("resize", handleResize);
      chartInstance?.dispose();
    };
  }, []);

  return (
    <div className="w-full h-[800px]">
      <div ref={chartRef} className="w-full h-full bg-white overflow-hidden" />
    </div>
  );
}
