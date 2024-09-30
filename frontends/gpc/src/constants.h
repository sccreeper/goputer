#ifndef GP_CONSTANTS
#define GP_CONSTANTS

#include <unistd.h>

// General purpose registers

typedef enum {

    GP_REGISTER_r00 = 0U,
    GP_REGISTER_r01 = 1U,
    GP_REGISTER_r02 = 2U,
    GP_REGISTER_r03 = 3U,
    GP_REGISTER_r04 = 4U,
    GP_REGISTER_r05 = 5U,
    GP_REGISTER_r06 = 6U,
    GP_REGISTER_r07 = 7U,
    GP_REGISTER_r08 = 8U,
    GP_REGISTER_r09 = 9U,
    GP_REGISTER_r10 = 10U,
    GP_REGISTER_r11 = 11U,
    GP_REGISTER_r12 = 12U,
    GP_REGISTER_r13 = 13U,
    GP_REGISTER_r14 = 14U,
    GP_REGISTER_r15 = 15U,

    //Video registers

    GP_REGISTER_vx0 = 16U,
    GP_REGISTER_vy0 = 17U,
    GP_REGISTER_vx1 = 18U,
    GP_REGISTER_vy1 = 19U,

    GP_REGISTER_vc = 20U,
    GP_REGISTER_vb = 21U,
    GP_REGISTER_vt = 22U,

    //Keyboard registers

    GP_REGISTER_kc = 23U,
    GP_REGISTER_kp = 24U,

    //Mouse registers

    GP_REGISTER_mx = 25U,
    GP_REGISTER_my = 26U,
    GP_REGISTER_mb = 27U,

    //Sound registers

    GP_REGISTER_st = 28U,
    GP_REGISTER_sv = 29U,
    GP_REGISTER_sw = 55U,

    //Accumulator and data

    GP_REGISTER_a0 = 30U,

    GP_REGISTER_d0 = 31U,
    GP_REGISTER_dl = 53U,
    GP_REGISTER_dp = 54U,

    //Stack

    GP_REGISTER_stk = 32U,
    GP_REGISTER_stz = 33U,

    GP_REGISTER_cstk = 51U,
    GP_REGISTER_cstz = 52U,

    //IO Registers

    GP_REGISTER_io00 = 34U,
    GP_REGISTER_io01 = 35U,
    GP_REGISTER_io02 = 36U,
    GP_REGISTER_io03 = 37U,
    GP_REGISTER_io04 = 38U,
    GP_REGISTER_io05 = 39U,
    GP_REGISTER_io06 = 40U,
    GP_REGISTER_io07 = 41U,

    GP_REGISTER_io08 = 42U,
    GP_REGISTER_io09 = 43U,
    GP_REGISTER_io10 = 44U,
    GP_REGISTER_io11 = 45U,
    GP_REGISTER_io12 = 46U,
    GP_REGISTER_io13 = 47U,
    GP_REGISTER_io14 = 48U,
    GP_REGISTER_io15 = 49U,

    // Program counter

    GP_REGISTER_prc = 50U,

} register_t_;

typedef enum {

    // Interrupts

    GP_INTERRUPT_ss = 0U,
    GP_INTERRUPT_sf = 1U,
    GP_INTERRUPT_va = 2U,
    GP_INTERRUPT_vp = 3U,
    GP_INTERRUPT_vt = 4U,
    GP_INTERRUPT_vc = 5U,
    GP_INTERRUPT_vl = 6U,
    GP_INTERRUPT_iof = 7U,
    GP_INTERRUPT_ioc = 8U,

    // Subscribable interrupts

    GP_INTERRUPT_mm = 9U,
    GP_INTERRUPT_mu = 10U,
    GP_INTERRUPT_md = 11U,

    GP_INTERRUPT_io08 = 12U,
    GP_INTERRUPT_io09 = 13U,
    GP_INTERRUPT_io10 = 14U,
    GP_INTERRUPT_io11 = 15U,
    GP_INTERRUPT_io12 = 16U,
    GP_INTERRUPT_io13 = 17U,
    GP_INTERRUPT_io14 = 18U,
    GP_INTERRUPT_io15 = 19U,

    GP_INTERRUPT_ku = 20U,
    GP_INTERRUPT_kd = 21U,

} interrupt_t;

typedef enum {

    // Instructions

    GP_INSTRUCTION_mov = 0U,
    GP_INSTRUCTION_jmp = 1U,

    GP_INSTRUCTION_add = 2U,
    GP_INSTRUCTION_mul = 3U,
    GP_INSTRUCTION_div = 4U,
    GP_INSTRUCTION_sub = 5U,

    GP_INSTRUCTION_cndjmp = 6U,

    GP_INSTRUCTION_gt = 7U,
    GP_INSTRUCTION_lt = 8U,

    GP_INSTRUCTION_or = 9U,
    GP_INSTRUCTION_xor = 10U,
    GP_INSTRUCTION_and = 11U,

    GP_INSTRUCTION_inv = 12U,

    GP_INSTRUCTION_eq = 13U,
    GP_INSTRUCTION_neq = 14U,

    GP_INSTRUCTION_sl = 15U,
    GP_INSTRUCTION_sr = 16U,

    GP_INSTRUCTION_int = 17U,

    GP_INSTRUCTION_lda = 18U,
    GP_INSTRUCTION_sta = 19U,

    GP_INSTRUCTION_push = 20U,
    GP_INSTRUCTION_pop = 21U,

    GP_INSTRUCTION_incr = 22U,
    GP_INSTRUCTION_decr = 23U,

    GP_INSTRUCTION_hlt = 24U,

    GP_INSTRUCTION_sqrt = 25U,

    GP_INSTRUCTION_call = 26U,
    GP_INSTRUCTION_cndcall = 27U,

    GP_INSTRUCTION_pow = 28U,

    GP_INSTRUCTION_clr = 29U,

    GP_INSTRUCTION_mod = 30U,

    GP_INSTRUCTION_emi = 31U,


} instruction_t;

#endif