/*
 * fake_edr.c
 * A target for injection.
 * By J. Stuart McMurray
 * Created 20240909
 * Last Modified 20240910
 */

#include <unistd.h>

int
main(void) {
        for (;;)
                sleep(1024);
        return -1;
}
