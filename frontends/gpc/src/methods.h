#ifndef GP_METHODS
#define GP_METHODS

void (*gpInit) (char*, __int32_t);
__uint32_t (*gpGetInterrupt) (void);
void (*gpSendInterrupt) (__uint32_t);
__uint32_t (*gpGetRegister) (__uint32_t);
void (*gpSetRegister) (__uint32_t, __uint32_t);
char* (*gpGetBuffer) (__uint32_t);
__uint32_t (*gpIsSubscribed) (__uint32_t);
__uint32_t (*gpIsFinished) (void);
void (*gpStep) (void);
__uint32_t (*gpGetCurrentInstruction) (void);
__uint32_t (*gpGetArgs) (void);

#endif