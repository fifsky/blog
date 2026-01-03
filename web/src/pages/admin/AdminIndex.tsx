import React, { useEffect, useState } from "react";
import { settingApi, settingUpdateApi } from "@/service";

export default function AdminIndex() {
  const [formdata, setFormdata] = useState<Record<string, string>>({});
  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    await settingUpdateApi({ kv: formdata });
    alert("保存成功");
  };
  useEffect(() => {
    (async () => {
      const data = await settingApi();
      setFormdata(data.kv || {});
    })();
  }, []);
  return (
    <div id="settings">
      <h2>站点设置</h2>
      <div className="message">保存成功</div>
      <form className="nf" method="post" autoComplete="off" onSubmit={submit}>
        <p>
          <label className="label_input">站点名称</label>
          <input
            type="text"
            className="input_text"
            size={50}
            name="site_name"
            value={formdata.site_name || ""}
            onChange={(e) =>
              setFormdata((prev) => ({ ...prev, site_name: e.target.value }))
            }
          />
          <span className="hint">站点的名称将显示在网页的标题处。</span>
        </p>
        <p>
          <label className="label_input">站点描述</label>
          <textarea
            name="site_desc"
            rows={3}
            cols={50}
            value={formdata.site_desc || ""}
            onChange={(e) =>
              setFormdata((prev) => ({ ...prev, site_desc: e.target.value }))
            }
          ></textarea>
          <span className="hint">站点描述将显示在网页代码的头部。</span>
        </p>
        <p>
          <label className="label_input">关键字</label>
          <input
            type="text"
            className="input_text"
            size={50}
            name="site_keyword"
            value={formdata.site_keyword || ""}
            onChange={(e) =>
              setFormdata((prev) => ({
                ...prev,
                site_keyword: e.target.value,
              }))
            }
          />
          <span className="hint">请以半角逗号","分割多个关键字。</span>
        </p>
        <p>
          <label className="label_input">每页显示文章数</label>
          <input
            className="input_text"
            style={{ width: 50 }}
            name="post_num"
            type="text"
            value={formdata.post_num || ""}
            onChange={(e) =>
              setFormdata((prev) => ({ ...prev, post_num: e.target.value }))
            }
          />
        </p>
        <p className="act">
          <input className="formbutton" type="submit" value="保存" />
        </p>
      </form>
    </div>
  );
}
