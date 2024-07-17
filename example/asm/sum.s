
// func add(a int, b int) int{
//   return a + b
// }

// TEXT ·add(sb) $0-24
//    MOVQ args1-$0-8, CX
//    MOVQ arg2-$0-16, DP
//    ADDQ CX, DP
//    MOVQ result-$16(BF)
//    RET


TEXT ·add(SB), $0-24
    MOVQ args1+0(FP), CX  // 从参数区域加载第一个参数到CX寄存器
    MOVQ args2+8(FP), DX  // 从参数区域加载第二个参数到DX寄存器
    ADDQ DX, CX           // 将DX寄存器的值加到CX寄存器
    MOVQ CX, ret+16(FP)   // 将结果存储到返回值区域
    RET                   // 返回
