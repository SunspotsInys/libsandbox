#include <stdlib.h>
#include <stdio.h>
int
main()
{
	char *p;
	int a [100000];
	p = (char *)malloc(2048*sizeof(char));
	if (p == NULL){
		printf("内存分配出错!");
		exit(1);
	}
	free(p);
	p = NULL;
	return 0;
}
