import { DesktopApp } from "../bindings/github.com/hduhelp/hdu-cli/desktop/wails"
import type { SettingsInput } from "../bindings/github.com/hduhelp/hdu-cli/desktop/wails/models"

export function currentState() {
  return DesktopApp.CurrentState()
}

export function saveSettings(input: SettingsInput) {
  return DesktopApp.SaveSettings(input)
}

export function connectNow() {
  return DesktopApp.ConnectNow()
}

export function disconnectNow() {
  return DesktopApp.DisconnectNow()
}

export function refreshStatus() {
  return DesktopApp.ReconnectStatus()
}

export function toggleLaunchAtLogin(enabled: boolean) {
  return DesktopApp.ToggleLaunchAtLogin(enabled)
}
