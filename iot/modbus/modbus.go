package modbus

const FunReadCoilStatus = 0x01          // 读线圈寄存器  位  读写
const FunReadInputStatus = 0x02         // 读状态寄存器  位  只读
const FunReadHoldingRegisters = 0x03    // 读保持寄存器  2byte  读写
const FunReadInputRegisters = 0x04      // 读输入寄存器  2byte  只读
const FunForceSingleCoil = 0x05         // 写单个线圈
const FunPresetSingleRegister = 0x06    // 写单个寄存器
const FunForceMultipleCoils = 0x0F      // 写多个线圈
const FuncForceMultipleRegisters = 0x10 // 写多个寄存器
