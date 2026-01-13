export type Locale = 'zh-CN' | 'en-US'

const messages = {
  'zh-CN': {
    app: {
      title: 'TestFlight 监控',
      active: '活跃',
      settings: '设置',
    },
    header: {
      nextCheck: '下次检测',
      checking: '检测中...',
      waiting: '等待中...',
    },
    sidebar: {
      newMonitor: '添加监控',
      urlsLabel: 'TestFlight 链接',
      urlsPlaceholder: 'https://testflight.apple.com/join/...',
      urlsHint: '每行一个链接',
      intervalLabel: '检测间隔（秒）',
      durationLabel: '监控时长',
      autoStart: '添加后自动开始监控',
      addButton: '添加监控',
      telegram: 'Telegram 通知',
      botToken: 'Bot Token',
      chatId: 'Chat ID',
      enableNotify: '启用通知',
      save: '保存',
      testSend: '测试发送',
    },
    monitor: {
      available: '有位',
      full: '已满',
      checking: '检测中',
      error: '错误',
      expired: '已过期',
      loading: '加载中...',
      pause: '暂停',
      resume: '恢复',
      delete: '删除',
      edit: '编辑',
      confirmDelete: '确定要删除这个监控吗？',
      interval: '检测间隔',
      duration: '监控时长',
      forever: '永久',
      hours: '小时',
      seconds: '秒',
      cancel: '取消',
    },
    empty: {
      title: '暂无监控',
      desc: '从左侧添加 TestFlight 链接开始监控',
    },
    settings: {
      title: '设置',
      language: '语言',
      proxy: '代理设置',
      proxyEnabled: '启用代理',
      proxyUrl: '代理地址',
      proxyPlaceholder: 'http://127.0.0.1:7890',
      telegram: 'Telegram 通知',
      close: '关闭',
      saved: '保存成功',
      testSuccess: '测试消息已发送，请检查 Telegram',
      testFailed: '发送失败',
    },
  },
  'en-US': {
    app: {
      title: 'TestFlight Monitor',
      active: 'Active',
      settings: 'Settings',
    },
    header: {
      nextCheck: 'Next check in',
      checking: 'Checking...',
      waiting: 'Waiting...',
    },
    sidebar: {
      newMonitor: 'New Monitor',
      urlsLabel: 'TestFlight URLs',
      urlsPlaceholder: 'https://testflight.apple.com/join/...',
      urlsHint: 'One URL per line',
      intervalLabel: 'Interval (seconds)',
      durationLabel: 'Duration',
      autoStart: 'Auto start monitoring',
      addButton: 'Add Monitor',
      telegram: 'Telegram Notification',
      botToken: 'Bot Token',
      chatId: 'Chat ID',
      enableNotify: 'Enable Notifications',
      save: 'Save',
      testSend: 'Test Send',
    },
    monitor: {
      available: 'Available',
      full: 'Full',
      checking: 'Checking',
      error: 'Error',
      expired: 'Expired',
      loading: 'Loading...',
      pause: 'Pause',
      resume: 'Resume',
      delete: 'Delete',
      edit: 'Edit',
      confirmDelete: 'Are you sure you want to delete this monitor?',
      interval: 'Interval',
      duration: 'Duration',
      forever: 'Forever',
      hours: 'hours',
      seconds: 'seconds',
      cancel: 'Cancel',
    },
    empty: {
      title: 'No monitors',
      desc: 'Add a TestFlight URL from the sidebar to start monitoring',
    },
    settings: {
      title: 'Settings',
      language: 'Language',
      proxy: 'Proxy Settings',
      proxyEnabled: 'Enable Proxy',
      proxyUrl: 'Proxy URL',
      proxyPlaceholder: 'http://127.0.0.1:7890',
      telegram: 'Telegram Notification',
      close: 'Close',
      saved: 'Saved successfully',
      testSuccess: 'Test message sent! Check your Telegram.',
      testFailed: 'Failed to send',
    },
  },
}

export type Messages = typeof messages['zh-CN']

export function getMessages(locale: Locale): Messages {
  return messages[locale]
}

export function getStoredLocale(): Locale {
  const stored = localStorage.getItem('locale')
  if (stored === 'zh-CN' || stored === 'en-US') {
    return stored
  }
  return navigator.language.startsWith('zh') ? 'zh-CN' : 'en-US'
}

export function setStoredLocale(locale: Locale) {
  localStorage.setItem('locale', locale)
}
