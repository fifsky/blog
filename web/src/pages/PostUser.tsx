import React, { useEffect, useState } from "react";
import { userGetApi, userCreateApi, userUpdateApi } from "@/service";
import { useLocation, useNavigate, Link } from "react-router-dom";

export default function PostUser() {
  const [user, setUser] = useState<any>({ type: 1 });
  const location = useLocation();
  const navigate = useNavigate();
  const params = new URLSearchParams(location.search);

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (user.password1 !== user.password2) {
      alert("两次输入的密码不一致");
      return;
    }
    const { id, name, nick_name, password1, email, type } = user;
    const data = {
      id,
      name,
      nick_name,
      password: password1,
      email,
      type: parseInt(type),
    };
    if (id) await userUpdateApi(data);
    else await userCreateApi(data);
    navigate("/admin/users");
  };

  useEffect(() => {
    (async () => {
      if (params.get("id")) {
        const u = await userGetApi({ id: parseInt(params.get("id")!) });
        setUser({ ...u, nick_name: u.nick_name });
      }
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div>
      <h2>
        {user.id ? "编辑" : "新增"}用户
        <Link to="/admin/users">
          <i className="iconfont icon-undo" style={{ color: "#444" }}></i>
          返回列表
        </Link>
      </h2>
      <div className="message"></div>
      <form className="vf" method="post" autoComplete="off" onSubmit={submit}>
        <p>
          <label className="label_input">
            用户名 <span className="desc">(必填)</span>
          </label>
          <input
            className="input_text"
            size={50}
            value={user.name || ""}
            onChange={(e) => setUser({ ...user, name: e.target.value })}
          ></input>
        </p>
        <span className="hint">
          此用户名将作为用户登录时所用的名称，请不要与系统中现有的用户名重复。
        </span>
        <p></p>
        <p>
          <label className="label_input">
            邮箱 <span className="desc">(必填)</span>
          </label>
          <input
            className="input_text"
            size={50}
            value={user.email || ""}
            onChange={(e) => setUser({ ...user, email: e.target.value })}
          ></input>
        </p>
        <span className="hint">
          电子邮箱地址将作为此用户的主要联系方式，请不要与系统中现有的电子邮箱地址重复。
        </span>
        <p></p>
        <p>
          <label className="label_input">昵称</label>
          <input
            className="input_text"
            size={50}
            value={user.nick_name || ""}
            onChange={(e) => setUser({ ...user, nick_name: e.target.value })}
          ></input>
        </p>
        <span className="hint">
          用户昵称可以与用户名不同,
          用于前台显示，如果你将此项留空，将默认使用用户名。
        </span>
        <p>
          <label className="label_input">
            密码 <span className="desc">(必填)</span>
          </label>
          <input
            type="password"
            className="input_text"
            size={50}
            value={user.password1 || ""}
            onChange={(e) => setUser({ ...user, password1: e.target.value })}
          />
        </p>
        <span className="hint">为用户分配一个密码。</span>
        <p>
          <label className="label_input">
            确认密码 <span className="desc">(必填)</span>
          </label>
          <input
            type="password"
            className="input_text"
            size={50}
            value={user.password2 || ""}
            onChange={(e) => setUser({ ...user, password2: e.target.value })}
          />
          <span className="hint">
            请确认你的密码，与上面输入的密码保持一致。
          </span>
        </p>
        <p>
          <label className="label_input">
            角色 <span className="desc">(必填)</span>
          </label>
          <select
            name="type"
            value={user.type}
            onChange={(e) => setUser({ ...user, type: e.target.value })}
          >
            <option value="1">管理员</option>
            <option value="2">编辑</option>
          </select>
          <span className="hint">
            管理员具有所有操作权限，编辑仅能包含文章、评论、心情的操作权限。
          </span>
        </p>
        <p className="act">
          <input className="formbutton" type="submit" value="保存" />
        </p>
      </form>
    </div>
  );
}
