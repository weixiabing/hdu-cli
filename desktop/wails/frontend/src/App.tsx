import { useEffect, useMemo, useState } from "react"
import { connectNow, currentState, disconnectNow, refreshStatus, saveSettings, toggleLaunchAtLogin } from "./api"
import type { DesktopState, SettingsInput } from "./types"
import "./styles.css"

type PageId = "home" | "connect" | "logs" | "me"

const navItems: Array<{ id: PageId; label: string; icon: string }> = [
  { id: "home", label: "主页", icon: "⌂" },
  { id: "connect", label: "连接", icon: "◫" },
  { id: "logs", label: "日志", icon: "≣" },
  { id: "me", label: "我的", icon: "◉" },
]

const fallbackState: DesktopState = {
  phase: "disconnected",
  message: "",
  online: false,
  username: "",
  endpoint: "http://192.168.112.30",
  acid: "0",
  checkIntervalSeconds: 60,
  autoConnect: true,
  autoReconnect: true,
  launchAtLogin: false,
}

function App() {
  const [activePage, setActivePage] = useState<PageId>("home")
  const [state, setState] = useState<DesktopState>(fallbackState)
  const [form, setForm] = useState<SettingsInput>({
    username: "",
    password: "",
    endpoint: fallbackState.endpoint,
    acid: fallbackState.acid,
    checkIntervalSeconds: fallbackState.checkIntervalSeconds,
    autoConnect: true,
    autoReconnect: true,
    launchAtLogin: false,
  })
  const [busy, setBusy] = useState(false)

  useEffect(() => {
    currentState()
      .then((next) => {
        setState(next)
        syncForm(next)
      })
      .catch(() => undefined)
  }, [])

  const onlineLabel = useMemo(() => {
    if (state.online) {
      return "已连接"
    }
    if (state.phase === "connecting") {
      return "连接中"
    }
    if (state.phase === "failed") {
      return "连接失败"
    }
    return "未连接"
  }, [state.online, state.phase])

  function syncForm(next: DesktopState) {
    setForm((prev) => ({
      ...prev,
      username: next.username,
      endpoint: next.endpoint,
      acid: next.acid,
      checkIntervalSeconds: next.checkIntervalSeconds,
      autoConnect: next.autoConnect,
      autoReconnect: next.autoReconnect,
      launchAtLogin: next.launchAtLogin,
    }))
  }

  async function handleSaveAndConnect() {
    setBusy(true)
    try {
      const saved = await saveSettings(form)
      setState(saved)
      const connected = await connectNow()
      setState(connected)
      setActivePage("home")
    } finally {
      setBusy(false)
    }
  }

  async function handleDisconnect() {
    setBusy(true)
    try {
      setState(await disconnectNow())
    } finally {
      setBusy(false)
    }
  }

  async function handleRefresh() {
    setBusy(true)
    try {
      setState(await refreshStatus())
    } finally {
      setBusy(false)
    }
  }

  async function handleLaunchAtLogin(enabled: boolean) {
    setForm((prev) => ({ ...prev, launchAtLogin: enabled }))
    setBusy(true)
    try {
      setState(await toggleLaunchAtLogin(enabled))
    } finally {
      setBusy(false)
    }
  }

  return (
    <main className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <div className="brand-mark">H</div>
          <div>
            <strong>HDU</strong>
            <span>校园网客户端</span>
          </div>
        </div>

        <nav className="nav-list">
          {navItems.map((item) => (
            <button
              key={item.id}
              className={`nav-item ${activePage === item.id ? "active" : ""}`}
              onClick={() => setActivePage(item.id)}
            >
              <span className="nav-icon">{item.icon}</span>
              <span>{item.label}</span>
            </button>
          ))}
        </nav>
      </aside>

      <section className="content">
        <header className="content-head">
          <div>
            <p className="eyebrow">应用程序</p>
            <h1>{navItems.find((item) => item.id === activePage)?.label}</h1>
          </div>
          <div className="head-status">
            <span className={`dot ${state.online ? "online" : ""}`} />
            <span>{onlineLabel}</span>
          </div>
        </header>

        {activePage === "home" && (
          <section className="home-view">
            <div className={`hero-status phase-${state.phase}`}>
              <div className="hero-badge">{state.online ? "⏸" : "▶"}</div>
              <div className="hero-text">
                <strong>{onlineLabel}</strong>
                <span>{state.message || "等待操作"}</span>
              </div>
            </div>

            <div className="hero-actions">
              <button className="primary" disabled={busy} onClick={handleSaveAndConnect}>
                保存并连接
              </button>
              <button className="secondary" disabled={busy} onClick={handleRefresh}>
                刷新状态
              </button>
              <button className="ghost" disabled={busy} onClick={handleDisconnect}>
                断开连接
              </button>
            </div>

            <div className="stat-grid">
              <article className="stat-card">
                <span>学号</span>
                <strong>{state.username || "未设置"}</strong>
              </article>
              <article className="stat-card">
                <span>认证地址</span>
                <strong>{state.endpoint}</strong>
              </article>
              <article className="stat-card">
                <span>AC ID</span>
                <strong>{state.acid}</strong>
              </article>
              <article className="stat-card">
                <span>巡检间隔</span>
                <strong>{state.checkIntervalSeconds}s</strong>
              </article>
            </div>
          </section>
        )}

        {activePage === "connect" && (
          <section className="stack">
            <article className="panel">
              <h2>校园网账号</h2>
              <label>
                学号
                <input
                  value={form.username}
                  onChange={(event) => setForm((prev) => ({ ...prev, username: event.target.value }))}
                  placeholder="例如 20230001"
                />
              </label>
              <label>
                密码
                <input
                  type="password"
                  value={form.password}
                  onChange={(event) => setForm((prev) => ({ ...prev, password: event.target.value }))}
                  placeholder="校园网密码"
                />
              </label>
            </article>

            <article className="panel">
              <h2>连接参数</h2>
              <label>
                认证地址
                <input
                  value={form.endpoint}
                  onChange={(event) => setForm((prev) => ({ ...prev, endpoint: event.target.value }))}
                />
              </label>
              <label>
                AC ID
                <input
                  value={form.acid}
                  onChange={(event) => setForm((prev) => ({ ...prev, acid: event.target.value }))}
                />
              </label>
              <label>
                巡检间隔（秒）
                <input
                  type="number"
                  min={5}
                  value={form.checkIntervalSeconds}
                  onChange={(event) =>
                    setForm((prev) => ({ ...prev, checkIntervalSeconds: Number(event.target.value) || 60 }))
                  }
                />
              </label>
            </article>

            <article className="panel">
              <h2>连接策略</h2>
              <SettingRow
                title="开机自动连接"
                description="系统启动后自动完成首次连接"
                checked={form.autoConnect}
                onChange={(checked) => setForm((prev) => ({ ...prev, autoConnect: checked }))}
              />
              <SettingRow
                title="断线自动重连"
                description="掉线后自动恢复连接"
                checked={form.autoReconnect}
                onChange={(checked) => setForm((prev) => ({ ...prev, autoReconnect: checked }))}
              />
              <div className="row-actions">
                <button className="primary" disabled={busy} onClick={handleSaveAndConnect}>
                  保存并连接
                </button>
              </div>
            </article>
          </section>
        )}

        {activePage === "logs" && (
          <section className="stack">
            <article className="panel">
              <h2>最近状态</h2>
              <div className="log-list">
                <div className="log-item">
                  <span>当前</span>
                  <strong>{onlineLabel}</strong>
                  <p>{state.message || "暂无消息"}</p>
                </div>
                <div className="log-item">
                  <span>掉线重连</span>
                  <strong>{state.autoReconnect ? "开启" : "关闭"}</strong>
                  <p>用于保持校园网长时间在线</p>
                </div>
                <div className="log-item">
                  <span>巡检间隔</span>
                  <strong>{state.checkIntervalSeconds}s</strong>
                  <p>定时检查连通性</p>
                </div>
              </div>
            </article>
          </section>
        )}

        {activePage === "me" && (
          <section className="stack">
            <article className="panel">
              <h2>应用程序</h2>
              <SettingRow
                title="退出时最小化"
                description="关闭窗口后不退出，只隐藏到托盘"
                checked={true}
                onChange={() => undefined}
              />
              <SettingRow
                title="开机自启"
                description="跟随系统自动启动应用"
                checked={form.launchAtLogin}
                onChange={handleLaunchAtLogin}
              />
              <SettingRow
                title="显示连接页面"
                description="在导航栏中显示连接页"
                checked={true}
                onChange={() => undefined}
              />
            </article>

            <article className="panel">
              <h2>关于</h2>
              <div className="about-grid">
                <div>
                  <span>版本</span>
                  <strong>HDU Network</strong>
                </div>
                <div>
                  <span>平台</span>
                  <strong>macOS / Windows</strong>
                </div>
              </div>
            </article>
          </section>
        )}
      </section>
    </main>
  )
}

function SettingRow({
  title,
  description,
  checked,
  onChange,
}: {
  title: string
  description: string
  checked: boolean
  onChange: (checked: boolean) => void
}) {
  return (
    <div className="setting-row">
      <div className="setting-copy">
        <strong>{title}</strong>
        <span>{description}</span>
      </div>
      <label className="switch">
        <input type="checkbox" checked={checked} onChange={(event) => onChange(event.target.checked)} />
        <span className="switch-track" />
      </label>
    </div>
  )
}

export default App
