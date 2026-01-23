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
import { useEffect, useState, useCallback } from "react";
import { regionListApi, ossPresignApi } from "@/service";
import FileUploadCompact from "./file-upload/compact-upload";
import { FileWithPreview } from "@/hooks/use-file-upload";

interface AdminPhotoDialogProps {
  isOpen: boolean;
  onClose: () => void;
  item: PhotoItem | undefined;
  onSubmit: (values: FormValues) => Promise<void>;
}

interface FormValues {
  title: string;
  description?: string;
  srcs: string[];
  province: number;
  city: number;
}

const formSchema = z.object({
  title: z.string().min(1, "请输入照片标题"),
  description: z.string().optional(),
  province: z.string().min(1, "请选择省份"),
  city: z.string().min(1, "请选择城市"),
});

export function AdminPhotoDialog({ isOpen, onClose, item, onSubmit }: AdminPhotoDialogProps) {
  const [loading, setLoading] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [provinces, setProvinces] = useState<RegionItem[]>([]);
  const [cities, setCities] = useState<RegionItem[]>([]);
  const [pendingFiles, setPendingFiles] = useState<FileWithPreview[]>([]);

  const isEdit = !!item?.id;

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      title: item?.title ?? "",
      description: item?.description ?? "",
      province: item?.province ? String(item.province) : "",
      city: item?.city ? String(item.city) : "",
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
    const provinceId = Number(selectedProvince);
    if (provinceId > 0) {
      const loadCities = async () => {
        const resp = await regionListApi({ parent_id: provinceId });
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
        province: item?.province ? String(item.province) : "",
        city: item?.city ? String(item.city) : "",
      });
      setPendingFiles([]);

      // Load cities for existing item
      if (item.province > 0) {
        regionListApi({ parent_id: item.province }).then((resp) => {
          setCities(resp.list || []);
        });
      }
    } else {
      form.reset({
        title: "",
        description: "",
        province: "",
        city: "",
      });
      setPendingFiles([]);
    }
  }, [item, form]);

  // Handle files change from upload component
  const handleFilesChange = useCallback((files: FileWithPreview[]) => {
    setPendingFiles(files);
  }, []);

  // Upload all pending files and return their CDN URLs
  const uploadAllFiles = async (): Promise<string[]> => {
    const urls: string[] = [];

    for (const fileItem of pendingFiles) {
      // Only upload actual File objects, skip FileMetadata
      if (!(fileItem.file instanceof File)) {
        continue;
      }

      try {
        // Get presigned URL from backend
        const presignResp = await ossPresignApi({ filename: fileItem.file.name });

        // Upload file directly using PUT request
        const uploadResp = await fetch(presignResp.url, {
          headers: {
            "Content-Type": "text/plain;charset=utf8",
          },
          method: "PUT",
          body: fileItem.file,
        });

        if (!uploadResp.ok) {
          throw new Error(`Upload failed, status: ${uploadResp.status}`);
        }

        urls.push(presignResp.cdn_url);
        console.log("File uploaded successfully to:", presignResp.cdn_url);
      } catch (error) {
        console.error("Upload failed:", error);
        throw error;
      }
    }

    return urls;
  };

  const handleSubmit = async (values: z.infer<typeof formSchema>) => {
    if (!isEdit && pendingFiles.length === 0) {
      alert("请先选择照片");
      return;
    }

    setLoading(true);
    setUploading(true);

    try {
      let srcs: string[] = [];

      if (!isEdit) {
        // Upload all files first
        srcs = await uploadAllFiles();
        if (srcs.length === 0) {
          alert("上传失败，请重试");
          return;
        }
      }

      await onSubmit({
        title: values.title,
        description: values.description,
        srcs,
        province: Number(values.province),
        city: Number(values.city),
      });

      onClose();
      form.reset();
      setPendingFiles([]);
    } catch (error) {
      console.error("Submit failed:", error);
      alert("提交失败，请重试");
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
                        onValueChange={(v) => field.onChange(v)}
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
                  <FileUploadCompact
                    maxFiles={10}
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
            )}

            <Field orientation="horizontal">
              <Button type="submit" size="sm" loading={loading || uploading}>
                {isEdit
                  ? "修改"
                  : `提交${pendingFiles.length > 0 ? ` (${pendingFiles.length} 张)` : ""}`}
              </Button>
              <Button
                size={"sm"}
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

export default AdminPhotoDialog;
