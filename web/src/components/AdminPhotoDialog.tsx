import { Controller, useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "./ui/select";
import { Field, FieldLabel, FieldError, FieldGroup, FieldContent } from "./ui/field";
import { Input } from "./ui/input";
import { Textarea } from "./ui/textarea";
import { Button } from "./ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "./ui/dialog";
import { PhotoItem, RegionItem } from "@/types/openapi";
import { useEffect, useState } from "react";
import { regionListApi, ossPresignApi } from "@/service";

interface AdminPhotoDialogProps {
  isOpen: boolean;
  onClose: () => void;
  item: PhotoItem | undefined;
  onSubmit: (values: z.infer<typeof formSchema>) => Promise<void>;
}

const formSchema = z.object({
  title: z.string().min(1, "请输入照片标题"),
  description: z.string().optional(),
  src: z.string().optional(),
  province: z.string().min(1, "请选择省份"),
  city: z.string().min(1, "请选择城市"),
});

export function AdminPhotoDialog({ isOpen, onClose, item, onSubmit }: AdminPhotoDialogProps) {
  const [loading, setLoading] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [provinces, setProvinces] = useState<RegionItem[]>([]);
  const [cities, setCities] = useState<RegionItem[]>([]);
  const [previewUrl, setPreviewUrl] = useState<string>("");

  const isEdit = !!item?.id;

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      title: item?.title ?? "",
      description: item?.description ?? "",
      src: item?.src ?? "",
      province: item?.province ?? "",
      city: item?.city ?? "",
    },
    mode: "onChange",
  });

  // Load provinces on mount
  useEffect(() => {
    const loadProvinces = async () => {
      const resp = await regionListApi({ parent_id: 0 });
      setProvinces(resp.list || []);
    };
    loadProvinces();
  }, []);

  // Load cities when province changes
  const selectedProvince = form.watch("province");
  useEffect(() => {
    if (selectedProvince) {
      const loadCities = async () => {
        const resp = await regionListApi({ parent_id: Number(selectedProvince) });
        setCities(resp.list || []);
      };
      loadCities();
    } else {
      setCities([]);
    }
  }, [selectedProvince]);

  // 当 item 属性变化时，更新表单值
  useEffect(() => {
    if (item) {
      form.reset({
        title: item?.title ?? "",
        description: item?.description ?? "",
        src: item?.src ?? "",
        province: item?.province ?? "",
        city: item?.city ?? "",
      });
      setPreviewUrl(item?.thumbnail || "");

      // Load cities for existing item
      if (item.province) {
        regionListApi({ parent_id: Number(item.province) }).then((resp) => {
          setCities(resp.list || []);
        });
      }
    } else {
      form.reset({
        title: "",
        description: "",
        src: "",
        province: "",
        city: "",
      });
      setPreviewUrl("");
    }
  }, [item, form]);

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploading(true);
    try {
      // Get presigned URL from backend
      const presignResp = await ossPresignApi({ filename: file.name });

      // Upload file directly using PUT request
      const uploadResp = await fetch(presignResp.url, {
        headers: {
          "Content-Type": "text/plain;charset=utf8",
        },
        method: "PUT",
        body: file,
      });

      if (!uploadResp.ok) {
        throw new Error(`Upload failed, status: ${uploadResp.status}`);
      }

      // Set the CDN URL as src
      form.setValue("src", presignResp.cdn_url);
      setPreviewUrl(presignResp.cdn_url + "!photothumb");

      console.log("File uploaded successfully to:", presignResp.cdn_url);
    } catch (error) {
      console.error("Upload failed:", error);
      alert("上传失败，请重试");
    } finally {
      setUploading(false);
    }
  };

  const handleSubmit = async (values: z.infer<typeof formSchema>) => {
    if (!isEdit && !values.src) {
      alert("请先上传照片");
      return;
    }

    setLoading(true);
    try {
      await onSubmit(values);
      onClose();
      form.reset();
      setPreviewUrl("");
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    onClose();
    form.reset();
    setPreviewUrl("");
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>{isEdit ? "编辑照片" : "新增照片"}</DialogTitle>
        </DialogHeader>
        <form
          className="w-full px-1"
          method="post"
          autoComplete="off"
          onSubmit={form.handleSubmit(handleSubmit)}
        >
          <FieldGroup>
            <Controller
              name="title"
              control={form.control}
              render={({ field, fieldState }) => (
                <Field orientation="vertical" data-invalid={fieldState.invalid}>
                  <FieldLabel htmlFor={field.name}>照片标题</FieldLabel>
                  <FieldContent>
                    <Input {...field} id={field.name} placeholder="请输入照片标题" />
                    {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                  </FieldContent>
                </Field>
              )}
            />

            <Controller
              name="description"
              control={form.control}
              render={({ field, fieldState }) => (
                <Field orientation="vertical" data-invalid={fieldState.invalid}>
                  <FieldLabel htmlFor={field.name}>照片描述</FieldLabel>
                  <FieldContent>
                    <Textarea {...field} id={field.name} rows={3} placeholder="请输入照片描述" />
                    {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                  </FieldContent>
                </Field>
              )}
            />

            <div className="grid grid-cols-2 gap-4">
              <Controller
                name="province"
                control={form.control}
                render={({ field, fieldState }) => (
                  <Field orientation="vertical" data-invalid={fieldState.invalid}>
                    <FieldLabel htmlFor="province">省份</FieldLabel>
                    <FieldContent>
                      <Select
                        value={field.value}
                        onValueChange={(v) => {
                          field.onChange(v);
                          form.setValue("city", "");
                        }}
                      >
                        <SelectTrigger size="sm">
                          <SelectValue placeholder="请选择省份" />
                        </SelectTrigger>
                        <SelectContent>
                          {provinces.map((p) => (
                            <SelectItem key={p.region_id} value={String(p.region_id)}>
                              {p.region_name}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                    </FieldContent>
                  </Field>
                )}
              />

              <Controller
                name="city"
                control={form.control}
                render={({ field, fieldState }) => (
                  <Field orientation="vertical" data-invalid={fieldState.invalid}>
                    <FieldLabel htmlFor="city">城市</FieldLabel>
                    <FieldContent>
                      <Select
                        value={field.value}
                        onValueChange={field.onChange}
                        disabled={!selectedProvince}
                      >
                        <SelectTrigger size="sm">
                          <SelectValue placeholder="请选择城市" />
                        </SelectTrigger>
                        <SelectContent>
                          {cities.map((c) => (
                            <SelectItem key={c.region_id} value={String(c.region_id)}>
                              {c.region_name}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                    </FieldContent>
                  </Field>
                )}
              />
            </div>

            {!isEdit && (
              <Field orientation="vertical">
                <FieldLabel>上传照片</FieldLabel>
                <FieldContent>
                  <div className="flex items-center gap-4">
                    <input
                      type="file"
                      accept="image/*"
                      onChange={handleFileUpload}
                      disabled={uploading}
                      className="text-sm"
                    />
                    {uploading && <span className="text-sm text-gray-500">上传中...</span>}
                  </div>
                  {previewUrl && (
                    <div className="mt-2">
                      <img
                        src={previewUrl}
                        alt="预览"
                        className="max-w-[200px] max-h-[150px] object-cover rounded"
                      />
                    </div>
                  )}
                </FieldContent>
              </Field>
            )}

            {isEdit && previewUrl && (
              <Field orientation="vertical">
                <FieldLabel>当前照片</FieldLabel>
                <FieldContent>
                  <img
                    src={previewUrl}
                    alt="当前照片"
                    className="max-w-[200px] max-h-[150px] object-cover rounded"
                  />
                  <p className="text-xs text-gray-500 mt-1">编辑模式不支持重新上传</p>
                </FieldContent>
              </Field>
            )}

            <Field orientation="horizontal">
              <Button type="submit" size="sm" loading={loading || uploading}>
                {isEdit ? "修改" : "添加"}
              </Button>
              <Button
                size={"sm"}
                variant="link"
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

export default AdminPhotoDialog;
