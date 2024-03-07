#include <sys/types.h>
#include <sys/xattr.h>
#include <errno.h>
#include <grp.h>
#include <pwd.h>
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
    fprintf( stderr, "\nUsage: %s [user][:group] file [file ...]\n", program );

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
set_attr ( char * filename, char * attr_name, char * value )
{
    if (setxattr( filename, attr_name, value, strlen( value ), 0 ) == -1) {
        error_message( "Cannot add/change attribute %s of file %s: %s", attr_name, filename, strerror( errno ) );
        return -1;
    }

    return 0;
}

int
main ( int argc, char * argv[] )
{
    char * user, * group, * colon;
    uid_t uid;
    gid_t gid;
    size_t len;
    int error;

    program = argv[0];

    if (argc < 3) {
        usage( "Not enough parameters" );
        exit( 1 );
    }

    user = argv[1];

    colon = index( argv[1], ':' );

    if (colon == user) { // no user
        user = 0;
    }

    if (colon != 0) {
        *colon = 0;
        group = colon + 1;

        if (strlen( group ) == 0) {
            group = 0;
        }
    }
    else { // no group
        group = 0;
    }

    if (user == 0 && group == 0) { // Only : was provided
        usage( "No user or group was provided" );
        exit( 2 );
    }

    if (user != 0 && sscanf( user, "%d", &uid ) != 1) {
        struct passwd * entry = getpwnam( user );

        if (entry == 0) {
            usage( "User %s does not exist", user );
            exit( 2 );
        }

        uid = entry->pw_uid;
    }

    if (group != 0 && sscanf( group, "%d", &gid ) != 1) {
        struct group * entry = getgrnam( group );

        if (entry == 0) {
            usage( "Group %s does not exist", group );
            exit( 2 );
        }

        gid = entry->gr_gid;
    }

    // printf( "UID = %d, GID = %d\n", user ? uid : -1, group ? gid : -1 );

    if (user != 0) { // Convert user to a string with UID
        len = snprintf( user, 0, "%d", uid );
        len++;
        user = malloc( len );
        snprintf( user, len, "%d", uid );
    }

    if (group != 0) { // Convert group to a string with GID
        len = snprintf( group, 0, "%d", gid );
        len++;
        group = malloc( len );
        snprintf( group, len, "%d", gid );
    }

    error = 0;

    for (int i = 2; i < argc; i++) {
        if (user != 0) {
            if (set_attr( argv[i], NEW_OWNER_ATTR, user ) == -1) {
                error = 3;
            }
        }

        if (group != 0) {
            if (set_attr( argv[i], NEW_GROUP_ATTR, group ) == -1) {
                error = 3;
            }
        }
    }

    return error;
}


