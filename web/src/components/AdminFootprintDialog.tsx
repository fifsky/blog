import { Controller, useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Field, FieldLabel, FieldError, FieldGroup, FieldContent } from "./ui/field";
import { Input } from "./ui/input";
import { Textarea } from "./ui/textarea";
import { Button } from "./ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "./ui/dialog";
import { FootprintItem } from "@/types/openapi";
import { useEffect, useState, useRef, useCallback } from "react";
import { ossPresignApi } from "@/service";
import FileUploadCompact from "./file-upload/compact-upload";
import { FileWithPreview } from "@/hooks/use-file-upload";
import InputTags from "./custom/input-tag";

const MARKER_COLORS = [
  { value: "", label: "自动" },
  { value: "sunset", label: "日落", color: "#ff6b6b" },
  { value: "ocean", label: "海洋", color: "#4ecdc4" },
  { value: "violet", label: "紫罗兰", color: "#a55eea" },
  { value: "forest", label: "森林", color: "#26de81" },
  { value: "amber", label: "琥珀", color: "#fed330" },
  { value: "citrus", label: "柑橘", color: "#fd9644" },
];

const GRADIENT_MAP: Record<string, string> = {
  sunset: "linear-gradient(135deg, #ffb347, #ff6f61)",
  ocean: "linear-gradient(135deg, #06beb6, #48b1bf)",
  violet: "linear-gradient(135deg, #a18cd1, #fbc2eb)",
  forest: "linear-gradient(135deg, #5ee7df, #39a37c)",
  amber: "linear-gradient(135deg, #f6d365, #fda085)",
  citrus: "linear-gradient(135deg, #fdfb8f, #a1ffce)",
};

interface AdminFootprintDialogProps {
  isOpen: boolean;
  onClose: () => void;
  item: FootprintItem | undefined;
  onSubmit: (values: FormValues) => Promise<void>;
}

export interface FormValues {
  name: string;
  description?: string;
  longitude: string;
  latitude: string;
  date?: string;
  marker_color: string;
  categories: string[];
  url?: string;
  url_label?: string;
  photo_urls: string[];
}

const formSchema = z.object({
  name: z.string().min(1, "请输入地点名称"),
  longitude: z.string().min(1, "请选择或输入经度"),
  latitude: z.string().min(1, "请选择或输入纬度"),
});

