#include <stdio.h>
#include <stdlib.h>
void checkCPUendian();

int main()
{
        checkCPUendian();
        return 0;
}

void checkCPUendian()
{
        union{
                unsigned int i;
                unsigned char s[4];
        }c;

        c.i = 0x12345678;
        printf("%s\n", (0x12 == c.s[0]) ? "大端模式" : "小端模式");
}
