# 複雜數學公式測試

## 行內數學 (Inline Math)

愛因斯坦質能方程式：$E = mc^2$

黃金比例：$\phi = \frac{1 + \sqrt{5}}{2} \approx 1.6180339887$

歐拉公式：$e^{i\pi} + 1 = 0$

傅立葉轉換：$\mathcal{F}\{f(t)\} = F(\omega) = \int_{-\infty}^{\infty} f(t) e^{-i\omega t} dt$

## 區塊數學 (Display Math)

### 馬克士威方程式 (Maxwell's Equations)

$$\nabla \cdot \mathbf{E} = \frac{\rho}{\varepsilon_0}$$

$$\nabla \times \mathbf{E} = -\frac{\partial \mathbf{B}}{\partial t}$$

$$\nabla \cdot \mathbf{B} = 0$$

$$\nabla \times \mathbf{B} = \mu_0 \mathbf{J} + \mu_0 \varepsilon_0 \frac{\partial \mathbf{E}}{\partial t}$$

### 積分範例

$$\int_{0}^{\infty} x^2 e^{-x^2} dx = \frac{\sqrt{\pi}}{4}$$

$$\oint_C \mathbf{F} \cdot d\mathbf{r} = \iint_S (\nabla \times \mathbf{F}) \cdot d\mathbf{S}$$

### 矩陣運算

$$\begin{pmatrix} a_{11} & a_{12} \\ a_{21} & a_{22} \end{pmatrix}^{-1} = \frac{1}{ad - bc} \begin{pmatrix} d & -b \\ -c & a \end{pmatrix}$$

### 求和與連乘

$$\sum_{n=1}^{\infty} \frac{1}{n^2} = \frac{\pi^2}{6}$$

$$\prod_{k=1}^{n} k = n!$$

### 極限

$$\lim_{x \to 0} \frac{\sin x}{x} = 1$$

$$\lim_{n \to \infty} \left(1 + \frac{x}{n}\right)^n = e^x$$

### 希爾伯特空間

$$\langle f | g \rangle = \int_0^1 f(x) \overline{g(x)} dx$$

### 機率論

$$\mathbb{E}[X] = \int_{-\infty}^{\infty} x f(x) dx$$

$$\text{Var}(X) = \mathbb{E}[X^2] - (\mathbb{E}[X])^2$$

### 多行公式 (Aligned)

$$\begin{aligned}
\frac{\partial}{\partial t} (\rho \mathbf{v}) + \nabla \cdot (\rho \mathbf{v} \otimes \mathbf{v}) &= -\nabla p + \nabla \cdot \mathbf{\tau} + \rho \mathbf{g} \\
\frac{\partial \rho}{\partial t} + \nabla \cdot (\rho \mathbf{v}) &= 0
\end{aligned}$$

### 希臘字母與特殊符號

- $\alpha, \beta, \gamma, \delta, \epsilon, \zeta, \eta, \theta, \iota, \kappa, \lambda, \mu, \nu, \xi, \pi, \rho, \sigma, \tau, \upsilon, \phi, \chi, \psi, \omega$
- $\Gamma, \Delta, \Theta, \Lambda, \Xi, \Pi, \Sigma, \Upsilon, \Phi, \Psi, \Omega$

### 巢狀分數

$$\frac{1}{\sqrt{2\pi\sigma^2}} \exp\left(-\frac{(x-\mu)^2}{2\sigma^2}\right)$$

### 拉拉喳喳

$$A \xleftarrow{\text{下方}} B \xrightarrow[\text{下方}]{\text{上下都有}} C \xrightarrow{\text{上方}} D$$

$$\frac{\partial^2 u}{\partial t^2} = c^2 \frac{\partial^2 u}{\partial x^2}$$

### 克羅內克函數

$$\delta_{ij} = \begin{cases} 0 & \text{if } i \neq j \\ 1 & \text{if } i = j \end{cases}$$

## 混合內容

這是一段普通文字，中間有行內數學 $\int_a^b f(x)\,dx = F(b) - F(a)$，再繼續普通文字。

$$e^x = \sum_{n=0}^{\infty} \frac{x^n}{n!} = 1 + x + \frac{x^2}{2!} + \frac{x^3}{3!} + \cdots$$

再多一些行內表達 $\alpha + \beta = \gamma$，看看行內與區塊數學是否能正確區分。

### 字元計數（測試正則表達式是否正常運作）

$$E = mc^2$$ (簡單)
$$x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}$$ (二次公式)
$\sin^2\theta + \cos^2\theta = 1$ (三角函數)
$\nabla \times (\nabla \times \mathbf{A}) = \nabla(\nabla \cdot \mathbf{A}) - \nabla^2 \mathbf{A}$

## 錯誤的價格（測試跳脫）

價格：\$100
打折後：\$80

## 總結

以上測試涵蓋了：
- 基本代數
- 微積分（積分、極限）
- 矩陣運算
- 希臘字母
- 特殊函數
- 多行對齊公式
- 嵌套結構
- 跳脫字元