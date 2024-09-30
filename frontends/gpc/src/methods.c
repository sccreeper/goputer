#include <unistd.h>
#include <dlfcn.h>
#include <methods.h>

void initMethods(void*);

void initMethods(void *handle) {

    gpInit = (void (*)(char*, __int32_t)) dlsym(handle, "Init");
    gpGetInterrupt = (__uint32_t (*)(void)) dlsym(handle, "GetInterrupt");
    gpSendInterrupt = (void (*)(__uint32_t)) dlsym(handle, "SendInterrupt");
    gpGetRegister = (__uint32_t (*)(__uint32_t)) dlsym(handle, "GetRegister");
    gpSetRegister = (void (*)(__uint32_t, __uint32_t)) dlsym(handle, "SetRegister");
    gpGetBuffer = (char* (*)(__uint32_t)) dlsym(handle, "GetBuffer");
    gpIsSubscribed = (__uint32_t (*)(__uint32_t)) dlsym(handle, "IsSubscribed");
    gpIsFinished = (__uint32_t (*)(void)) dlsym(handle, "IsFinished");
    gpStep = (void (*)(void)) dlsym(handle, "Step");
    gpGetCurrentInstruction = (__uint32_t (*)(void)) dlsym(handle, "GetCurrentInstruction");
    gpGetArgs = (__uint32_t (*)(void)) dlsym(handle, "GetArgs");

}