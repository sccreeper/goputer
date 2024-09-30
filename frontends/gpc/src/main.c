#include <stdio.h>
#include <raylib.h>
#include <unistd.h>
#include <string.h>
#include <dlfcn.h>
#include <stdlib.h>
#include <constants.h>
#include <methods.h>
#include <interrupts.h>

#define FILE_MAGIC "GPTR"
#define BUFFER_SIZE 128
#define DEFAULT_CPS 240
#define LIBRARY_NAME "./bindings.so"

int main(int argc, char *argv[]) {

    FILE *program_file_ptr;
    char program_magic[5];
    int file_size;
    char* program_file;
    
    void *gp_handle;
    char *error;

    // Load goputer library & symbols

    gp_handle = dlopen(LIBRARY_NAME, RTLD_LAZY);
    if (!gp_handle) {

        fprintf(stderr, "%s\n", dlerror());
        exit(EXIT_FAILURE);

    }

    initMethods(gp_handle);
    

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
    if (NULL == program_file)
    {   
        printf("File allocation failed\n");
        return -1;
    }
    
    fgets(program_file, file_size, program_file_ptr);

    // Start VM and window

    InitWindow(640, 480, "Goputer C");

    printf("Starting VM...");

    gpInit(program_file, file_size);

    ClearBackground(BLACK);
    SetTargetFPS(DEFAULT_CPS);

    while (!WindowShouldClose()) {

        BeginDrawing();

        handleInterrupt(gpGetInterrupt());
        
        EndDrawing();

        if (gpIsFinished())
        {
            break;
        }

        gpStep();
        
    }

    while (1)
    {
    }
    

    CloseWindow();

    return 0;
}