import { useEffect, useRef, useState } from "react";
import { footprintsApi } from "@/service";
import type { FootprintItem } from "@/types/openapi";

declare global {
  interface Window {
    FootprintMap: {
      initWithData: (
        container: string | HTMLElement,
        data: { locations: unknown[] },
      ) => void;
    };
  }
}

function toXiaoTenFormat(items: FootprintItem[]) {
  return {
    locations: items.map((fp) => ({
      name: fp.name,
      coordinates: `${fp.longitude},${fp.latitude}`,
      description: fp.description,
      date: fp.date,
      url: fp.url,
      urlLabel: fp.url_label || "查看相关内容",
      photos: (fp.photos || []).map((p) => p.src),
      categories: fp.categories || [],
      markerPreset: fp.marker_color || undefined,
    })),
  };
}

export default function TravelMap() {
  const containerRef = useRef<HTMLDivElement>(null);
  const [error, setError] = useState<string>("");

  useEffect(() => {
    let destroyed = false;

    async function load() {
      try {
        const resp = await footprintsApi();
        if (destroyed) return;
        const data = toXiaoTenFormat(resp.footprints || []);
        if (containerRef.current && window.FootprintMap) {
          window.FootprintMap.initWithData(containerRef.current, data);
        }
      } catch (e: unknown) {
        if (!destroyed) setError(e instanceof Error ? e.message : "加载失败");
      }
    }

    load();
    return () => {
      destroyed = true;
    };
  }, []);

  return (
    <>
      <title>山海漫记</title>
      <h1 className="text-center text-lg font-medium py-4">山海漫记，皆是旅途</h1>
      <div
        ref={containerRef}
        className="footprint-map footprint-map--loading"
        style={{ width: "100%", height: "800px" }}
      />
      {error && (
        <div className="footprint-map__error">加载失败: {error}</div>
      )}
    </>
  );
}
