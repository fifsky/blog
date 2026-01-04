import React, { useEffect, useState } from "react";
import { settingApi, settingUpdateApi } from "@/service";

export default function AdminIndex() {
  const [formdata, setFormdata] = useState<Record<string, string>>({});
  const [showMessage, setShowMessage] = useState(false);
  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    await settingUpdateApi({ kv: formdata });
    setShowMessage(true);
    setTimeout(() => setShowMessage(false), 3000);
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
      {showMessage && <div className="message">保存成功</div>}
      <form className="nf" method="post" autoComplete="off" onSubmit={submit}>
        <p className="flex">
          <label className="label_input">站点名称</label>
          <div>
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
          </div>
        </p>
        <p className="flex">
          <label className="label_input">站点描述</label>
          <div>
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
          </div>
        </p>
        <p className="flex">
          <label className="label_input">关键字</label>
          <div>
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
          </div>
        </p>
        <p className="flex">
          <label className="label_input">每页显示文章数</label>
          <div>
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
          </div>
        </p>
        <p className="flex items-center">
          <div className="label_input"></div>
          <div>
            <button className="formbutton" type="submit">
              保存
            </button>
          </div>
        </p>
      </form>
    </div>
  );
}
