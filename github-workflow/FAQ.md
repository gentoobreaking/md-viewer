FAQ.1

! [remote rejected] main -> main (refusing to allow an OAuth App to create or update workflow .github/workflows/build.yml without workflow scope)

ANSWER:

這是 GitHub 的權限問題，你的 Git credential 沒有 `workflow` scope。

**最快的解法：重新產生一個有 workflow 權限的 Personal Access Token（PAT）**

1. 去 **GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic)**
2. 點 **Generate new token (classic)**
3. 勾選 `repo` 和 **`workflow`** 這兩個 scope
4. 產生後複製 token

然後更新你本機的 credential：

```bash
git remote set-url origin https://<你的token>@github.com/<username>/<repo>.git
```

再重新 push 就可以了。

---

或者如果你用的是 **GitHub Desktop**，它本身不支援 workflow scope，建議改用 command line 的方式 push `.github/workflows/` 的檔案。
