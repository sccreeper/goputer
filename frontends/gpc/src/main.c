#include <stdio.h>
#include <raylib.h>
#include <unistd.h>
#include <string.h>
#include <dlfcn.h>
#include <stdlib.h>

#define FILE_MAGIC "GPTR"
#define BUFFER_SIZE 128
#define LIBRARY_NAME "./bindings.so"

int main(int argc, char *argv[]) {

    FILE *program_file_ptr;
    char program_magic[5];
    int file_size;
    char* program_file;
    
    void *gp_handle;
    char *error;

    // Goputer methods
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

    // Load goputer library & symbols

    gp_handle = dlopen(LIBRARY_NAME, RTLD_LAZY);
    if (!gp_handle) {

        fprintf(stderr, "%s\n", dlerror());
        exit(EXIT_FAILURE);

    }

    gpInit = (void (*)(char*, __int32_t)) dlsym(gp_handle, "Init");
    gpGetInterrupt = (__uint32_t (*)(void)) dlsym(gp_handle, "GetInterrupt");
    gpSendInterrupt = (void (*)(__uint32_t)) dlsym(gp_handle, "SendInterrupt");
    gpGetRegister = (__uint32_t (*)(__uint32_t)) dlsym(gp_handle, "GetRegister");
    gpSetRegister = (void (*)(__uint32_t, __uint32_t)) dlsym(gp_handle, "SetRegister");
    gpGetBuffer = (char* (*)(__uint32_t)) dlsym(gp_handle, "GetBuffer");
    gpIsSubscribed = (__uint32_t (*)(__uint32_t)) dlsym(gp_handle, "IsSubscribed");
    gpIsFinished = (__uint32_t (*)(void)) dlsym(gp_handle, "IsFinished");
    gpStep = (void (*)(void)) dlsym(gp_handle, "Step");
    gpGetCurrentInstruction = (__uint32_t (*)(void)) dlsym(gp_handle, "GetCurrentInstruction");
    gpGetArgs = (__uint32_t (*)(void)) dlsym(gp_handle, "GetArgs");
    

    printf("Goputer C \n");

    // Open program file and check if valid

    if (access(argv[0], F_OK) != 0) {
        printf("Cannot open file %s\n", argv[1]);
        return 1;
    }
    else {
        printf("Opening file %s...\n", argv[1]);
    }

    program_file_ptr = fopen(argv[1], "r");
    fgets(program_magic, 5, program_file_ptr);

    if (strcmp(FILE_MAGIC, program_magic) != 0) {
        fprintf(stderr,"%s\n",FILE_MAGIC);
        fprintf(stderr,"%s\n", program_magic);
        fprintf(stderr,"File %s is invalid\n", argv[1]);
        exit(EXIT_FAILURE);
    }

    // Get file size & read it

    fseek(program_file_ptr, 0L, SEEK_END);
    file_size = ftell(program_file_ptr);
    rewind(program_file_ptr);

    program_file = malloc(file_size * sizeof(char));
    fgets(program_file, file_size, program_file_ptr);

    // Start VM and window

    InitWindow(640, 480, "Goputer C");

    printf("Starting VM...");

    gpInit(program_file, file_size);

    while (!WindowShouldClose()) {

        BeginDrawing();

        ClearBackground(RAYWHITE);

        DrawRectangle(0, 0, 64, 64, YELLOW);

        EndDrawing();
    }

    CloseWindow();

    return 0;
}