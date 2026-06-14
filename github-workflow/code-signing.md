Code signing 分兩個部分：**把憑證存進 GitHub Secrets**，以及**在 workflow 裡使用它**。

---

如果之後決定申請，流程是：

去 developer.apple.com 加入 Apple Developer Program（$99/年）
在 Xcode 或 developer.apple.com 申請 Developer ID Application 憑證
憑證會自動存進 Keychain Access
再回來做剛才說的 codesign 步驟

---

## 第一部分：準備憑證

你需要從 Mac 的 Keychain 把憑證匯出：

1. 打開 **Keychain Access**
2. 找到 `Developer ID Application: Your Name (XXXXXXXXXX)`
3. 右鍵 → Export → 存成 `.p12` 檔，設一個密碼（記住這個密碼）
4. 把 `.p12` 轉成 base64：
```bash
base64 -i certificate.p12 | pbcopy
```
這樣就複製到剪貼簿了。

---

## 第二部分：存進 GitHub Secrets

去你的 repo → **Settings → Secrets and variables → Actions → New repository secret**，加這三個：

| Secret 名稱 | 內容 |
|---|---|
| `DEVELOPER_ID_CERT_BASE64` | 剛才 base64 的內容 |
| `DEVELOPER_ID_CERT_PASSWORD` | .p12 的密碼 |
| `DEVELOPER_ID_IDENTITY` | 完整名稱，如 `Developer ID Application: Your Name (XXXXXXXXXX)` |

---

## 第三部分：workflow 加入 codesign 步驟

在 Build 步驟之後、Create DMG 之前加入：

```yaml
# ── Import certificate into temporary keychain ────────
- name: Import certificate
  env:
    CERT_BASE64: ${{ secrets.DEVELOPER_ID_CERT_BASE64 }}
    CERT_PASSWORD: ${{ secrets.DEVELOPER_ID_CERT_PASSWORD }}
  run: |
    # 建立暫時的 keychain
    security create-keychain -p "" build.keychain
    security default-keychain -s build.keychain
    security unlock-keychain -p "" build.keychain
    security set-keychain-settings -lut 21600 build.keychain

    # 匯入憑證
    echo "$CERT_BASE64" | base64 --decode > certificate.p12
    security import certificate.p12 \
      -k build.keychain \
      -P "$CERT_PASSWORD" \
      -T /usr/bin/codesign
    security set-key-partition-list \
      -S apple-tool:,apple: -s \
      -k "" build.keychain

# ── Code sign ─────────────────────────────────────────
- name: Code sign
  env:
    IDENTITY: ${{ secrets.DEVELOPER_ID_IDENTITY }}
  run: |
    codesign --deep --force --verify --verbose \
      --sign "$IDENTITY" \
      --options runtime \
      "md-viewer.app"
```

---

**`--options runtime` 是必要的**，沒有這個後面 Notarization 會直接拒絕。

之後如果要加 Notarization，還需要額外的 Apple ID / App-specific password，那是另一個步驟，有需要再說。
