#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>
#include <errno.h>
#include <string.h>
#include <string.h>
#include <wait.h>

char * pid2str( pid_t pid )
{
    int len = 2;
    char * str;
    pid_t scout = pid;

    while (scout /= 10) len++;

    str = malloc( len );
    snprintf( str, len, "%u", pid );

    return str;
}

void usage( char * command )
{
    fprintf( stderr, "Usage: %s cgroup [cgroup list] -c command [command args]\n", command );
    exit( 1 );
}

int main( int argc, char ** argv )
{
    char ** command = 0;
    pid_t pid;
    char * pid_str;
    int pipe_fds[2];

    for (int i = 1; i < argc; i++) {
        if (strcmp( argv[i], "-c" ) == 0) {
            command = argv + i + 1;
            argv[i] = 0;
        }
    }

    if (command == 0 || command == argv + 2) {
        fprintf( stderr, "No cgroups where provided\n");
        usage( argv[0] );
    }

    for (int i = 1; argv[i] != 0; i++) {
        char * cgroup_pids = malloc( strlen( argv[i] ) + 14 ); // "/cgroup.procs"
        sprintf( cgroup_pids, "%s/cgroup.procs", argv[i] );

        if (access( cgroup_pids, W_OK) == -1) {
             switch(errno) {
             case ENOENT:
                 fprintf( stderr, "cgroup %s not found\n\t(looking for file %s)\n", argv[i], cgroup_pids );
                 exit( 1 );
             case EACCES:
                 fprintf( stderr, "No permission to add PID to cgroup %s\n\t(looking for file %s)\n", argv[i], cgroup_pids );
                 exit( 1 );
             default:
                 fprintf( stderr, "cgroup %s access error %d\n\t(looking for file %s)\n", argv[i], errno, cgroup_pids );
                 exit( 1 );
             }
        }

        argv[i] = cgroup_pids;
    }

    printf( "Execute the command %s with these cgroups:\n", *command );
    for (int i = 1; argv[i] != 0; i++) {
        printf( "\t%s\n", argv[i] );
    }

    pipe( pipe_fds ); // To sync parent and child

    pid = fork();

    if (pid == -1) { // error
        fprintf( stderr, "Could not fork, errno = %d\n", errno );
        exit( 2 );
    }

    if (pid == 0) { // child
        char c;
        read( pipe_fds[0], &c, 1 ); // indication to proceed from parent
        close( pipe_fds[0] );
        close( pipe_fds[1] );
        
        execv( command[0], command );

        fprintf( stderr, "Could not exec, errno = %d\n", errno );
        exit( 3 );
    }

    pid_str = pid2str( pid );

    for (int i = 1; argv[i] != 0; i++) {
        int fd = open( argv[i], O_WRONLY );

        if (fd == -1) {
            fprintf( stderr, "Cannot open to write cgroup file %s (errno = %d)\n", argv[i], errno );
            exit( 1 );
        }

        if (write( fd, pid_str, strlen( pid_str ) ) == -1) {
            fprintf( stderr, "Child process coulnt not be included in cgroup\n" );
            kill( pid, SIGKILL );
        }
        close( fd );
    }

    write( pipe_fds[1], pid_str, 1 ); // Send indication to child to proceed
    wait( 0 );

    return 0;
}
