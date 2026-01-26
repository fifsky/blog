import Taro from "@tarojs/taro";

export async function putFileToPresignUrl(presignUrl: string, filePath: string) {
  const fs = Taro.getFileSystemManager();
  const buf = await new Promise<ArrayBuffer>((resolve, reject) => {
    fs.readFile({
      filePath,
      success(res) {
        resolve(res.data as ArrayBuffer);
      },
      fail(err) {
        reject(err);
      },
    });
  });

  const resp = await Taro.request({
    url: presignUrl,
    method: "PUT",
    data: buf,
    header: {
      "Content-Type": "text/plain;charset=utf8",
    },
  });

  if (resp.statusCode >= 200 && resp.statusCode < 300) {
    return;
  }
  throw new Error(`上传文件失败(${resp.statusCode})`);
}
