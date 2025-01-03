// Copyright (c) 2023 ACOAUTO Team.
// All rights reserved.
//
// Detailed license information can be found in the LICENSE file.
//
// File: message_arm64.s Vehicle SOA protocal package.
//
// Author: Cheng.siyuan <chengsiyuan@acoinfo.com>

#include "textflag.h"

TEXT    ·smn(SB), NOSPLIT, $0
  MOVD     smn+0(FP), R0
  MOVW     $0x29, R1
  MOVB     R1, 0(R0)
  RET


TEXT    ·cmn(SB), NOSPLIT, $0
  MOVB     cmn+0(FP), R0
  MOVW     $0x29, R1
  CMP      R1, R0
  BEQ      cmny
  MOVW     $0x00, ret+16(FP)  
  RET
cmny:
  MOVW     $0xef, R3
  MOVW     R3, ret+16(FP)
  RET

TEXT    ·vs(SB), NOSPLIT, $0
  MOVB     vs+0(FP), R1
  LSR      $4, R1
  MOVW     $0x0f, R2
  AND      R2, R1
  MOVB     R1, ret+16(FP)  
  RET

TEXT    ·msgt(SB), NOSPLIT, $0
  MOVB     msgt+1(FP), R1
  MOVB     R1, ret+16(FP)  
  RET

TEXT    ·smt(SB), NOSPLIT, $0
  MOVD     smt+0(FP), R0
  MOVB     mt+8(FP), R2
  MOVB     R2, 1(R0)
  RET

TEXT    ·spe(SB), NOSPLIT, $0
  MOVD     smt+0(FP), R0
  MOVW     $0xff, R1
  MOVB     R1, 1(R0)
  RET

TEXT    ·ir(SB), NOSPLIT, $0
  MOVB     ir+1(FP), R0
  MOVW     $0x01, R1
  CMP      R1, R0
  BEQ      iry
  MOVW     $0x00, ret+16(FP)  
  RET
iry:
  MOVW     $0x01, R3
  MOVW     R3, ret+16(FP)
  RET

TEXT    ·in(SB), NOSPLIT, $0
  MOVB     in+1(FP), R0
  MOVW     $0xfe, R1
  CMP      R1, R0
  BEQ      iny
  MOVW     $0x00, ret+16(FP)  
  RET
iny:
  MOVW     $0x01, R3
  MOVW     R3, ret+16(FP)
  RET

TEXT    ·ipe(SB), NOSPLIT, $0
  MOVB     ipe+1(FP), R0
  MOVW     $0xff, R1
  CMP      R1, R0
  BEQ      ipey
  MOVW     $0x00, ret+16(FP)  
  RET
ipey:
  MOVW     $0x01, R3
  MOVW     R3, ret+16(FP)
  RET

TEXT    ·isi(SB), NOSPLIT, $0
  MOVB     isi+1(FP), R0
  MOVW     $0x00, R1
  CMP      R1, R0
  BEQ      isiy
  MOVW     $0x00, ret+16(FP)  
  RET
isiy:
  MOVW     $0x01, R3
  MOVW     R3, ret+16(FP)
  RET

TEXT    ·iss(SB), NOSPLIT, $0
  MOVB     iss+1(FP), R0
  MOVW     $0x02, R1
  CMP      R1, R0
  BEQ      issy
  MOVW     $0x00, ret+16(FP)  
  RET
issy:
  MOVW     $0x01, R3
  MOVW     R3, ret+16(FP)
  RET

TEXT    ·ius(SB), NOSPLIT, $0
  MOVB     ius+1(FP), R0
  MOVW     $0x03, R1
  CMP      R1, R0
  BEQ      iusy
  MOVW     $0x00, ret+16(FP)  
  RET
iusy:
  MOVW     $0x01, R3
  MOVW     R3, ret+16(FP)
  RET

TEXT    ·id(SB), NOSPLIT, $0
  MOVB     id+1(FP), R0
  MOVW     $0x05, R1
  CMP      R1, R0
  BEQ      idy
  MOVW     $0x00, ret+16(FP)  
  RET
idy:
  MOVW     $0x01, R3
  MOVW     R3, ret+16(FP)
  RET

TEXT    ·ip(SB), NOSPLIT, $0
  MOVB     ip+1(FP), R0
  MOVW     $0x04, R1
  CMP      R1, R0
  BEQ      ipy
  MOVW     $0x00, ret+16(FP)  
  RET
ipy:
  MOVW     $0x01, R3
  MOVW     R3, ret+16(FP)
  RET

TEXT    ·iR(SB), NOSPLIT, $0
  MOVB     iR+2(FP), R0
  AND      $0x01, R0
  MOVW     R0, ret+16(FP)  
  RET

