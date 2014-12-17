#include <sys/resource.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int f(int i){
    if(i<10){
        i++;
        char a[512*1024];
        memset(a,0,512*1024);
        struct rusage r_usage;
        getrusage(RUSAGE_SELF,&r_usage);
        printf("Memory minflt = %ld\n",r_usage.ru_minflt);
        printf("Memory jarflt = %ld\n",r_usage.ru_majflt);
        printf("Memory srss= %ld\n",r_usage.ru_isrss);
        printf("Memory usage = %ld\n",r_usage.ru_maxrss);
        sleep (1);
        return f(i);
    }else{
        return i;
    }
}
int main() {
    f(1);
    return 0;
}
