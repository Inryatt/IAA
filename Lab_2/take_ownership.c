#include <sys/types.h>
#include <sys/xattr.h>
#include <errno.h>
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "ownership_attrs.h"

static char * program;

void
usage( char * message, ... )
{
    va_list ap;

    va_start( ap, message );

    vfprintf( stderr, message,  ap );
    fprintf( stderr, "\nUsage: %s file [file ...]\n", program );

    va_end( ap );
}

void
error_message( char * message, ... )
{
    va_list ap;

    va_start( ap, message );

    fprintf( stderr, "%s: ", program );
    vfprintf( stderr, message, ap );
    fprintf( stderr, "\n" );

    va_end( ap );
}

int
get_attr( char * filename, char * attr_name )
{
    int len;
    char * value;

    for (;;) {
        if ((len = getxattr( filename, attr_name, 0, 0 )) > 0) { // Attribute exists
            len++;
            value = alloca( len );
            if (getxattr( filename, attr_name, value, len ) == -1) {
                if (errno != ERANGE) {
                    error_message( "Cannot get attribute %s of file %s: %s", attr_name, filename, strerror( errno ) );
                    return -1;
                }
                // else continue the cycle untyil succeeding
            }
            else {
                return atoi( value );
            }
        }
    }
}

int
main ( int argc, char * argv[] )
{
    char * attr_name;
    int error;

    program = argv[0];

    if (argc < 2) {
        usage( "Not enough parameters" );
        exit( 1 );
    }


    error = 0;

    for (int i = 1; i < argc; i++) {
        uid_t uid = get_attr( argv[i], NEW_OWNER_ATTR );
        gid_t gid = get_attr( argv[i], NEW_GROUP_ATTR );

        if (uid != -1 && getuid() != uid) {
            error_message( "Cannot take ownership of file %s: your UID (%d) is different from the one allowed (%d)", argv[i], getuid(), uid );
            uid = -1;
            error = 3;
        }

        if (gid != -1 && getgid() != gid) {
            error_message( "Cannot take ownership of file %s: your GID (%d) is different from the one allowed (%d)", argv[i], getuid(), gid );
            gid = -1;
            error = 3;
        }

        chown( argv[i], uid, gid );

        if (uid != -1) {
            removexattr( argv[i], NEW_OWNER_ATTR );
        }

        if (gid != -1) {
            removexattr( argv[i], NEW_GROUP_ATTR );
        }
    }

    return error;
}

