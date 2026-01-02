import React, { useState } from "react";
import { CHeader } from "@/components/CHeader";
import { CFooter } from "@/components/CFooter";
import { useStore } from "@/store/context";
import { useNavigate } from "react-router";
import { LoginRequest } from "@/types/openapi";

export default function Login() {
  const [formdata, setFormdata] = useState<LoginRequest>({
    password: "",
    user_name: "",
  });
  const { loginAction } = useStore();
  const navigate = useNavigate();
  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    await loginAction(formdata);
    navigate("/admin/index");
  };
  return (
    <div id="container">
      <CHeader />
      <div className="admin">
        <div id="content">
          <div id="sign-in">
            <h2>博客管理登录</h2>
            <div className="message"></div>
            <form method="post" onSubmit={submit} className="vf lf">
              <p>
                <label className="label_input">用户名：</label>
                <input
                  type="text"
                  className="input_text"
                  value={formdata.user_name || ""}
                  onChange={(e) =>
                    setFormdata((prev) => ({
                      ...prev,
                      user_name: e.target.value,
                    }))
                  }
                />
              </p>
              <p>
                <label className="label_input">密码：</label>
                <input
                  type="password"
                  className="input_text"
                  value={formdata.password || ""}
                  onChange={(e) =>
                    setFormdata((prev) => ({
                      ...prev,
                      password: e.target.value,
                    }))
                  }
                />
              </p>
              <p className="act">
                <input type="submit" className="formbutton" value="登录" />
              </p>
            </form>
          </div>
        </div>
      </div>
      <CFooter />
    </div>
  );
}
