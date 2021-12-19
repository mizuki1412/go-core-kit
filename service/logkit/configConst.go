package logkit

// 日志目录；默认在project.dir下
const ConfigKeyLogPath = "logger.path"

// 文件名，无后缀
const ConfigKeyLogName = "logger.name"

// 最大保留天数
const ConfigKeyLogMaxRemain = "logger.max-remain"

// 最大保留个数
const ConfigKeyLogMaxBackups = "logger.max-backups"

// 单文件最大size： megabytes
const ConfigKeyLogMaxSize = "logger.max-size"

// 关闭文件日志,只留console
const ConfigKeyLogFileOff = "logger.file-off"

const ConfigKeyLogLevel = "logger.level"
