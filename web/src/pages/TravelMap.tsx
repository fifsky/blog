import { useEffect, useRef, useState } from "react";
import { settingChinaMapApi, footprintsApi, cityPhotosApi } from "@/service";
import Lightbox from "yet-another-react-lightbox";
import Captions from "yet-another-react-lightbox/plugins/captions";
import Thumbnails from "yet-another-react-lightbox/plugins/thumbnails";
import Zoom from "yet-another-react-lightbox/plugins/zoom";
import "yet-another-react-lightbox/styles.css";
import "yet-another-react-lightbox/plugins/captions.css";
import "yet-another-react-lightbox/plugins/thumbnails.css";
import { TravelPhoto } from "@/types/openapi";

export default function TravelMap() {
  const chartRef = useRef<HTMLDivElement>(null);
  const [open, setOpen] = useState(false);
  const [photos, setPhotos] = useState<TravelPhoto[]>([]);

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
        // Fetch China Map Data and Footprints
        const [chinaJson, footprintsResp] = await Promise.all([
          settingChinaMapApi(),
          footprintsApi(),
        ]);

        // Register Map
        (window as any).echarts.registerMap("china", chinaJson);

        // Initialize Chart
        chartInstance = (window as any).echarts.init(chartRef.current);

        // Store cities data for later use
        const cities = footprintsResp.cities || [];

        // Map regions (provinces) for highlighting
        const regionsList = (footprintsResp.provinces || []).map((p) => p.name);

        // Map footprints (cities) for scatter points
        const data = cities.map((city) => ({
          name: city.name,
          value: [parseFloat(city.longitude), parseFloat(city.latitude)],
          regionId: city.region_id,
        }));

        const selectBg = "#fff";
        const regions = regionsList.map((name) => ({
          name,
          itemStyle: { areaColor: selectBg },
        }));

        const option = {
          backgroundColor: "#fff",
          title: {
            text: "山海漫记，皆是旅途",
            left: "center",
            top: 20,
            textStyle: {
              color: "#374151",
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
                areaColor: "#e5e7eb",
                borderColor: "#9ca3af",
              },
              emphasis: {
                areaColor: "#d1d5db",
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
              label: {
                formatter: "{b}",
                position: "right",
                show: true,
                color: "#1f2937",
                fontSize: 10,
              },
              itemStyle: {
                color: "#60a5fa",
                shadowBlur: 10,
                shadowColor: "rgba(96, 165, 250, 0.5)",
                cursor: "pointer",
              },
              zlevel: 1,
            },
          ],
        };

        chartInstance.setOption(option);

        chartInstance.on("click", async (params: any) => {
          if (params.componentType === "series" && params.seriesName === "足迹") {
            const regionId = params.data?.regionId;
            if (typeof regionId === "number" && regionId > 0) {
              try {
                const resp = await cityPhotosApi({ region_id: regionId });
                const cityPhotos = resp.photos || [];
                if (cityPhotos.length > 0) {
                  setPhotos(cityPhotos);
                  setOpen(true);
                } else {
                  console.log("No photos for this city");
                }
              } catch (error) {
                console.error("Failed to load city photos", error);
              }
            }
          }
        });

        window.addEventListener("resize", handleResize);
      } catch (error) {
        console.error("Failed to load map data", error);
      }
    };

    initChart();

    return () => {
      window.removeEventListener("resize", handleResize);
      chartInstance?.dispose();
    };
  }, []);

  // Convert photos to lightbox slides format
  const slides = photos.map((photo) => ({
    src: photo.src,
    thumbnail: photo.thumbnail,
    title: photo.title,
    description: photo.description,
  }));

  return (
    <>
      <div className="w-full h-[800px]">
        <div ref={chartRef} className="w-full h-full bg-white overflow-hidden" />
      </div>
      <Lightbox
        open={open}
        close={() => setOpen(false)}
        slides={slides}
        plugins={[Captions, Thumbnails, Zoom]}
        render={{
          thumbnail: ({ slide, rect }) => {
            const s = slide as any;
            return (
              <img
                src={s.thumbnail || s.src}
                alt={s.title || ""}
                width={rect.width}
                height={rect.height}
                className="w-full h-full object-cover"
              />
            );
          },
        }}
      />
    </>
  );
}
