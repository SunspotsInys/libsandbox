#include <stdio.h>

void main()
{
	int s,a,b,i;
	scanf("%d",&s);
	for(i = 0 ;i< s;i++)
	{
		scanf("%d %d",a,b);
		printf("%d ",a+b);
	}
}
