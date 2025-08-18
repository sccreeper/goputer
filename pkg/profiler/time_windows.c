#include <windows.h>

unsigned long long get_perf_freq() {

    LARGE_INTEGER result;

    if (!QueryPerformanceFrequency(&result)) {
        return 0;
    }

    return (unsigned long long)result.QuadPart;

}


unsigned long long get_perf_counter() {

    LARGE_INTEGER result;

    if (!QueryPerformanceCounter(&result)) {
        return 0;
    }

    return (unsigned long long)result.QuadPart;

}