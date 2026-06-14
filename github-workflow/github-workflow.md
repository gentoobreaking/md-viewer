用 GitHub Actions 來自動執行這段 build script，建立一個 workflow 檔案。
這是一個 GitHub Actions workflow，用來在 macOS runner 上自動執行你的 build script：
檔案放到repo的 `.github/workflows/build.yml` 

**幾個重點說明：**

**Trigger 條件**
- push 或 PR 到 `main`/`master` 時自動觸發
- `workflow_dispatch` 讓你可以在 GitHub UI 手動觸發，並且可以帶入自訂的 `MARKETING_VERSION`

**環境**
- `runs-on: macos-latest` — Swift 和 `install_name_tool` 都是 macOS 內建工具，必須用 macOS runner
- Go 版本從你的 `go.mod` 自動讀取（`go-version-file: go.mod`）

**Build number 持久化**
- 用 `actions/cache` 把 `.build_number` 快取起來，讓每次 build 的號碼能夠累加，而不是每次都從 1 開始

**產出物 (Artifacts)**
- `.app` bundle 和獨立 binary 都會上傳，在 Actions run 頁面可以直接下載，保留 14 天

**如果你還需要：**
- **Release 自動發布** — 加上 `softprops/action-gh-release` 在打 tag 時自動建立 GitHub Release 並附上 zip
- **Code signing** — 需要把 Developer ID 憑證存到 Secrets，再加 `codesign` 步驟
- **Notarization** — 需要 Apple ID / App-specific password 存到 Secrets，用 `xcrun notarytool` 公證

-----
最簡單的方式是用 workflow_dispatch，讓你手動決定什麼時候跑，而不是每次 push 都自動觸發：

```
on:
  workflow_dispatch:    # 只有手動點按鈕才會跑
    inputs:
      marketing_version:
        description: 'Marketing version (e.g. 1.0.0)'
        required: false
        default: '0.2.2'
```

----

可以！但有個限制：GitHub Actions **無法直接讀取你 repo 裡 shell script 的變數**，因為那個變數只有在 script 執行時才存在。

不過有幾種方式可以達到「不用手動輸入版本號」的效果：

---

** 方案: 從 `.build_number` 或獨立的版本檔讀取 **

在 repo 裡建一個 `VERSION` 檔：
```
0.2.2
```

然後 `build.sh` 改成讀這個檔：
```bash
MARKETING_VERSION="${MARKETING_VERSION:-$(cat VERSION)}"
```

workflow 也改成讀這個檔：
```yaml
- name: Read version
  run: echo "MARKETING_VERSION=$(cat VERSION)" >> $GITHUB_ENV

- name: Build
  run: ./build.sh
  # 不用再傳入 env，build.sh 會自己讀 VERSION 檔
```

這樣只需要改 `VERSION` 這一個檔案，build.sh 和 workflow 都會自動跟著用，也不需要 `workflow_dispatch` 的 input。

---

GitHub → 你的 repo → 上方選單點 Actions → 左側選 Build md-viewer → 右側會出現 Run workflow 按鈕。

---

先在 GitHub 建立一個有 `repo` 權限的新 PAT，然後存到 **md-viewer repo** 的 Secrets：

- 名稱：`TAP_GITHUB_TOKEN`
- 內容：新的 PAT

---

然後在 `build.yml` 的最後加上這個步驟，放在 `Create GitHub Release` 之後：

```yaml
      - name: Update Homebrew tap
        if: ${{ github.event.inputs.create_release == 'true' }}
        env:
          TAP_GITHUB_TOKEN: ${{ secrets.TAP_GITHUB_TOKEN }}
          VERSION: ${{ env.MARKETING_VERSION }}
          GITHUB_USERNAME: YOUR_GITHUB_USERNAME
        run: |
          git clone https://$TAP_GITHUB_TOKEN@github.com/$GITHUB_USERNAME/homebrew-tap.git
          cd homebrew-tap
          sed -i '' "s/version \".*\"/version \"$VERSION\"/" Casks/md-viewer.rb
          git config user.name "GitHub Actions"
          git config user.email "actions@github.com"
          git add Casks/md-viewer.rb
          git commit -m "Update md-viewer to v$VERSION"
          git push
```

`YOUR_GITHUB_USERNAME` 換成你的帳號名稱就好。

這樣每次勾選 `Create a GitHub Release` 跑完之後，`homebrew-tap` 的版本號會自動跟著更新，完全不用手動改。


