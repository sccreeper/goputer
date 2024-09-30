#include <unistd.h>
#include <constants.h>

#ifndef GP_INTERRUPTS
#define GP_INTERRUPTS

void handleInterrupt(interrupt_t interrupt);

#endif