export function AdminFootprintDialog({ isOpen, onClose, item, onSubmit }: AdminFootprintDialogProps) {
  const [loading, setLoading] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [pendingFiles, setPendingFiles] = useState<FileWithPreview[]>([]);
  const [existingPhotos, setExistingPhotos] = useState<string[]>([]);
  const [selectedColor, setSelectedColor] = useState("");
  const [categories, setCategories] = useState<string[]>([]);
  const [mapSearchKeyword, setMapSearchKeyword] = useState("");
  const mapContainerRef = useRef<HTMLDivElement>(null);
  const mapInstanceRef = useRef<any>(null);
  const markerRef = useRef<any>(null);
  const placeSearchRef = useRef<any>(null);
  const descRef = useRef<HTMLTextAreaElement>(null);
  const dateRef = useRef<HTMLInputElement>(null);
  const urlRef = useRef<HTMLInputElement>(null);
  const urlLabelRef = useRef<HTMLInputElement>(null);

  const isEdit = !!item?.id;

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: item?.name ?? "",
      longitude: item?.longitude ?? "",
      latitude: item?.latitude ?? "",
    },
    mode: "onChange",
  });

  const initMap = useCallback(() => {
    const AMap = (window as any).AMap;
    if (!mapContainerRef.current || !AMap) return;

    const lng = item?.longitude ? parseFloat(item.longitude) : 116.397428;
    const lat = item?.latitude ? parseFloat(item.latitude) : 39.90923;

    const map = new AMap.Map(mapContainerRef.current, {
      zoom: 12,
      center: [lng, lat],
      viewMode: "2D",
    });

    map.addControl(new AMap.Scale());

    if (item?.longitude && item?.latitude) {
      const marker = new AMap.Marker({ position: [lng, lat], draggable: true });
      map.add(marker);
      markerRef.current = marker;
    }

    map.on("click", (e: any) => {
      const { lng: clickLng, lat: clickLat } = e.lnglat;
      form.setValue("longitude", clickLng.toFixed(6), { shouldValidate: true });
      form.setValue("latitude", clickLat.toFixed(6), { shouldValidate: true });

      if (markerRef.current) {
        markerRef.current.setPosition([clickLng, clickLat]);
      } else {
        const marker = new AMap.Marker({ position: [clickLng, clickLat], draggable: true });
        map.add(marker);
        markerRef.current = marker;
      }
    });

    mapInstanceRef.current = map;

    AMap.plugin(["AMap.PlaceSearch"], () => {
      placeSearchRef.current = new AMap.PlaceSearch({
        pageSize: 5,
        pageIndex: 1,
        city: "全国",
      });
    });
  }, [form, item]);

  useEffect(() => {
    if (isOpen) {
      setTimeout(() => initMap(), 100);
    }
    return () => {
      if (mapInstanceRef.current) {
        mapInstanceRef.current.destroy();
        mapInstanceRef.current = null;
        markerRef.current = null;
        placeSearchRef.current = null;
      }
    };
  }, [isOpen, initMap]);

  useEffect(() => {
    if (item) {
      form.reset({
        name: item.name ?? "",
        longitude: item.longitude ?? "",
        latitude: item.latitude ?? "",
      });
      setSelectedColor(item.marker_color ?? "");
      setCategories(item.categories || []);
      setExistingPhotos((item.photos || []).map((p) => p.src));
      setPendingFiles([]);
    } else {
      form.reset({ name: "", longitude: "", latitude: "" });
      setSelectedColor("");
      setCategories([]);
      setExistingPhotos([]);
      setPendingFiles([]);
    }
  }, [item, form]);

  const handleSearch = useCallback(() => {
    if (!placeSearchRef.current || !mapSearchKeyword.trim()) return;
    const AMap = (window as any).AMap;
    placeSearchRef.current.search(mapSearchKeyword.trim(), (status: string, result: any) => {
      if (status !== "complete" || !result.poiList?.pois?.length) return;
      const poi = result.poiList.pois[0];
      const { lng, lat } = poi.location;
      form.setValue("longitude", lng.toFixed(6), { shouldValidate: true });
      form.setValue("latitude", lat.toFixed(6), { shouldValidate: true });
      if (mapInstanceRef.current) {
        mapInstanceRef.current.setZoomAndCenter(14, [lng, lat]);
        if (markerRef.current) {
          markerRef.current.setPosition([lng, lat]);
        } else {
          const marker = new AMap.Marker({ position: [lng, lat], draggable: true });
          mapInstanceRef.current.add(marker);
          markerRef.current = marker;
        }
      }
      if (!form.getValues("name")) form.setValue("name", poi.name);
    });
  }, [mapSearchKeyword, form]);

  const handleFilesChange = useCallback((files: FileWithPreview[]) => {
    setPendingFiles(files);
  }, []);

  const uploadAllFiles = async (): Promise<string[]> => {
    const urls: string[] = [];
    for (const fileItem of pendingFiles) {
      if (!(fileItem.file instanceof File)) continue;
      try {
        const presignResp = await ossPresignApi({ filename: fileItem.file.name });
        const uploadResp = await fetch(presignResp.url, {
          headers: { "Content-Type": "text/plain;charset=utf8" },
          method: "PUT",
          body: fileItem.file,
        });
        if (!uploadResp.ok) throw new Error(`Upload failed, status: ${uploadResp.status}`);
        urls.push(presignResp.cdn_url);
      } catch (error) {
        console.error("Upload failed:", error);
        throw error;
      }
    }
    return urls;
  };

  const handleSubmit = async (values: z.infer<typeof formSchema>) => {
    setLoading(true);
    setUploading(true);
    try {
      let newUrls: string[] = [];
      if (pendingFiles.length > 0) {
        newUrls = await uploadAllFiles();
      }
      const allPhotoUrls = [...existingPhotos, ...newUrls];

      await onSubmit({
        name: values.name,
        description: descRef.current?.value || undefined,
        longitude: values.longitude,
        latitude: values.latitude,
        date: dateRef.current?.value || undefined,
        marker_color: selectedColor,
        categories,
        url: urlRef.current?.value || undefined,
        url_label: urlLabelRef.current?.value || undefined,
        photo_urls: allPhotoUrls,
      });
      onClose();
      form.reset();
      setPendingFiles([]);
      setExistingPhotos([]);
      setCategories([]);
      setSelectedColor("");
    } catch (error) {
      console.error("Submit failed:", error);
    } finally {
      setLoading(false);
      setUploading(false);
    }
  };

  const handleCancel = () => {
    onClose();
    form.reset();
    setPendingFiles([]);
  };

  return (
    <Dialog open={isOpen} onOpenChange={(open) => { if (!open) onClose(); }}>
      <DialogContent className="sm:max-w-[650px] max-h-[90vh] flex flex-col overflow-hidden" onInteractOutside={(e) => e.preventDefault()}>
        <DialogHeader>
          <DialogTitle>{isEdit ? "编辑足迹" : "新增足迹"}</DialogTitle>
        </DialogHeader>
        <form
          className="w-full px-1 overflow-y-auto"
          method="post"
          autoComplete="off"
          onSubmit={form.handleSubmit(handleSubmit)}
        >
          <FieldGroup>
            <Controller
              name="name"
              control={form.control}
              render={({ field, fieldState }) => (
                <Field orientation="vertical" data-invalid={fieldState.invalid}>
                  <FieldLabel htmlFor={field.name}>地点名称</FieldLabel>
                  <FieldContent>
                    <Input {...field} id={field.name} placeholder="请输入地点名称" />
                    {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                  </FieldContent>
                </Field>
              )}
            />

            <Field orientation="vertical">
              <FieldLabel>坐标选择</FieldLabel>
              <FieldContent>
                <div className="flex gap-2 mb-2">
                  <Input
                    value={mapSearchKeyword}
                    onChange={(e) => setMapSearchKeyword(e.target.value)}
                    placeholder="搜索地点"
                    className="flex-1"
                    onKeyDown={(e) => {
                      if (e.key === "Enter") {
                        e.preventDefault();
                        handleSearch();
                      }
                    }}
                  />
                  <Button type="button" size="sm" onClick={handleSearch}>
                    搜索
                  </Button>
                </div>
                <div ref={mapContainerRef} className="w-full h-[250px] rounded-lg border" />
              </FieldContent>
            </Field>

            <div className="grid grid-cols-2 gap-4">
              <Controller
                name="longitude"
                control={form.control}
                render={({ field, fieldState }) => (
                  <Field orientation="vertical" data-invalid={fieldState.invalid}>
                    <FieldLabel htmlFor="longitude">经度</FieldLabel>
                    <FieldContent>
                      <Input {...field} id="longitude" placeholder="经度" />
                      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                    </FieldContent>
                  </Field>
                )}
              />
              <Controller
                name="latitude"
                control={form.control}
                render={({ field, fieldState }) => (
                  <Field orientation="vertical" data-invalid={fieldState.invalid}>
                    <FieldLabel htmlFor="latitude">纬度</FieldLabel>
                    <FieldContent>
                      <Input {...field} id="latitude" placeholder="纬度" />
                      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                    </FieldContent>
                  </Field>
                )}
              />
            </div>

            <Field orientation="vertical">
              <FieldLabel>描述</FieldLabel>
              <FieldContent>
                <Textarea
                  ref={descRef}
                  rows={2}
                  placeholder="描述（可选）"
                  defaultValue={item?.description ?? ""}
                />
              </FieldContent>
            </Field>

            <div className="grid grid-cols-2 gap-4">
              <Field orientation="vertical">
                <FieldLabel>日期</FieldLabel>
                <FieldContent>
                  <Input
                    ref={dateRef}
                    placeholder="如 2024-05-01 或 2024, 2025"
                    defaultValue={item?.date ?? ""}
                  />
                </FieldContent>
              </Field>
              <Field orientation="vertical">
                <FieldLabel>标记颜色</FieldLabel>
                <FieldContent>
                  <div className="flex gap-2 flex-wrap">
                    {MARKER_COLORS.map((c) => (
                      <button
                        key={c.value}
                        type="button"
                        title={c.label}
                        className={`w-7 h-7 rounded-full border-2 transition-all ${
                          selectedColor === c.value
                            ? "border-gray-800 scale-110"
                            : "border-gray-300"
                        }`}
                        style={{
                          background: c.value === ""
                            ? "linear-gradient(135deg, #ff6b6b, #4ecdc4, #a55eea)"
                            : GRADIENT_MAP[c.value],
                        }}
                        onClick={() => setSelectedColor(c.value)}
                      />
                    ))}
                  </div>
                </FieldContent>
              </Field>
            </div>

            <Field orientation="vertical">
              <FieldLabel>分类标签</FieldLabel>
              <FieldContent>
                <InputTags
                  value={categories}
                  onChange={setCategories}
                  placeholder="输入标签，回车/逗号确认"
                />
              </FieldContent>
            </Field>

            <div className="grid grid-cols-2 gap-4">
              <Field orientation="vertical">
                <FieldLabel>关联链接</FieldLabel>
                <FieldContent>
                  <Input
                    ref={urlRef}
                    placeholder="https://..."
                    defaultValue={item?.url ?? ""}
                  />
                </FieldContent>
              </Field>
              <Field orientation="vertical">
                <FieldLabel>链接文案</FieldLabel>
                <FieldContent>
                  <Input
                    ref={urlLabelRef}
                    placeholder="如：阅读游记"
                    defaultValue={item?.url_label ?? ""}
                  />
                </FieldContent>
              </Field>
            </div>

            <Field orientation="vertical">
              <FieldLabel>照片</FieldLabel>
              <FieldContent>
                {existingPhotos.length > 0 && (
                  <div className="flex gap-2 mb-2 flex-wrap">
                    {existingPhotos.map((url, i) => (
                      <div key={i} className="relative group">
                        <img
                          src={url + "!photothumb"}
                          alt=""
                          className="w-16 h-12 object-cover rounded border"
                        />
                        <button
                          type="button"
                          onClick={() =>
                            setExistingPhotos(existingPhotos.filter((_, idx) => idx !== i))
                          }
                          className="absolute -top-1 -right-1 w-4 h-4 bg-red-500 text-white rounded-full text-xs opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center"
                        >
                          ×
                        </button>
                      </div>
                    ))}
                  </div>
                )}
                <FileUploadCompact
                  maxFiles={20}
                  maxSize={10 * 1024 * 1024}
                  accept="image/*"
                  multiple={true}
                  onFilesChange={handleFilesChange}
                />
                {uploading && (
                  <p className="text-sm text-gray-500 mt-2">
                    正在上传 {pendingFiles.length} 张照片...
                  </p>
                )}
              </FieldContent>
            </Field>

            <Field orientation="horizontal">
              <Button type="submit" size="sm" loading={loading || uploading}>
                {isEdit ? "修改" : "提交"}
              </Button>
              <Button
                size="sm"
                variant="secondary"
                onClick={(e) => {
                  e.preventDefault();
                  handleCancel();
                }}
                disabled={loading || uploading}
              >
                取消
              </Button>
            </Field>
          </FieldGroup>
        </form>
      </DialogContent>
    </Dialog>
  );
}

export default AdminFootprintDialog;
