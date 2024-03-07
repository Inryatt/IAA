#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <errno.h>

int main()
{
    printf( "Initial PID: %d\n", getpid() );
    
    for (int i = 0; i < 100; i++) {
        switch (fork()) {
        case 0:
             sleep( 10 );
             exit( 0 );
        case -1:
            printf( "\nError creating process %d\n", i );
            exit( errno );
        default:
            putchar( '.' );
            fflush( stdout );
            continue;
        }
    }

    printf( "\n" );

    return 0;
}
