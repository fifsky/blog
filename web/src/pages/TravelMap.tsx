import { useEffect, useRef } from "react";
import { settingApi, settingChinaMapApi } from "@/service";

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
        // Fetch China Map Data and Settings
        const [chinaJson, settingResponse] = await Promise.all([
          settingChinaMapApi(),
          settingApi(),
        ]);

        // Register Map
        (window as any).echarts.registerMap("china", chinaJson);

        // Initialize Chart
        chartInstance = (window as any).echarts.init(chartRef.current);

        let data = [];
        let regionsList = [];

        if (settingResponse.kv?.map_footprints) {
          try {
            const parsed = JSON.parse(settingResponse.kv.map_footprints);
            if (Array.isArray(parsed)) data = parsed;
          } catch (e) {
            console.error("Failed to parse map_footprints", e);
          }
        }

        if (settingResponse.kv?.map_regions) {
          try {
            const parsed = JSON.parse(settingResponse.kv.map_regions);
            if (Array.isArray(parsed)) regionsList = parsed;
          } catch (e) {
            console.error("Failed to parse map_regions", e);
          }
        }

        const selectBg = "#fff";
        const regions = regionsList.map((name) => ({
          name,
          itemStyle: { areaColor: selectBg },
        }));

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
