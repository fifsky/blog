import { createRequire } from "module";
import { existsSync, readFileSync, writeFileSync } from "fs";
import path from "path";
import { execFileSync } from "child_process";

const require = createRequire(import.meta.url);

function resolvePkgDir(pkgName) {
  const pkgJsonPath = require.resolve(`${pkgName}/package.json`, { paths: [process.cwd()] });
  return path.dirname(pkgJsonPath);
}

function patchBindingSource(bindingDir) {
  const libPath = path.join(bindingDir, "src", "lib.rs");
  if (!existsSync(libPath)) return;
  const src = readFileSync(libPath, "utf8");
  let next = src;
  next = next.replace(/^\s*conf\.build_es5,\s*\n/gm, "");
  next = next.replace(/^\s*conf\.sub_pkg,\s*\n/gm, "");
  next = next.replace(/^\s*conf\.page_dir,\s*\n/gm, "");
  if (next !== src) {
    writeFileSync(libPath, next);
  }
}

function ensureDarwinX64(bindingDir) {
  const nodeFile = path.join(bindingDir, "taro.darwin-x64.node");
  if (existsSync(nodeFile)) return;

  patchBindingSource(bindingDir);

  const postinstall = path.join(bindingDir, "postinstall.js");
  if (!existsSync(postinstall)) return;

  execFileSync(process.execPath, [postinstall], {
    cwd: bindingDir,
    stdio: "inherit",
    env: {
      ...process.env,
      BUILD_TARO_FROM_SOURCE: "1",
    },
  });

  if (!existsSync(nodeFile)) {
    throw new Error("taro binding build succeeded but taro.darwin-x64.node not found");
  }
}

const bindingDir = (() => {
  try {
    return resolvePkgDir("@tarojs/binding");
  } catch {
    return "";
  }
})();

if (!bindingDir) {
  process.exitCode = 0;
} else if (process.platform === "darwin" && process.arch === "x64") {
  try {
    ensureDarwinX64(bindingDir);
  } catch (e) {
    console.error(e);
    process.exitCode = 1;
  }
}