TEXT    ·sRt(SB), NOSPLIT, $0
  MOVD     sRt+0(FP), R0
  MOVB     2(R0), R2
  MOVW     $0x01, R1
  ORR      R1, R2
  MOVB     R2, 2(R0)
  RET

TEXT    ·sRf(SB), NOSPLIT, $0
  MOVD     sRt+0(FP), R0
  MOVW     $0xfe, R1
  MOVB     2(R0), R2
  AND      R2, R1
  MOVB     R1, 2(R0)
  RET

TEXT    ·ivt(SB), NOSPLIT, $0
  MOVB    ivt+2(FP), R0
  AND     $0x02, R0
  LSR     $1, R0
  MOVW    R0, ret+16(FP)  
  RET

TEXT    ·svt(SB), NOSPLIT, $0
  MOVD     svt+0(FP), R0
  MOVB     2(R0), R2
  MOVW     $0x02, R1
  ORR      R1, R2
  MOVB     R2, 2(R0)
  RET

TEXT    ·mrm(SB), NOSPLIT, $0
  MOVB    mrm+2(FP), R0
  AND     $0x04, R0
  LSR     $2, R0
  MOVB    R0, r+16(FP)
  RET

TEXT    ·smrmg(SB), NOSPLIT, $0
  MOVD    smrmg+0(FP), R0
  MOVW    $0xfb, R1
  MOVB    2(R0), R2
  AND     R2, R1
  MOVB    R1, 2(R0)
  RET

TEXT    ·smrms(SB), NOSPLIT, $0
  MOVD     smrms+0(FP), R0
  MOVB     2(R0), R2
  MOVW     $0x04, R1
  ORR      R1, R2
  MOVB     R2, 2(R0)
  RET

TEXT    ·pl(SB), NOSPLIT, $0
  MOVB     pl+2(FP), R1
  LSR      $6, R1
  MOVW     $0x03, R2
  AND      R2, R1
  MOVB     R1, ret+16(FP)  
  RET

TEXT    ·spl(SB), NOSPLIT, $0
  MOVD     spl+0(FP), R0
  MOVB     2(R0), R2
  MOVB     spl+8(FP), R4
  MOVW     $0x3f, R1
  AND      R1, R2
  LSL      $6, R4
  MOVW     $0xc0, R1
  AND      R1, R4
  ORR      R4, R2
  MOVB     R2, 2(R0)
  RET

TEXT    ·st(SB), NOSPLIT, $0
  MOVB     st+3(FP), R0
  MOVB     R0, r+16(FP)
  RET

TEXT    ·sst(SB), NOSPLIT, $0
  MOVD     sst+0(FP), R0
  MOVD     sst+8(FP), R1
  MOVB     R1, 3(R0)
  RET

TEXT    ·sn(SB), NOSPLIT, $0
  MOVB     sn+4(FP), R0
  MOVB     sn+5(FP), R1
  MOVB     sn+6(FP), R2
  MOVB     sn+7(FP), R3
  LSL      $8, R0
  AND      $0xff00, R0
  AND      $0xff, R1
  ORR      R1, R0
  LSL      $8, R0
  AND      $0xffff00, R0
  AND      $0xff, R2
  ORR      R2, R0
  LSL      $8, R0
  AND      $0xffffff00, R0
  AND      $0xff, R3
  ORR      R3, R0
  MOVD     $0x00, ret+16(FP)
  MOVW     R0, ret+16(FP)
  RET

TEXT    ·ssn(SB), NOSPLIT, $0
  MOVD     ssn+0(FP), R0
  MOVD     ssn+8(FP), R1
  MOVD     R1, R2
  AND      $0xff ,R1
  MOVB     R1, 7(R0)
  MOVD     R2, R1
  LSR      $8, R1
  AND      $0xff ,R1
  MOVB     R1, 6(R0)
  MOVD     R2, R1
  LSR      $16, R1
  AND      $0xff ,R1
  MOVB     R1, 5(R0)
  MOVD     R2, R1
  LSR      $24, R1
  AND      $0xff ,R1
  MOVB     R1, 4(R0)
  RET

TEXT    ·tid(SB), NOSPLIT, $0
  MOVB     sn+8(FP), R0
  MOVB     sn+9(FP), R1
  LSL      $8, R0
  AND      $0xff00, R0
  AND      $0xff, R1
  ORR      R1, R0
  MOVW     R0, ret+16(FP)
  RET
  
TEXT    ·stid(SB), NOSPLIT, $0
  MOVD     ssn+0(FP), R0
  MOVD     ssn+8(FP), R1
  MOVD     R1, R2
  AND      $0xff ,R1
  MOVB     R1, 9(R0)
  MOVD     R2, R1
  LSR      $8, R1
  AND      $0xff ,R1
  MOVB     R1, 8(R0)
  RET
