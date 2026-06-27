import { useEffect, useState } from "react"
import { connectNow, currentState, disconnectNow, refreshStatus, saveSettings, toggleLaunchAtLogin } from "./api"
import type { DesktopState, SettingsInput } from "./types"
import "./styles.css"

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
    <main className="shell">
      <section className="hero">
        <div className="hero-copy">
          <p className="eyebrow">HDU Campus Network</p>
          <h1>托盘常驻，开机自动连，掉线自动重连。</h1>
          <p className="lede">
            这是给 macOS 和 Windows 用的校园网桌面客户端设置窗口。保存一次账号后，应用会在托盘常驻，关闭窗口不会退出。
          </p>
        </div>
        <div className={`status-card phase-${state.phase}`}>
          <span className="status-pill">{state.online ? "ONLINE" : state.phase.toUpperCase()}</span>
          <strong>{state.phase}</strong>
          <p>{state.message || "等待操作"}</p>
        </div>
      </section>

      <section className="panel-grid">
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
              onChange={(event) => setForm((prev) => ({ ...prev, checkIntervalSeconds: Number(event.target.value) || 60 }))}
            />
          </label>
        </article>
      </section>

      <section className="panel preferences">
        <h2>常驻策略</h2>
        <label className="toggle">
          <input
            type="checkbox"
            checked={form.autoConnect}
            onChange={(event) => setForm((prev) => ({ ...prev, autoConnect: event.target.checked }))}
          />
          开机自动连接
        </label>
        <label className="toggle">
          <input
            type="checkbox"
            checked={form.autoReconnect}
            onChange={(event) => setForm((prev) => ({ ...prev, autoReconnect: event.target.checked }))}
          />
          断线自动重连
        </label>
        <label className="toggle">
          <input
            type="checkbox"
            checked={form.launchAtLogin}
            onChange={(event) => handleLaunchAtLogin(event.target.checked)}
          />
          开机自启
        </label>
      </section>

      <section className="actions">
        <button className="primary" disabled={busy} onClick={handleSaveAndConnect}>
          保存并连接
        </button>
        <button className="secondary" disabled={busy} onClick={handleRefresh}>
          刷新状态
        </button>
        <button className="ghost" disabled={busy} onClick={handleDisconnect}>
          断开连接
        </button>
      </section>
    </main>
  )
}

export default App
