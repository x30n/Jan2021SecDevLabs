#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void flag3() {
    system("/bin/cat flag3.txt");
}

int main(void)  {

    printf("flag3() is at %p\n", blah);

    struct {
        char flag4[100];
        char buf[512];
        char root_flag;
        char stuff[100];
    }mem;

    mem.root_flag = 0;
    memcpy(mem.stuff, "stuff\n\0", strlen("stuff\n")+1);
    strcpy(mem.flag4, "flag 4 is NOT HERE!");
    printf("flag4 is at %p\n", mem.flag4);
    printf("stuff is at %p\n", mem.stuff);
    printf("Enter the password: \n");
    fflush(stdout);
    scanf("%s", mem.buf);

    if(strcmp(mem.buf, "HMMM WHAT COULD IT BE??")) {
        printf ("Wrong Password!\n");
    } else {
        printf ("Correct Password!\n");
        mem.root_flag = 1;
    }
    if(mem.root_flag)   {
        system("/bin/cat flag.txt");
    }
    printf(mem.stuff);
    if(mem.stuff[7] == '9' && mem.stuff[9] =='\x9a') {
        system("/bin/cat flag2.txt");
    }

    return 0;
}