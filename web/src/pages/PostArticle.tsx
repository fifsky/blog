import React, { useEffect, useState } from "react";
import {
  articleDetailApi,
  articleCreateApi,
  articleUpdateApi,
  cateListApi,
} from "@/service";
import { useLocation, useNavigate, Link } from "react-router-dom";
import "@wangeditor/editor/dist/css/style.css";
import { Editor, Toolbar } from "@wangeditor/editor-for-react";
import type {
  IDomEditor,
  IEditorConfig,
  IToolbarConfig,
} from "@wangeditor/editor";
import { getApiUrl, getAccessToken } from "@/utils/common";

export default function PostArticle() {
  const [article, setArticle] = useState<any>({ type: 1 });
  const [cates, setCates] = useState<any[]>([]);
  const [editor, setEditor] = useState<IDomEditor | null>(null);

  const location = useLocation();
  const navigate = useNavigate();
  const params = new URLSearchParams(location.search);

  const toolbarConfig: Partial<IToolbarConfig> = {
    excludeKeys: ["uploadVideo", "fontFamily", "lineHeight", "group-indent"],
  };

  const editorConfig: Partial<IEditorConfig> = {
    placeholder: "请输入内容...",
    MENU_CONF: {
      uploadImage: {
        server: getApiUrl("/api/admin/upload"),
        fieldName: "uploadFile",
        headers: { "Access-Token": getAccessToken() },
        withCredentials: false,
        maxFileSize: 10 * 1024 * 1024,
        allowedFileTypes: ["image/*"],
        onFailed(file: File, res: any) {
          alert(`${file.name} 上传失败` + (res.message || ""));
        },
        onError(file: File, err: any) {
          alert(`${file.name} 上传失败` + (err || ""));
        },
      },
    },
  };

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    const { id, cate_id, title, content, type, url } = article;
    if (!title || title.length === 0) {
      alert("标题不能为空");
      return;
    }
    if (!cate_id || cate_id < 1) {
      alert("请选择分类");
      return;
    }
    const data = { id, cate_id, title, content, type, url };
    if (id)
      await articleUpdateApi({
        id,
        cate_id,
        title,
        content,
        type,
        url,
      });
    else await articleCreateApi(data);
    navigate("/admin/articles");
  };

  useEffect(() => {
    (async () => {
      if (params.get("id")) {
        const a = await articleDetailApi({ id: parseInt(params.get("id")!) });
        setArticle(a);
      }
      const ret = await cateListApi({});
      setCates(ret.list || []);
      setArticle((prev: any) => ({
        ...prev,
        cate_id: prev.cate_id || ret.list?.[0]?.id || 0,
      }));
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);
  useEffect(() => {
    return () => {
      if (editor) editor.destroy();
    };
  }, [editor]);

  return (
    <div id="articles">
      <h2>
        {article.id ? "编辑" : "撰写"}文章
        <Link to="/admin/articles">
          <i className="iconfont icon-undo" style={{ color: "#444" }}></i>
          返回列表
        </Link>
      </h2>
      <div className="message"></div>
      <form className="vf" method="post" autoComplete="off" onSubmit={submit}>
        <div className="clearfix">
          <div className="col-left">
            <p>
              <label className="label_input">标题</label>
              <input
                type="text"
                className="input_text"
                maxLength={200}
                size={50}
                name="title"
                value={article.title || ""}
                onChange={(e) =>
                  setArticle((prev: any) => ({
                    ...prev,
                    title: e.target.value,
                  }))
                }
              />
            </p>
            <p>
              <label className="label_input">分类</label>
              {cates.length > 0 && (
                <select
                  name="cate_id"
                  value={article.cate_id || ""}
                  onChange={(e) =>
                    setArticle((prev: any) => ({
                      ...prev,
                      cate_id: Number(e.target.value),
                    }))
                  }
                >
                  {cates.map((v) => (
                    <option key={v.id} value={v.id}>
                      {v.name}
                    </option>
                  ))}
                </select>
              )}
            </p>
            {article.type == 2 && (
              <p>
                <label className="label_input">缩略名</label>
                <input
                  type="text"
                  className="input_text"
                  maxLength={200}
                  size={50}
                  name="url"
                  value={article.url || ""}
                  onChange={(e) =>
                    setArticle((prev: any) => ({
                      ...prev,
                      url: e.target.value,
                    }))
                  }
                />
                <span className="hint">
                  页面的URL名称，如红色部分http://domain.com/
                  <span style={{ color: "red" }}>about</span>
                </span>
              </p>
            )}
            <input type="hidden" name="id" value={article.id || ""} />
          </div>
          <div className="col-right">
            <p>
              <label className="label_input">类型</label>
              <input
                className="input_check"
                name="type"
                type="radio"
                value={1}
                checked={article.type === 1}
                onChange={() =>
                  setArticle((prev: any) => ({ ...prev, type: 1 }))
                }
              />
              文章
              <input
                className="input_check"
                name="type"
                type="radio"
                value={2}
                checked={article.type === 2}
                onChange={() =>
                  setArticle((prev: any) => ({ ...prev, type: 2 }))
                }
              />
              页面
            </p>
          </div>
        </div>
        <div id="editor">
          <div style={{ border: "1px solid #ddd" }}>
            <Toolbar
              editor={editor}
              defaultConfig={toolbarConfig}
              mode="default"
              style={{ borderBottom: "1px solid #ddd" }}
            />
            <Editor
              style={{ height: 500, overflowY: "hidden" }}
              defaultConfig={editorConfig}
              value={article.content || ""}
              onCreated={(ed: IDomEditor) => setEditor(ed)}
              onChange={(ed: IDomEditor) =>
                setArticle((prev: any) => ({ ...prev, content: ed.getHtml() }))
              }
              mode="default"
            />
            <input type="hidden" name="content" value={article.content || ""} />
          </div>
        </div>
        <p className="act">
          <input className="formbutton" type="submit" value="发布" />
          <a
            id="_save_draft"
            href="#"
            className="ml10"
            onClick={(e) => e.preventDefault()}
          >
            保存草稿
          </a>
        </p>
      </form>
    </div>
  );
}
