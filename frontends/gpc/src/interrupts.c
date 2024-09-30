#include <unistd.h>
#include <constants.h>
#include <methods.h>
#include <raylib.h>
#include <util.h>

void handleInterrupt(interrupt_t interrupt) {

    switch (interrupt)
    {
    case GP_INTERRUPT_vp:
        
        DrawPixel(
            gpGetRegister(GP_REGISTER_vx0), 
            gpGetRegister(GP_REGISTER_vy0), 
            convertColour(gpGetRegister(GP_REGISTER_vc))
        );

        break;
    case GP_INTERRUPT_va:

        DrawRectangle(
            gpGetRegister(GP_REGISTER_vx0), 
            gpGetRegister(GP_REGISTER_vy0),

            gpGetRegister(GP_REGISTER_vx1)-gpGetRegister(GP_REGISTER_vx0),
            gpGetRegister(GP_REGISTER_vy1)-gpGetRegister(GP_REGISTER_vy0),

            convertColour(gpGetRegister(GP_REGISTER_vc))
            
        );

        break;
    
    default:
        break;
    }

}