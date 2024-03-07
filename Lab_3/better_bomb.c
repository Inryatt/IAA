#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <errno.h>

int main()
{
    printf( "Initial PID: %d\n", getpid() );
    
    for (;;) {
        switch (fork()) {
        default:
            putchar( '.' );
            fflush( stdout );
            continue;
        }
    }

    printf( "\n" );

    return 0;
}
