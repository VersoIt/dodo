export enum LogLevel {
  DEBUG = 0,
  INFO = 1,
  WARN = 2,
  ERROR = 3,
}

class Logger {
  private level: LogLevel = LogLevel.INFO

  constructor() {
    if (import.meta.env.DEV) {
      this.level = LogLevel.DEBUG
    }
  }

  private formatMessage(level: string, message: string, data?: any) {
    const timestamp = new Date().toISOString()
    const color = this.getColor(level)
    console.log(
      `%c[${timestamp}] [${level}] ${message}`,
      `color: ${color}; font-weight: bold`,
      data || ''
    )
  }

  private getColor(level: string): string {
    switch (level) {
      case 'DEBUG': return '#7f8c8d'
      case 'INFO': return '#2ecc71'
      case 'WARN': return '#f1c40f'
      case 'ERROR': return '#e74c3c'
      default: return '#34495e'
    }
  }

  debug(msg: string, data?: any) {
    if (this.level <= LogLevel.DEBUG) this.formatMessage('DEBUG', msg, data)
  }

  info(msg: string, data?: any) {
    if (this.level <= LogLevel.INFO) this.formatMessage('INFO', msg, data)
  }

  warn(msg: string, data?: any) {
    if (this.level <= LogLevel.WARN) this.formatMessage('WARN', msg, data)
  }

  error(msg: string, error?: any) {
    if (this.level <= LogLevel.ERROR) {
      const timestamp = new Date().toISOString()
      console.error(`%c[${timestamp}] [ERROR] ${msg}`, 'color: #e74c3c; font-weight: bold', error || '')
    }
  }
}

export const logger = new Logger()
